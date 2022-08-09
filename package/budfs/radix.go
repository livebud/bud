package budfs

import (
	goradix "github.com/armon/go-radix"
)

func newRadix() *radix {
	return &radix{
		tree: goradix.New(),
	}
}

type radix struct {
	tree *goradix.Tree
	keys []string
}

func (t *radix) Keys() []string {
	return t.keys
}

// Set fn
func (t *radix) Set(key string, generator Generator) {
	if _, ok := t.tree.Get(key); ok {
		return
	}
	t.tree.Insert(key, generator)
	t.keys = append(t.keys, key)
}

// Get fn
func (t *radix) Get(key string) (generator Generator, ok bool) {
	existing, ok := t.tree.Get(key)
	if !ok {
		return nil, false
	}
	return existing.(Generator), true
}

// GetByPrefix fn
func (t *radix) GetByPrefix(prefix string) (key string, generator Generator, ok bool) {
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

func (r *radix) Upsert(key string, generator Generator) {
	if _, ok := r.tree.Get(key); ok {
		return
	}
	r.tree.Insert(key, generator)
	r.keys = append(r.keys, key)
}
