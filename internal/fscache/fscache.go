package fscache

import (
	"fmt"
	"io/fs"
	"path"
	"sync"
)

func New() *Cache {
	return &Cache{}
}

type Cache struct {
	sm sync.Map
}

func (c *Cache) Has(path string) (ok bool) {
	_, ok = c.sm.Load(path)
	return ok
}

func (c *Cache) Set(path string, entry Entry) {
	c.sm.Store(path, entry)
}

func (c *Cache) Open(path string) (fs.File, error) {
	value, ok := c.sm.Load(path)
	if !ok {
		return nil, fs.ErrNotExist
	}
	entry, ok := value.(Entry)
	if !ok {
		return nil, fmt.Errorf("virtual: invalid map entry")
	}
	return entry.open(), nil
}

func (c *Cache) Keys() (keys []string) {
	c.sm.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(string))
		return true
	})
	return keys
}

func (c *Cache) Clear() {
	c.sm = sync.Map{}
}

func (c *Cache) Wrap(name string, fsys fs.FS) fs.FS {
	return &Wrapped{name, fsys, c}
}

// File events

// Update event
func (c *Cache) Update(name string) {
	c.sm.Delete(name)
}

// Delete event
func (c *Cache) Delete(name string) {
	c.sm.Delete(name)
	c.sm.Delete(path.Dir(name))
}

// Create event
func (c *Cache) Create(name string) {
	c.sm.Delete(path.Dir(name))
}
