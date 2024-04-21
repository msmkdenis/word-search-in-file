package searcher

import (
	"bufio"
	"context"
	"golang.org/x/sync/errgroup"
	"io/fs"
	"regexp"
	"slices"
	"strings"
	"sync"
)

type Searcher struct {
	FS fs.FS
}

type index struct {
	words map[string]map[string]struct{}
	mu    sync.RWMutex
}

func (i *index) Add(word string, path string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if _, ok := i.words[word]; !ok {
		i.words[word] = make(map[string]struct{})
		i.words[word][path] = struct{}{}
	} else {
		i.words[word][path] = struct{}{}
	}
}

func (s *Searcher) Search(ctx context.Context, word string) (files []string, err error) {
	var paths []string

	err = fs.WalkDir(s.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	idx := &index{words: make(map[string]map[string]struct{})}

	grp, ctx := errgroup.WithContext(context.Background())
	for _, path := range paths {
		grp.Go(func() error {
			if ctx.Done() != nil {
				return ctx.Err()
			}
			file, err := s.FS.Open(path)
			defer file.Close()

			if err != nil {
				return err
			}
			scanner := bufio.NewScanner(file)
			wordSet := make(map[string]struct{})
			for scanner.Scan() {
				line := scanner.Text()
				words := extractWords(line)
				for _, w := range words {
					wordSet[w] = struct{}{}
				}
				for w := range wordSet {
					idx.Add(w, path)
				}
			}
			return nil
		})
	}

	if err := grp.Wait(); err != nil {
		return nil, err
	}

	var answer []string
	for w := range idx.words {
		if w == strings.ToLower(word) {
			idx := idx.words[w]
			for p := range idx {
				answer = append(answer, strings.Split(p, ".")[0])
			}
		}
	}
	slices.Sort(answer)

	return answer, nil
}

func extractWords(text string) []string {
	// Удаляем знаки препинания и преобразуем текст в нижний регистр
	text = strings.ToLower(text)
	re := regexp.MustCompile(`[.,!?;:"'()-]`)
	text = re.ReplaceAllString(text, "")

	// Разделяем текст на слова
	words := strings.Fields(text)
	return words
}
