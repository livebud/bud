package plugin

import (
	"io/fs"

	"gitlab.com/mnm/bud/go/mod"
)

func New(module *mod.Module) fs.FS {
	return &fileSystem{module}
}

type fileSystem struct {
	module *mod.Module
}

var _ fs.FS = (*fileSystem)(nil)
var _ fs.ReadDirFS = (*fileSystem)(nil)

func (f *fileSystem) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (f *fileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return nil, fs.ErrNotExist
}
