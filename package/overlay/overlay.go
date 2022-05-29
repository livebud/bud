package overlay

import (
	"context"

	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/fscache"
	"github.com/livebud/bud/internal/pubsub"

	"io/fs"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/package/merged"

	"github.com/livebud/bud/package/conjure"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/pluginfs"
)

// Load the overlay filesystem
func Load(module *gomod.Module) (*FileSystem, error) {
	cache := fscache.New()
	pluginFS, err := pluginfs.Load(module)
	if err != nil {
		return nil, err
	}
	cfs := conjure.New()
	merged := merged.Merge(cache.Wrap("cfs", cfs), cache.Wrap("pluginfs", pluginFS))
	dag := dag.New()
	ps := pubsub.New()
	return &FileSystem{cache, cfs, dag, merged, module, ps}, nil
}

// Serve is just load without the cache
// TODO: consolidate
func Serve(module *gomod.Module) (*Server, error) {
	pluginFS, err := pluginfs.Load(module)
	if err != nil {
		return nil, err
	}
	cfs := conjure.New()
	merged := merged.Merge(cfs, pluginFS)
	dag := dag.New()
	ps := pubsub.New()
	return &FileSystem{fscache.New(), cfs, dag, merged, module, ps}, nil
}

type Server = FileSystem

type F interface {
	fs.FS
	Link(from, to string)
}

type FileSystem struct {
	cache  *fscache.Cache
	cfs    *conjure.FileSystem
	dag    *dag.Graph
	fsys   fs.FS
	module *gomod.Module
	ps     pubsub.Client
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
	f.cfs.GenerateFile(path, func(file *conjure.File) error {
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
	f.cfs.GenerateDir(path, func(dir *conjure.Dir) error {
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
	f.cfs.ServeFile(path, func(file *conjure.File) error {
		return fn(context.TODO(), f, &File{file})
	})
}

func (f *FileSystem) FileServer(path string, server FileServer) {
	f.ServeFile(path, server.ServeFile)
}

// Sync the overlay to the filesystem
func (f *FileSystem) Sync(dir string) error {
	// Clear the filesystem cache before syncing again
	f.cache.Clear()
	return dsync.Dir(f.fsys, dir, f.module.DirFS(dir), ".")
}
