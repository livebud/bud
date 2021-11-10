package parser

import (
	"go/ast"
	"sync"
)

// newCache initializes a new cache for the parser
func newCache() *cache {
	return &cache{
		files: make(map[string]*ast.File),
	}
}

// Cache is an opaque data structure for caching speeding up parsing
type cache struct {
	mu    sync.RWMutex
	files map[string]*ast.File
}

// Get a package from a file
func (c *cache) Get(filepath string) (*ast.File, bool) {
	c.mu.RLock()
	pkg, ok := c.files[filepath]
	c.mu.RUnlock()
	return pkg, ok
}

// Store a package in the cache from a file
func (c *cache) Set(filepath string, file *ast.File) {
	c.mu.Lock()
	c.files[filepath] = file
	c.mu.Unlock()
}
