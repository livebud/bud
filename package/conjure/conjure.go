package conjure

import (
	"io/fs"
	"testing/fstest"
)

func New() *FileSystem {
	f := &FileSystem{
		radix:  newRadix(),
		filler: fstest.MapFS{},
	}
	return f
}

type FileSystem struct {
	radix  *radix
	filler fstest.MapFS
}

type FS interface {
	fs.FS
}

func (f *FileSystem) File(path string, generator FileGenerator) {
	f.radix.Set(path, &fileGenerator{path, generator})
	f.filler[path] = &fstest.MapFile{}
}

func (f *FileSystem) Dir(path string, generator DirGenerator) {
	f.filler[path] = &fstest.MapFile{Mode: fs.ModeDir}
	f.radix.Set(path, &dirg{
		path:   path,
		gen:    generator,
		filler: f.filler,
	})
}

func (f *FileSystem) ServeFile(path string, server FileServer) {
	f.radix.Set(path, &serverg{
		path:   path,
		server: server,
	})
}

func (f *FileSystem) Open(target string) (fs.File, error) {
	dir := &Dir{
		gpath:  ".",
		tpath:  target,
		Mode:   fs.ModeDir,
		filler: f.filler,
		radix:  f.radix,
	}
	return dir.open(target)
}
