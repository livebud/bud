package genfs

import (
	"io/fs"

	"github.com/livebud/bud/package/virtual"
)

type ExternalFile struct {
	target string
}

func (e *ExternalFile) Target() string {
	return e.target
}

func (e *ExternalFile) Mode() fs.FileMode {
	return fs.FileMode(0)
}

type ExternalGenerator interface {
	GenerateExternal(fsys FS, file *ExternalFile) error
}

type externalGenerator struct {
	cache Cache
	fn    func(fsys FS, e *ExternalFile) error
	genfs fs.FS
	path  string
}

func (e *externalGenerator) Generate(target string) (fs.File, error) {
	if target != e.path {
		return nil, formatError(fs.ErrNotExist, "%q path doesn't match %q target", e.path, target)
	}
	if _, ok := e.cache.Get(target); ok {
		return nil, fs.ErrNotExist
	}
	scoped := &scopedFS{e.cache, e.genfs, target}
	file := &ExternalFile{target}
	if err := e.fn(scoped, file); err != nil {
		return nil, err
	}
	vfile := &virtual.File{
		Path: e.path,
	}
	e.cache.Set(target, vfile)
	return nil, fs.ErrNotExist
}
