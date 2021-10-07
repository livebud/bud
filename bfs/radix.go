package bfs

import "github.com/armon/go-radix"

func newRadix() *radixTree {
	return &radixTree{
		tree: radix.New(),
	}
}

type radixTree struct {
	tree *radix.Tree
}

// Set fn
func (t *radixTree) Set(key string, generator Generator) {
	t.tree.Insert(key, generator)
}

// Get fn
func (t *radixTree) Get(key string) (generator Generator, ok bool) {
	existing, ok := t.tree.Get(key)
	if !ok {
		return nil, false
	}
	return existing.(Generator), true
}

// GetByPrefix fn
func (t *radixTree) GetByPrefix(prefix string) (key string, generator Generator, ok bool) {
	key, existing, ok := t.tree.LongestPrefix(prefix)
	if !ok {
		return "", nil, false
	}
	generator, ok = existing.(Generator)
	if !ok {
		return "", nil, false
	}
	return key, generator, ok
}
