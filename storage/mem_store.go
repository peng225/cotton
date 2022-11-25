package storage

import "fmt"

type MemoryStore struct {
	data map[string][]byte
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string][]byte),
	}
}

func (ms *MemoryStore) Add(key string, value []byte) {
	ms.data[key] = value
}

func (ms *MemoryStore) Get(key string) ([]byte, error) {
	value, ok := ms.data[key]
	if !ok {
		return nil, fmt.Errorf("key not found: key = %s", key)
	}
	return value, nil
}

func (ms *MemoryStore) Delete(key string) {
	delete(ms.data, key)
}
