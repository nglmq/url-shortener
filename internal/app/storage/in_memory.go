package storage

import (
	"fmt"
	"sync"
)

type URLStore interface {
	Get(id string) (string, error)
	Add(id string, url string) error
}

type MemoryURLStore struct {
	URLs map[string]string
	mx   sync.RWMutex
}

func NewMemoryURLStore() *MemoryURLStore {
	return &MemoryURLStore{
		URLs: make(map[string]string),
	}
}

func (s *MemoryURLStore) Get(id string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	url, exists := s.URLs[id]
	if !exists {
		return "", fmt.Errorf("short URL not found")
	}
	return url, nil
}

func (s *MemoryURLStore) Add(id string, url string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.URLs[id] = url

	return nil
}
