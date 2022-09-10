package budfs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/glob"
	"github.com/livebud/bud/internal/linkmap"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/orderedset"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/budfs/genfs"
	"github.com/livebud/bud/package/budfs/mergefs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/virtual"
	"github.com/livebud/bud/package/virtual/vcache"
)

func New(cache vcache.Cache, module *gomod.Module, log log.Interface) *FileSystem {
	gen := genfs.New()
	merged := mergefs.Merge(module, vcache.Wrap(cache, gen, log))
	closer := new(once.Closer)
	lmap := linkmap.New()
	return &FileSystem{cache, closer, gen, lmap, log, module, merged}
}

type FileSystem struct {
	cache  vcache.Cache
	closer *once.Closer
	gen    *genfs.FileSystem
	lmap   *linkmap.Map
	log    log.Interface
	module *gomod.Module
	fsys   fs.FS
}

type genfsFile = genfs.File

type File struct {
	*genfsFile
	link *linkmap.List
}

func (f *File) Link(tos ...string) {
	match := map[string]bool{}
	for _, to := range tos {
		match[to] = true
	}
	f.link.Add(func(path string) bool {
		return match[path]
	})
}

type FS struct {
	closer *once.Closer
	ctx    context.Context
	fsys   fs.FS
	link   *linkmap.List
	log    log.Interface
	path   string
}

var _ fs.FS = (*FS)(nil)
var _ fs.GlobFS = (*FS)(nil)

// Context returns the context
func (f *FS) Context() context.Context {
	return f.ctx
}

// Open implement fs.FS
func (f *FS) Open(name string) (fs.File, error) {
	file, err := f.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	f.link.Add(func(path string) bool {
		return path == name
	})
	return file, nil
}

func (f *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	des, err := fs.ReadDir(f.fsys, name)
	if err != nil {
		return nil, err
	}
	f.link.Add(func(path string) bool {
		return path == name || filepath.Dir(path) == name
	})
	return des, nil
}

// Glob implements fs.GlobFS
func (f *FS) Glob(pattern string) (matches []string, err error) {
	// Compile the pattern into a glob matcher
	matcher, err := glob.Compile(pattern)
	if err != nil {
		return nil, err
	}
	f.link.Add(func(path string) bool {
		return matcher.Match(path)
	})
	// Base is a minor optimization to avoid walking the entire tree
	bases, err := glob.Bases(pattern)
	if err != nil {
		return nil, err
	}
	// Compute the matches for each base
	for _, base := range bases {
		results, err := f.glob(matcher, base)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		matches = append(matches, results...)
	}
	// Deduplicate the matches
	return orderedset.Strings(matches...), nil
}

func (f *FS) glob(matcher glob.Matcher, base string) (matches []string, err error) {
	// Walk the directory tree, filtering out non-valid paths
	err = fs.WalkDir(f.fsys, base, valid.WalkDirFunc(func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// If the paths match, add it to the list of matches
		if matcher.Match(path) {
			matches = append(matches, path)
		}
		return nil
	}))
	if err != nil {
		return nil, err
	}
	// return the list of matches
	return matches, nil
}

// Call the fn when we close the filesystem
func (f *FS) Defer(fn func() error) {
	f.closer.Closes = append(f.closer.Closes, fn)
}

type Dir struct {
	cache   vcache.Cache
	closer  *once.Closer
	fsys    fs.FS
	dir     *genfs.Dir
	linkset *linkmap.Map
	log     log.Interface
}

func (d *Dir) Target() string {
	return d.dir.Target()
}

func (d *Dir) Relative() string {
	return d.dir.Relative()
}

func (d *Dir) Path() string {
	return d.dir.Path()
}

func (d *Dir) GenerateFile(path string, fn func(fsys *FS, file *File) error) {
	link := d.linkset.Scope(path)
	fsys := &FS{d.closer, context.TODO(), d.fsys, link, d.log, path}
	d.dir.GenerateFile(path, func(file *genfs.File) error {
		return fn(fsys, &File{file, link})
	})
}

