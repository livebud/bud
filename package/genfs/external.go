package genfs

import (
	"io/fs"

	"github.com/livebud/bud/package/virtual"
)

type External struct {
	target string
}

func (e *External) Path() string {
	return e.target
}

func (e *External) Target() string {
	return e.target
}

func (e *External) Mode() fs.FileMode {
	return fs.FileMode(0)
}

type ExternalGenerator interface {
	GenerateExternal(fsys FS, file *External) error
}

type externalGenerator struct {
	cache Cache
	fn    func(fsys FS, e *External) error
	genfs fs.FS
	path  string
}

func (e *externalGenerator) Generate(target string) (fs.File, error) {
	if target != e.path {
		return nil, formatError(fs.ErrNotExist, "%q path doesn't match %q target", e.path, target)
	}
	if _, err := e.cache.Get(target); err == nil {
		return nil, fs.ErrNotExist
	}
	scoped := &scopedFS{e.cache, e.genfs, target}
	file := &External{target}
	if err := e.fn(scoped, file); err != nil {
		return nil, err
	}
	vfile := &virtual.File{
		Path: target,
		Mode: file.Mode(),
	}
	if err := e.cache.Set(target, vfile); err != nil {
		return nil, err
	}
	return nil, fs.ErrNotExist
}
