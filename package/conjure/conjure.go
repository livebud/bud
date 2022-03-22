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

func (f *FileSystem) GenerateFile(path string, fn func(f *File) error) {
	f.radix.Set(path, &fileGenerator{path, fn})
	f.filler[path] = &fstest.MapFile{}
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

// type Map map[string]Generator

// func (m Map) Generator(f *FileSystem) {
// 	for key, value := range m {
// 		value.Generate()
// 	}
// }

// type GenerateDir func(d *Dir) error

// func (fn GenerateDir) Generator(f *FileSystem, path string) {
// 	f.GenerateDir(path, fn)
// }

func (f *FileSystem) GenerateDir(path string, fn func(d *Dir) error) {
	f.filler[path] = &fstest.MapFile{Mode: fs.ModeDir}
	f.radix.Set(path, &dirg{
		path:   path,
		fn:     fn,
		filler: f.filler,
	})
}

func (f *FileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

func (f *FileSystem) ServeFile(path string, fn func(f *File) error) {
	f.radix.Set(path, &serverg{
		path: path,
		fn:   fn,
	})
}

func (f *FileSystem) FileServer(path string, server FileServer) {
	f.ServeFile(path, server.ServeFile)
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
