package storage

import (
	"fmt"
	"sync"
)

type URLStore interface {
	Get(id string) (string, error)
	Add(id string, url string) error
}

// MemoryURLStore is a simple in-memory storage
type MemoryURLStore struct {
	URLs map[string]string
	mx   sync.RWMutex
}

// NewMemoryURLStore creates a new in-memory URL store
func NewMemoryURLStore() *MemoryURLStore {
	return &MemoryURLStore{
		URLs: make(map[string]string),
	}
}

// Get returns the URL with the given ID
func (s *MemoryURLStore) Get(id string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	url, exists := s.URLs[id]
	if !exists {
		return "", fmt.Errorf("short URL not found")
	}
	return url, nil
}

// Add adds a new URL to the store
func (s *MemoryURLStore) Add(id string, url string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.URLs[id] = url

	return nil
}
