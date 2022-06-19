// Package fscache provides a caching layer on top of virtual filesystems. The
// cache itself looks like a regular fs.File, with one difference. It's
// idemptotent so it can be read many times and it will always return the same
// result. This is similar to how fstest.MapFS works. This cache greatly reduces
// the number of repeated file system calls in Bud.
//
// TODO: rework this package. The original plan was to share a single cache
// across all wrapped filesystems. This can lead to caching issues because the
// same valid file in two different filesystems will be cached as the same file.
// For example, "." should return different results for each filesystem. This
// has been worked around by creating a new cache each time we wrap a
// filesystem. However, when we start wanting to incrementally updating the
// cache,  we'll only have a single file path, so we'll need to be able to clear
// multiple caches at that time.
package fscache

import (
	"fmt"
	"io/fs"
	"path"
	"sync"

	"github.com/livebud/bud/package/log"
)

func New(log log.Interface) *Cache {
	return &Cache{log: log}
}

type Cache struct {
	log log.Interface
	sm  sync.Map
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

func (c *Cache) Wrap(fsname string, fsys fs.FS) fs.FS {
	return &Wrapped{c.log, fsname, fsys, New(c.log)}
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
