package searcher

import (
	"bufio"
	"fmt"
	"io/fs"
	"regexp"
	"slices"
	"strings"
	"sync"
)

type Searcher struct {
	FS fs.FS
}

func (s *Searcher) Search(word string) (files []string, err error) {
	type indexWordPath struct {
		word string
		path string
	}

	var wg sync.WaitGroup
	index := make(map[string]map[string]struct{})

	indexChan := make(chan indexWordPath)
	errChan := make(chan error)

	go func() {
		for v := range errChan {
			fmt.Println(v)
		}
	}()

	err = fs.WalkDir(s.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			wg.Add(1)
			go func() {
				defer wg.Done()
				file, err := s.FS.Open(path)
				if err != nil {
					errChan <- err
					return
				}
				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					line := scanner.Text()
					words := extractWords(line)
					for _, w := range words {
						indexChan <- indexWordPath{
							word: w,
							path: path,
						}
					}
				}
			}()
		}
		return nil
	})

	// Close fileChan when all goroutines are done
	go func() {
		wg.Wait()
		close(indexChan)
		close(errChan)
	}()

	for v := range indexChan {
		if _, ok := index[v.word]; !ok {
			index[v.word] = make(map[string]struct{})
			index[v.word][v.path] = struct{}{}
		} else {
			index[v.word][v.path] = struct{}{}
		}
	}

	var answer []string
	for w := range index {
		if w == strings.ToLower(word) {
			idx := index[w]
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
