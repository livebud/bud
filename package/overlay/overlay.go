package overlay

import (
	"context"

	"gitlab.com/mnm/bud/internal/dsync"

	"gitlab.com/mnm/bud/internal/dag"
	"gitlab.com/mnm/bud/package/fs"
	"gitlab.com/mnm/bud/package/merged"

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
	merged := merged.Merge(cfs, pluginFS)
	dag := dag.New()
	return &FileSystem{cfs, dag, merged, module}, nil
}

type F interface {
	fs.OpenFS
	Link(from, to string)
}

type FileSystem struct {
	cfs    *conjure.FileSystem
	dag    *dag.Graph
	fsys   *merged.FS
	module *gomod.Module
}

func (f *FileSystem) Link(from, to string) {
}

func (f *FileSystem) Open(name string) (fs.File, error) {
	return f.fsys.Open(name)
}

func (f *FileSystem) OpenContext(ctx context.Context, name string) (fs.File, error) {
	return f.fsys.OpenContext(ctx, name)
}

var _ fs.FS = (*FileSystem)(nil)

func (f *FileSystem) GenerateFile(path string, fn func(ctx context.Context, fsys F, file *File) error) {
	f.cfs.GenerateFile(path, func(ctx context.Context, file *conjure.File) error {
		return fn(ctx, f, &File{File: file})
	})
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

func (f *FileSystem) GenerateDir(path string, fn func(ctx context.Context, fsys F, dir *Dir) error) {
	f.cfs.GenerateDir(path, func(ctx context.Context, dir *conjure.Dir) error {
		return fn(ctx, f, &Dir{f, dir})
	})
}

func (f *FileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

// Sync the overlay to the filesystem
func (f *FileSystem) Sync(dir string) error {
	return dsync.Dir(f.fsys, dir, f.module.DirFS(dir), ".")
}
