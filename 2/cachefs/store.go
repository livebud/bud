package cachefs

import (
	"io/fs"
	"path"
	"sync"
	"testing/fstest"
)

// Cache storage
func Cache() *Store {
	return &Store{
		mu: sync.RWMutex{},
		fs: fstest.MapFS{},
	}
}

type Store struct {
	mu sync.RWMutex
	fs fstest.MapFS
}

func (s *Store) Has(name string) (ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok = s.fs[name]
	return ok
}

func (s *Store) Set(name string, file *fstest.MapFile) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fs[name] = file
}

func (s *Store) Open(name string) (file fs.File, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.fs.Open(name)
}

// Events

// Update event
func (s *Store) Update(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.fs, name)
}

// Delete event
func (s *Store) Delete(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.fs, name)
	delete(s.fs, path.Dir(name))
}

// Create event
func (s *Store) Create(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.fs, path.Dir(name))
}
