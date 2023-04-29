//go:generate mockgen -source=storage.go -destination=mock.go -package=counter
package counter

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("counter not found")

type Storage interface {
	Set(counter *Counter) error
	Get(id string) (*Counter, error)
	Delete(id string) error
}

type MemoryStorage struct {
	mu       sync.RWMutex
	counters map[string]*Counter
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{counters: map[string]*Counter{}}
}

func (s *MemoryStorage) Set(counter *Counter) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.counters[counter.ID] = counter

	return nil
}

func (s *MemoryStorage) Get(id string) (*Counter, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	counter, ok := s.counters[id]
	if !ok {
		return nil, ErrNotFound
	}

	return counter, nil
}

func (s *MemoryStorage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.counters[id]; !ok {
		return ErrNotFound
	}

	delete(s.counters, id)

	return nil
}