func (d *Dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *Dir) GenerateDir(path string, fn func(fsys *FS, dir *Dir) error) {
	fsys := &FS{d.closer, context.TODO(), d.fsys, d.linkset.Scope(path), d.log, path}
	d.dir.GenerateDir(path, func(dir *genfs.Dir) error {
		d := &Dir{d.cache, d.closer, d.fsys, dir, d.linkset, d.log}
		if err := fn(fsys, d); err != nil {
			return err
		}
		d.cache.Set(d.dir.Path(), &virtual.Dir{
			Path:    d.dir.Path(),
			Mode:    d.dir.Mode(),
			Entries: d.dir.Entries(),
		})
		return nil
	})
}

func (d *Dir) DirGenerator(path string, generator DirGenerator) {
	d.GenerateDir(path, generator.GenerateDir)
}

type FileGenerator interface {
	GenerateFile(fsys *FS, file *File) error
}

type GenerateFile func(fsys *FS, file *File) error

func (fn GenerateFile) GenerateFile(fsys *FS, file *File) error {
	return fn(fsys, file)
}

type DirGenerator interface {
	GenerateDir(fsys *FS, dir *Dir) error
}

type GenerateDir func(fsys *FS, dir *Dir) error

func (fn GenerateDir) GenerateDir(fsys *FS, dir *Dir) error {
	return fn(fsys, dir)
}

type EmbedFile genfs.EmbedFile

var _ FileGenerator = (*EmbedFile)(nil)

func (e *EmbedFile) GenerateFile(fsys *FS, file *File) error {
	file.Data = e.Data
	return nil
}

func (f *FileSystem) Open(name string) (fs.File, error) {
	f.log.Debug("budfs: open", "name", name)
	return f.fsys.Open(name)
}

func (f *FileSystem) Close() error {
	return f.closer.Close()
}

func (f *FileSystem) GenerateFile(path string, fn func(fsys *FS, file *File) error) {
	link := f.lmap.Scope(path)
	fsys := &FS{f.closer, context.TODO(), f, link, f.log, path}
	f.gen.GenerateFile(path, func(file *genfs.File) error {
		return fn(fsys, &File{file, link})
	})
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

func (f *FileSystem) GenerateDir(path string, fn func(fsys *FS, dir *Dir) error) {
	fsys := &FS{f.closer, context.TODO(), f, f.lmap.Scope(path), f.log, path}
	f.gen.GenerateDir(path, func(dir *genfs.Dir) error {
		d := &Dir{f.cache, f.closer, f, dir, f.lmap, f.log}
		if err := fn(fsys, d); err != nil {
			return err
		}
		f.cache.Set(d.dir.Path(), &virtual.Dir{
			Path:    d.dir.Path(),
			Mode:    d.dir.Mode(),
			Entries: d.dir.Entries(),
		})
		return nil
	})
}

func (f *FileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

func (f *FileSystem) ServeFile(path string, fn func(fsys *FS, file *File) error) {
	link := f.lmap.Scope(path)
	fsys := &FS{f.closer, context.TODO(), f, link, f.log, path}
	f.gen.ServeFile(path, func(file *genfs.File) error {
		return fn(fsys, &File{file, link})
	})
}

func (f *FileSystem) FileServer(path string, generator FileGenerator) {
	f.ServeFile(path, generator.GenerateFile)
}

// Sync the overlay to the filesystem
func (f *FileSystem) Sync(to string) error {
	fmt.Println("syncing...")
	err := dsync.To(f, f.module.DirFS("."), to)
	fmt.Println("synced...")
	return err
}

func (f *FileSystem) Mount(path string, fsys fs.FS) {
	f.gen.Mount(path, fsys)
}

// Change updates the cache
func (f *FileSystem) Change(changes ...string) {
	for i := 0; i < len(changes); i++ {
		change := changes[i]
		if f.cache.Has(change) {
			f.log.Debug("budfs: cache", "delete", change)
			f.cache.Delete(change)
		}
		for path, fn := range f.lmap.Map() {
			if f.cache.Has(path) && fn(change) {
				changes = append(changes, path)
			}
		}
	}
}
