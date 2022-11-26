package storage

import (
	"fmt"
	"sync"
)

type MemoryStore struct {
	mu   sync.Mutex
	data map[string][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string][]byte),
	}
}

func (ms *MemoryStore) Add(key string, value []byte) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.data[key] = value
}

func (ms *MemoryStore) Get(key string) ([]byte, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	value, ok := ms.data[key]
	if !ok {
		return nil, fmt.Errorf("key not found: key = %s", key)
	}
	return value, nil
}

func (ms *MemoryStore) Delete(key string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	delete(ms.data, key)
}
