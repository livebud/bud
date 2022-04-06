package fscache

import (
	"io/fs"
)

type Wrapped struct {
	name string
	fs   fs.FS
	c    *Cache
}

func (w *Wrapped) Open(name string) (fs.File, error) {
	if w.c.Has(name) {
		// fmt.Println("  ", w.name, "cache hit", name)
		return w.c.Open(name)
	}
	// fmt.Println("  ", w.name, "cache miss", name)
	file, err := w.fs.Open(name)
	if err != nil {
		return nil, err
	}
	entry, err := From(file)
	if err != nil {
		return nil, err
	}
	w.c.Set(name, entry)
	return w.c.Open(name)
}
