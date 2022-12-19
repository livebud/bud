package genfs

import (
	"io/fs"
	"path"

	"github.com/livebud/bud/package/virtual"
)

func newFiller() *filler {
	return &filler{virtual.Tree{}}
}

type filler struct {
	m virtual.Tree
}

func (f *filler) Open(name string) (fs.File, error) {
	return f.m.Open(name)
}

func (f *filler) Has(path string) bool {
	return f.m[path] != nil
}

func (f *filler) Insert(fpath string, mode fs.FileMode) {
	if fpath == "." || f.m[fpath] != nil {
		return
	}
	f.m[fpath] = &virtual.File{Path: fpath, Mode: mode}
	f.Insert(path.Dir(fpath), fs.ModeDir)
}
