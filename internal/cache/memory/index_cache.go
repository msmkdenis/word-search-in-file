package memory

import (
	"sync"
)

type IndexCache struct {
	index map[string]map[string]map[string]struct{}
	mu    sync.RWMutex
}

func NewIndexCache() *IndexCache {
	return &IndexCache{
		index: make(map[string]map[string]map[string]struct{}),
	}
}

func (i *IndexCache) GetIndex(path string) map[string]map[string]struct{} {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.index[path]
}

func (i *IndexCache) AddIndex(path string, idx map[string]map[string]struct{}) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.index[path] = idx
}
