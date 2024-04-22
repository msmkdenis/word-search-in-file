package model

import (
	"slices"
	"strings"
	"sync"
)

type Index struct {
	mu        sync.RWMutex
	words     map[string]map[string]struct{}
	searchIdx map[string][]string
}

func NewIndex() *Index {
	return &Index{
		words:     make(map[string]map[string]struct{}),
		searchIdx: make(map[string][]string),
	}
}

func (i *Index) GetIndex() map[string][]string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.searchIdx
}

func (i *Index) SetIndex(idx map[string]map[string]struct{}) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.words = idx
}

func (i *Index) BuildSearchIndex() {
	i.mu.Lock()
	defer i.mu.Unlock()
	for k, v := range i.words {
		var files []string
		for f := range v {
			files = append(files, strings.Split(f, ".")[0])
		}
		slices.Sort(files)
		i.searchIdx[k] = files
	}
}

func (i *Index) GetFiles(word string) []string {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.searchIdx[strings.ToLower(word)]
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
