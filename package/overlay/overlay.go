package overlay

import (
	"context"

	"github.com/livebud/bud/package/budfs/genfs"

	"github.com/livebud/bud/package/budfs/mergefs"

	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/virtual"
	"github.com/livebud/bud/internal/virtual/vcache"

	"io/fs"

	"github.com/livebud/bud/internal/dag"
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
	cache := vcache.New()
	gen := genfs.New(cache)
	merged := mergefs.Merge(gen, pluginFS)
	syncfs := wrapFS(cache, merged, log)
	dag := dag.New()
	clear := func() {
		cache.Clear()
	}
	return &FileSystem{gen, dag, syncfs, module, clear}, nil
}

// Serve is just load without the cache
// TODO: consolidate
func Serve(log log.Interface, module *gomod.Module) (*Server, error) {
	pluginFS, err := pluginfs.Load(module)
	if err != nil {
		return nil, err
	}
	cache := vcache.Discard
	gen := genfs.New(cache)
	merged := mergefs.Merge(gen, pluginFS)
	dag := dag.New()
	clear := func() {
		cache.Clear()
	}
	return &FileSystem{gen, dag, merged, module, clear}, nil
}

type Server = FileSystem

type F interface {
	fs.FS
	Link(from, to string)
}

type FileSystem struct {
	gen    *genfs.FileSystem
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
	f.gen.GenerateFile(path, func(file *genfs.File) error {
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
	f.gen.GenerateDir(path, func(dir *genfs.Dir) error {
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
	f.gen.ServeFile(path, func(file *genfs.File) error {
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
	f.gen.Mount(dir, fsys)
}

// wrapFS wraps a filesystem with a cache
func wrapFS(c vcache.Cache, fsys fs.FS, log log.Interface) fs.FS {
	return virtual.Opener(func(name string) (fs.File, error) {
		entry, ok := c.Get(name)
		if ok {
			log.Debug("cachefs: cache hit", "file", name)
			return entry.Open(), nil
		}
		log.Debug("cachefs: cache miss", "file", name)
		file, err := fsys.Open(name)
		if err != nil {
			return nil, err
		}
		entry, err = virtual.From(file)
		if err != nil {
			return nil, err
		}
		c.Set(name, entry)
		return entry.Open(), nil
	})
}
