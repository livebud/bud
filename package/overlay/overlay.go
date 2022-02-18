package overlay

import (
	"io/fs"

	"gitlab.com/mnm/bud/internal/dsync"

	"gitlab.com/mnm/bud/internal/dag"
	"gitlab.com/mnm/bud/package/mergefs"

	"gitlab.com/mnm/bud/package/conjure"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/pluginfs"
)

// Load the overlay filesystem
func Load(module *gomod.Module) (*FileSystem, error) {
	pluginFS, err := pluginfs.Load(module)
	if err != nil {
		return nil, err
	}
	cfs := conjure.New()
	fsys := mergefs.Merge(cfs, pluginFS)
	dag := dag.New()
	return &FileSystem{cfs, dag, fsys, module}, nil
}

type F interface {
	fs.FS
	Link(from, to string)
}

type FileSystem struct {
	cfs    *conjure.FileSystem
	dag    *dag.Graph
	fsys   fs.FS
	module *gomod.Module
}

func (f *FileSystem) overlay() {}

func (f *FileSystem) Link(from, to string) {
}

func (f *FileSystem) Open(name string) (fs.File, error) {
	return f.fsys.Open(name)
}

var _ fs.FS = (*FileSystem)(nil)

func (f *FileSystem) GenerateFile(path string, fn func(fsys F, file *File) error) {
	f.cfs.GenerateFile(path, func(file *conjure.File) error {
		return fn(f, &File{File: file})
	})
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

func (f *FileSystem) GenerateDir(path string, fn func(fsys F, dir *Dir) error) {
	f.cfs.GenerateDir(path, func(dir *conjure.Dir) error {
		return fn(f, &Dir{f, dir})
	})
}

func (f *FileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

// Sync the overlay to the filesystem
func (f *FileSystem) Sync(dir string) error {
	return dsync.Dir(f.fsys, dir, f.module.DirFS(dir), ".")
}
