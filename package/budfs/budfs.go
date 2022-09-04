package budfs

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"path"

	"github.com/livebud/bud/package/virtual/vcache"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/package/budfs/genfs"
	"github.com/livebud/bud/package/budfs/mergefs"
	"github.com/livebud/bud/package/log"
)

func New(cache vcache.Cache, module *gomod.Module, log log.Interface) *FileSystem {
	gen := genfs.New()
	merged := mergefs.Merge(gen, module)
	return &FileSystem{
		cache:  cache,
		closer: new(once.Closer),
		dag:    dag.New(),
		gen:    gen,
		log:    log,
		module: module,
		fsys:   merged,
	}
}

type FileSystem struct {
	cache  vcache.Cache
	closer *once.Closer
	dag    *dag.Graph
	gen    *genfs.FileSystem
	log    log.Interface
	module *gomod.Module
	fsys   fs.FS
}

type File = genfs.File

type FS struct {
	closer *once.Closer
	ctx    context.Context
	dag    *dag.Graph
	fsys   fs.FS
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
	f.dag.Link(f.path, name)
	return file, nil
}

// Glob implements fs.GlobFS
func (f *FS) Glob(name string) ([]string, error) {
	matches, err := fs.Glob(f.fsys, name)
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		f.dag.Link(f.path, match)
	}
	return matches, nil
}

// Link a dependency where "from" depends on "to".
func (f *FS) Link(from, to string) {
	f.dag.Link(from, to)
}

func (f *FS) Defer(fn func() error) {
	f.closer.Closes = append(f.closer.Closes, fn)
}

type Dir struct {
	cache  vcache.Cache
	closer *once.Closer
	dag    *dag.Graph
	fsys   fs.FS
	dir    *genfs.Dir
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
	fsys := &FS{d.closer, context.TODO(), d.dag, d.fsys, path}
	d.dir.GenerateFile(path, func(file *genfs.File) error { return fn(fsys, file) })
}

func (d *Dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *Dir) GenerateDir(path string, fn func(fsys *FS, dir *Dir) error) {
	fsys := &FS{d.closer, context.TODO(), d.dag, d.fsys, path}
	d.dir.GenerateDir(path, func(dir *genfs.Dir) error {
		d := &Dir{d.cache, d.closer, d.dag, d.fsys, dir}
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
	entry, ok := f.cache.Get(name)
	f.log.Debug("budfs: cache hit", "name", name)
	if ok {
		return entry.Open(), nil
	}
	f.log.Debug("budfs: cache miss", "name", name)
	file, err := f.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	entry, err = f.toEntry(name, file)
	if err != nil {
		return nil, err
	}
	f.cache.Set(name, entry)
	return entry.Open(), nil
}

func (f *FileSystem) Close() error {
	return f.closer.Close()
}

func (f *FileSystem) GenerateFile(path string, fn func(fsys *FS, file *File) error) {
	fsys := &FS{f.closer, context.TODO(), f.dag, f.fsys, path}
	f.gen.GenerateFile(path, func(file *genfs.File) error { return fn(fsys, file) })
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

func (f *FileSystem) GenerateDir(path string, fn func(fsys *FS, dir *Dir) error) {
	fsys := &FS{f.closer, context.TODO(), f.dag, f.fsys, path}
	f.gen.GenerateDir(path, func(dir *genfs.Dir) error {
		d := &Dir{f.cache, f.closer, f.dag, f, dir}
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
	fsys := &FS{f.closer, context.TODO(), f.dag, f.fsys, path}
	f.gen.ServeFile(path, func(file *genfs.File) error { return fn(fsys, file) })
}

func (f *FileSystem) FileServer(path string, generator FileGenerator) {
	f.ServeFile(path, generator.GenerateFile)
}

// Sync the overlay to the filesystem
func (f *FileSystem) Sync(to string) error {
	// Clear the filesystem cache before syncing again
	f.cache.Clear()
	return dsync.To(f, f.module.DirFS("."), to)
}

func (f *FileSystem) Mount(path string, fsys fs.FS) {
	f.gen.Mount(path, fsys)
}

func (f *FileSystem) Print() string {
	return f.dag.String()
}

// Create event
func (f *FileSystem) Create(filepath string) {
	dir := path.Dir(filepath)
	f.cache.Delete(dir)
	// Delete all downstream cache entries that depended on dir
	for _, ancestor := range f.dag.Ancestors(dir) {
		f.cache.Delete(ancestor)
	}
}

// Update event
func (f *FileSystem) Update(filepath string) {
	f.cache.Delete(filepath)
	// Delete all downstream cache entries that depended on filepath
	for _, ancestor := range f.dag.Ancestors(filepath) {
		f.cache.Delete(ancestor)
	}
}

// Delete event
func (f *FileSystem) Delete(filepath string) {
	dir := path.Dir(filepath)
	f.cache.Delete(filepath)
	f.cache.Delete(dir)
	// Delete all downstream cache entries that depended on filepath and dir
	for _, ancestor := range f.dag.Ancestors(filepath) {
		f.cache.Delete(ancestor)
	}
	for _, ancestor := range f.dag.Ancestors(dir) {
		f.cache.Delete(ancestor)
	}
}

// toEntry converts a fs.File into a reusable virtual entry
func (f *FileSystem) toEntry(fullpath string, file fs.File) (virtual.Entry, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	// Handle files
	if !stat.IsDir() {
		// Read the data fully
		data, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return &virtual.File{
			Path:    fullpath,
			Data:    data,
			ModTime: stat.ModTime(),
			Mode:    stat.Mode(),
			Sys:     stat.Sys(),
		}, nil
	}
	vdir := &virtual.Dir{
		Path:    fullpath,
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
		Sys:     stat.Sys(),
	}
	dir, ok := file.(fs.ReadDirFile)
	if !ok {
		return nil, fmt.Errorf("budfs: dir does not have ReadDir: %s", fullpath)
	}
	des, err := dir.ReadDir(-1)
	if err != nil {
		return nil, err
	}
	for _, de := range des {
		stat, err := f.toEntryInfo(fullpath, de)
		if err != nil {
			return nil, err
		}
		vdir.Entries = append(vdir.Entries, &virtual.DirEntry{
			Path:    de.Name(),
			ModTime: stat.ModTime(),
			Mode:    stat.Mode(),
			Sys:     stat.Sys(),
			Size:    stat.Size(),
		})
	}
	return vdir, nil
}

func (f *FileSystem) toEntryInfo(fullpath string, de fs.DirEntry) (fs.FileInfo, error) {
	entryPath := path.Join(fullpath, de.Name())
	f.log.Debug("budfs: entry info", "name", entryPath)
	entry, ok := f.cache.Get(entryPath)
	if ok {
		f.log.Debug("budfs: cache hit", "name", entryPath)
		return entry.Open().Stat()
	}
	f.log.Debug("budfs: cache miss", "name", entryPath)
	return de.Info()
}
