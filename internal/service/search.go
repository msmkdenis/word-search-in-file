package service

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"regexp"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/msmkdenis/word-search-in-file/internal/model"
)

type IndexCache interface {
	GetIndex(dirPath string) map[string]map[string]struct{}
	AddIndex(dirPath string, idx map[string]map[string]struct{})
}

type SearcherService struct {
	index    *model.Index
	idxCache IndexCache
	workers  int
}

func NewSearcher(idxCache IndexCache, workers int) *SearcherService {
	idx := model.NewIndex()
	return &SearcherService{
		index:    idx,
		idxCache: idxCache,
		workers:  workers,
	}
}

func (s *SearcherService) Search(ctx context.Context, word string, fs model.FileSystem) (files []string, err error) {
	if idx := s.idxCache.GetIndex(fs.DirPath); idx != nil {
		s.index.SetIndex(idx)
		answer := s.index.SearchFiles(word)
		if len(answer) == 0 {
			return nil, nil
		}
		return answer, nil
	}

	paths, err := s.getFilePaths(fs.FS)
	if err != nil {
		return nil, fmt.Errorf("get file paths: %w", err)
	}

	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(s.workers)
	for _, path := range paths {
		path := path
		grp.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				file, err := fs.FS.Open(path)
				defer func() {
					if errClose := file.Close(); err != nil {
						err = errClose
					}
				}()

				if err != nil {
					return fmt.Errorf("open file: %w", err)
				}

				scanner := bufio.NewScanner(file)
				wordSet := make(map[string]struct{})
				for scanner.Scan() {
					line := scanner.Text()
					words := s.extractWords(line)
					for _, w := range words {
						wordSet[w] = struct{}{}
					}
					for w := range wordSet {
						s.index.AddWordFile(w, path)
					}
				}
				return nil
			}
		})
	}

	if err := grp.Wait(); err != nil {
		return nil, fmt.Errorf("make indexes errgroup: %w", err)
	}

	s.idxCache.AddIndex(fs.DirPath, s.index.GetIndex())

	answer := s.index.SearchFiles(word)
	if len(answer) == 0 {
		return nil, nil
	}

	return answer, nil
}

func (s *SearcherService) getFilePaths(fileSystem fs.FS) ([]string, error) {
	var paths []string
	err := fs.WalkDir(fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walk dir: %w", err)
		}
		if !d.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk dir: %w", err)
	}

	return paths, nil
}

func (s *SearcherService) extractWords(text string) []string {
	// Удаляем знаки препинания и преобразуем текст в нижний регистр
	text = strings.ToLower(text)
	re := regexp.MustCompile(`[.,!?;:"'()-]`)
	text = re.ReplaceAllString(text, "")

	// Разделяем текст на слова
	words := strings.Fields(text)
	return words
}
