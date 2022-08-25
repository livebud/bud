package overlay

import (
	"context"

	"github.com/livebud/bud/package/budfs/genfs"

	"github.com/livebud/bud/package/budfs/mergefs"

	"github.com/livebud/bud/internal/dsync"

	"io/fs"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/package/budfs/cachefs"
	"github.com/livebud/bud/package/log"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/pluginfs"
)

// Load the overlay filesystem
func Load(log log.Interface, module *gomod.Module) (*FileSystem, error) {
	pluginFS, err := pluginfs.Load(module)
	if err != nil {
		return nil, err
	}
	gfs := genfs.New()
	cacheFS := cachefs.New(log)
	merged := mergefs.Merge(cacheFS.Wrap(gfs), pluginFS)
	dag := dag.New()
	clear := func() {
		cacheFS.Clear()
	}
	return &FileSystem{gfs, dag, merged, module, clear}, nil
}

// Serve is just load without the cache
// TODO: consolidate
func Serve(log log.Interface, module *gomod.Module) (*Server, error) {
	pluginFS, err := pluginfs.Load(module)
	if err != nil {
		return nil, err
	}
	gfs := genfs.New()
	merged := mergefs.Merge(gfs, pluginFS)
	dag := dag.New()
	clear := func() {}
	return &FileSystem{gfs, dag, merged, module, clear}, nil
}

type Server = FileSystem

type F interface {
	fs.FS
	Link(from, to string)
}

type FileSystem struct {
	gfs    *genfs.FileSystem
	dag    *dag.Graph
	fsys   fs.FS
	module *gomod.Module
	clear  func() // Clear the cache
}

func (f *FileSystem) Link(from, to string) {
}

func (f *FileSystem) Open(name string) (fs.File, error) {
	file, err := f.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

var _ fs.FS = (*FileSystem)(nil)

type GenerateFile func(ctx context.Context, fsys F, file *File) error

func (fn GenerateFile) GenerateFile(ctx context.Context, fsys F, file *File) error {
	if err := fn(ctx, fsys, file); err != nil {
		return err
	}
	return nil
}

func (f *FileSystem) GenerateFile(path string, fn func(ctx context.Context, fsys F, file *File) error) {
	f.gfs.GenerateFile(path, func(file *genfs.File) error {
		if err := fn(context.TODO(), f, &File{File: file}); err != nil {
			return err
		}
		return nil
	})
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

type GenerateDir func(ctx context.Context, fsys F, dir *Dir) error

func (fn GenerateDir) GenerateDir(ctx context.Context, fsys F, dir *Dir) error {
	if err := fn(ctx, fsys, dir); err != nil {
		return err
	}
	return nil

}

func (f *FileSystem) GenerateDir(path string, fn func(ctx context.Context, fsys F, dir *Dir) error) {
	f.gfs.GenerateDir(path, func(dir *genfs.Dir) error {
		if err := fn(context.TODO(), f, &Dir{f, dir}); err != nil {
			return err
		}
		return nil
	})
}

func (f *FileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

func (f *FileSystem) ServeFile(path string, fn func(ctx context.Context, fsys F, file *File) error) {
	f.gfs.ServeFile(path, func(file *genfs.File) error {
		return fn(context.TODO(), f, &File{file})
	})
}

func (f *FileSystem) FileServer(path string, server FileServer) {
	f.ServeFile(path, server.ServeFile)
}

// Sync the overlay to the filesystem
func (f *FileSystem) Sync(dir string) error {
	// Clear the filesystem cache before syncing again
	f.clear()
	return dsync.To(f.fsys, f.module.DirFS("."), dir)
}

// Mount a filesystem to a dir
func (f *FileSystem) Mount(dir string, fsys fs.FS) {
	f.gfs.Mount(dir, fsys)
}
