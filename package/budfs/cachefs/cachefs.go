package cachefs

import (
	"fmt"
	"io/fs"
	"sync"

	"github.com/livebud/bud/internal/virtual"
	"github.com/livebud/bud/package/log"
)

func New(log log.Interface) *Cache {
	return &Cache{log, sync.Map{}}
}

type Cache struct {
	log log.Interface
	sm  sync.Map
}

func (c *Cache) Has(path string) (ok bool) {
	_, ok = c.sm.Load(path)
	return ok
}

func (c *Cache) Set(path string, entry virtual.Entry) {
	c.sm.Store(path, entry)
}

func (c *Cache) Get(path string) (entry virtual.Entry, ok bool) {
	value, ok := c.sm.Load(path)
	if !ok {
		return nil, false
	}
	entry, ok = value.(virtual.Entry)
	if !ok {
		return nil, false
	}
	return entry, ok
}

func (c *Cache) Delete(path string) {
	c.sm.Delete(path)
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

func (c *Cache) Wrap(fsys fs.FS) fs.FS {
	return &cachedFS{c, fsys, c.log}
}

type cachedFS struct {
	cache *Cache
	fsys  fs.FS
	log   log.Interface
}

func (f *cachedFS) open(name string) (fs.File, error) {
	entry, ok := f.cache.Get(name)
	if !ok {
		return nil, fmt.Errorf("cachefs: unable to open cached file %q. %w", name, fs.ErrNotExist)
	}
	return entry.Open(), nil
}

func (f *cachedFS) Open(name string) (fs.File, error) {
	if f.cache.Has(name) {
		f.log.Debug("cachefs: cache hit", "file", name)
		return f.open(name)
	}
	f.log.Debug("cachefs: cache miss", "file", name)
	file, err := f.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	entry, err := virtual.From(file)
	if err != nil {
		return nil, err
	}
	f.cache.Set(name, entry)
	return f.open(name)
}
