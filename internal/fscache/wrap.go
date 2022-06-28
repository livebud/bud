package fscache

import (
	"io/fs"

	"github.com/livebud/bud/package/log"
)

func Wrap(fsys fs.FS, log log.Interface, fsname string) *WrapFS {
	return &WrapFS{
		fsname: fsname,
		fsys:   fsys,
		log:    log,
		cache:  New(log),
	}
}

type WrapFS struct {
	fsname string
	fsys   fs.FS
	log    log.Interface
	cache  *Cache
}

func (w *WrapFS) mount(name string) string {
	return w.fsname + ":" + name
}

func (w *WrapFS) Open(name string) (fs.File, error) {
	if w.cache.Has(name) {
		w.log.Debug("fscache: cache hit", "file", name)
		return w.cache.Open(name)
	}
	w.log.Debug("fscache: cache miss", "file", name)
	file, err := w.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	entry, err := From(file)
	if err != nil {
		return nil, err
	}
	w.cache.Set(name, entry)
	return w.cache.Open(name)
}

func (w *WrapFS) Clear() {
	w.cache.Clear()
}
