package virtual

import (
	"fmt"
	"io/fs"
	"path"
	"sync"
)

func FileMap() *Map {
	return &Map{}
}

type Map struct {
	sm sync.Map
}

func (m *Map) Has(path string) (ok bool) {
	_, ok = m.sm.Load(path)
	return ok
}

func (m *Map) Set(path string, entry Entry) {
	m.sm.Store(path, entry)
}

func (m *Map) Open(path string) (fs.File, error) {
	value, ok := m.sm.Load(path)
	if !ok {
		return nil, fs.ErrNotExist
	}
	entry, ok := value.(Entry)
	if !ok {
		return nil, fmt.Errorf("virtual: invalid map entry")
	}
	return entry.open(), nil
}

// File events

// Update event
func (m *Map) Update(name string) {
	m.sm.Delete(name)
}

// Delete event
func (m *Map) Delete(name string) {
	m.sm.Delete(name)
	m.sm.Delete(path.Dir(name))
}

// Create event
func (m *Map) Create(name string) {
	m.sm.Delete(path.Dir(name))
}
