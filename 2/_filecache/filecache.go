package filecache

import (
	"io/fs"
	"path"
	"sync"

	"gitlab.com/mnm/bud/2/virtual"
)

func New() *Cache {
	return &Cache{
		mu:      sync.RWMutex{},
		entries: map[string]virtual.Entry{},
	}
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]virtual.Entry
}

func (c *Cache) Has(path string) (ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok = c.entries[path]
	return ok
}

func (c *Cache) Set(path string, entry virtual.Entry) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[path] = entry
}

func (c *Cache) Open(name string) (file fs.File, err error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return entry, nil
}

// Events

// Update event
func (c *Cache) Update(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, name)
}

// Delete event
func (c *Cache) Delete(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, name)
	delete(c.entries, path.Dir(name))
}

// Create event
func (c *Cache) Create(name string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, path.Dir(name))
}
