package vcache

import (
	"sync"

	"github.com/livebud/bud/package/virtual"
)

type Cache interface {
	Has(path string) (ok bool)
	Get(path string) (entry virtual.Entry, ok bool)
	Set(path string, entry virtual.Entry)
	Delete(path string)
	Range(fn func(path string, entry virtual.Entry) bool)
	Clear()
}

func New() Cache {
	return &memory{}
}

type memory struct {
	sm sync.Map
}

func (c *memory) Has(path string) (ok bool) {
	_, ok = c.sm.Load(path)
	return ok
}

func (c *memory) Set(path string, entry virtual.Entry) {
	c.sm.Store(path, entry)
}

func (c *memory) Get(path string) (entry virtual.Entry, ok bool) {
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

func (c *memory) Delete(path string) {
	c.sm.Delete(path)
}

func (c *memory) Range(fn func(path string, entry virtual.Entry) bool) {
	c.sm.Range(func(key, value interface{}) bool {
		path, ok := key.(string)
		if !ok {
			return true
		}
		entry, ok := value.(virtual.Entry)
		if !ok {
			return true
		}
		return fn(path, entry)
	})
}

func (c *memory) Clear() {
	c.sm.Range(func(key, value interface{}) bool {
		c.sm.Delete(key)
		return true
	})
}
