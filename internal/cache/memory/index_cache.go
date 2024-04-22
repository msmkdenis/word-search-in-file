package memory

import (
	"fmt"
	"strings"
	"sync"
)

type IndexCache struct {
	searchIdx map[string]map[string][]string
	mu        sync.RWMutex
}

func NewIndexCache() *IndexCache {
	return &IndexCache{
		searchIdx: make(map[string]map[string][]string),
	}
}

func (i *IndexCache) GetFiles(path, word string) ([]string, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	filesIdx, ok := i.searchIdx[path]
	if !ok {
		return nil, false
	}

	files := filesIdx[strings.ToLower(word)]
	return files, true
}

func (i *IndexCache) SetIndex(path string, idx map[string][]string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	fmt.Println(idx)
	i.searchIdx[path] = idx
}
