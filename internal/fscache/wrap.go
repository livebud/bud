package fscache

import (
	"context"

	"gitlab.com/mnm/bud/package/fs"
)

func Wrap(fs fs.FS) *Wrapped {
	return &Wrapped{fs, New()}
}

type Wrapped struct {
	fs fs.FS
	c  *Cache
}

func (w *Wrapped) Open(name string) (fs.File, error) {
	return w.OpenContext(context.Background(), name)
}

func (w *Wrapped) OpenContext(ctx context.Context, name string) (fs.File, error) {
	if w.c.Has(name) {
		return w.c.Open(name)
	}
	file, err := fs.Open(ctx, w.fs, name)
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
