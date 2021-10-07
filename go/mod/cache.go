package mod

import "sync"

// newcache initializes a new cache for mod
func newCache() *cache {
	return &cache{
		modfiles: make(map[string]File),
	}
}

// cache is an opaque data structure for caching speeding up parsing
type cache struct {
	mu       sync.RWMutex
	modfiles map[string]File
}

// Get a package from a directory
func (c *cache) Get(dir string) (File, bool) {
	c.mu.RLock()
	modfile, ok := c.modfiles[dir]
	c.mu.RUnlock()
	return modfile, ok
}

// Set a modfile in the cache from a directory
func (c *cache) Set(dir string, modfile File) {
	c.mu.Lock()
	c.modfiles[dir] = modfile
	c.mu.Unlock()
}
