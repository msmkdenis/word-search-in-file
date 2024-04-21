package model

import (
	"slices"
	"strings"
	"sync"
)

type Index struct {
	words map[string]map[string]struct{}
	mu    sync.RWMutex
}

func NewIndex() *Index {
	return &Index{words: make(map[string]map[string]struct{})}
}

func (i *Index) GetIndex() map[string]map[string]struct{} {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.words
}

func (i *Index) SetIndex(idx map[string]map[string]struct{}) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.words = idx
}

func (i *Index) AddWordFile(word string, path string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if _, ok := i.words[word]; !ok {
		i.words[word] = make(map[string]struct{})
		i.words[word][path] = struct{}{}
	} else {
		i.words[word][path] = struct{}{}
	}
}

func (i *Index) SearchFiles(word string) []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var answer []string
	for w := range i.words {
		if w == strings.ToLower(word) {
			idx := i.words[w]
			for p := range idx {
				answer = append(answer, strings.Split(p, ".")[0])
			}
		}
	}
	// сортировка необходимо для прохождения заданного теста
	slices.Sort(answer)

	return answer
}
