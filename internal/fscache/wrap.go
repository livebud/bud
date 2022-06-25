package fscache

import (
	"io/fs"

	"github.com/livebud/bud/package/log"
)

type Wrapped struct {
	log    log.Interface
	fsname string
	fs     fs.FS
	c      *Cache
}

func (w *Wrapped) mount(name string) string {
	return w.fsname + ":" + name
}

func (w *Wrapped) Open(name string) (fs.File, error) {
	if w.c.Has(name) {
		w.log.Debug("fscache: cache hit", "file", name)
		return w.c.Open(name)
	}
	w.log.Debug("fscache: cache miss", "file", name)
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
