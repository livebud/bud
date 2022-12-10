package genfs

import (
	"io/fs"

	rdx "github.com/armon/go-radix"
)

func newRadix() *radix {
	return &radix{rdx.New()}
}

type radix struct {
	tree *rdx.Tree
}

type generator interface {
	Generate(target string) (fs.File, error)
}

func (r *radix) Insert(path string, gen generator) {
	r.tree.Insert(path, gen)
}

func (r *radix) Get(path string) (gen generator, found bool) {
	value, found := r.tree.Get(path)
	if !found {
		return nil, false
	}
	return value.(generator), true
}

func (r *radix) FindByPrefix(path string) (gen generator, prefix string, found bool) {
	prefix, value, found := r.tree.LongestPrefix(path)
	if !found {
		return nil, prefix, false
	}
	return value.(generator), prefix, true
}
