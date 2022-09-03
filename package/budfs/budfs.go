package budfs

import (
	"context"
	"io/fs"
	"net/http"
	"path"

	"github.com/livebud/bud/package/vfs"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/virtual/vcache"
	"github.com/livebud/bud/package/budfs/genfs"
	"github.com/livebud/bud/package/budfs/mergefs"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/virtual"
)

func New(cache vcache.Cache, fsys vfs.ReadWritable, log log.Interface) *FileSystem {
	gen := genfs.New(cache)
	merged := mergefs.Merge(gen, fsys)
	server := http.FileServer(http.FS(merged))
	syncfs := wrapFS(cache, merged, log)
	return &FileSystem{
		cache:  cache,
		closer: new(once.Closer),
		dag:    dag.New(),
		gen:    gen,
		fsys:   fsys,
		log:    log,
		syncfs: syncfs,
		server: server,
	}
}

type FileSystem struct {
	cache  vcache.Cache
	closer *once.Closer
	dag    *dag.Graph
	gen    *genfs.FileSystem
	fsys   vfs.ReadWritable
	log    log.Interface
	syncfs fs.FS
	server http.Handler
}

type FS interface {
	fs.FS
	Link(from, to string)
	Defer(fn func() error)
}

type File = genfs.File

type EmbedFile genfs.EmbedFile

var _ FileGenerator = (*EmbedFile)(nil)

func (e *EmbedFile) GenerateFile(ctx context.Context, fsys FS, file *File) error {
	file.Data = e.Data
	return nil
}

type Dir struct {
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

func (d *Dir) GenerateFile(path string, fn func(ctx context.Context, fsys FS, file *File) error) {
	fsys := &linkedFS{d.closer, d.dag, d.fsys, path}
	d.dir.GenerateFile(path, func(file *genfs.File) error {
		return fn(context.TODO(), fsys, file)
	})
}

func (d *Dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *Dir) GenerateDir(path string, fn func(ctx context.Context, fsys FS, dir *Dir) error) {
	fsys := &linkedFS{d.closer, d.dag, d.fsys, path}
	d.dir.GenerateDir(path, func(dir *genfs.Dir) error {
		return fn(context.TODO(), fsys, &Dir{d.closer, d.dag, d.fsys, dir})
	})
}

func (d *Dir) DirGenerator(path string, generator DirGenerator) {
	d.GenerateDir(path, generator.GenerateDir)
}

type FileGenerator interface {
	GenerateFile(ctx context.Context, fsys FS, file *File) error
}

type DirGenerator interface {
	GenerateDir(ctx context.Context, fsys FS, dir *Dir) error
}

func (f *FileSystem) Open(name string) (fs.File, error) {
	file, err := f.syncfs.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *FileSystem) Close() error {
	return f.closer.Close()
}

type linkedFS struct {
	closer *once.Closer
	dag    *dag.Graph
	fsys   fs.FS
	path   string
}

var _ fs.FS = (*linkedFS)(nil)
var _ fs.GlobFS = (*linkedFS)(nil)

func (f *linkedFS) Open(name string) (fs.File, error) {
	file, err := f.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	f.dag.Link(f.path, name)
	return file, nil
}

func (f *linkedFS) Glob(name string) ([]string, error) {
	matches, err := fs.Glob(f.fsys, name)
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		f.dag.Link(f.path, match)
	}
	return matches, nil
}

// Link allows for arbitrary links
func (f *linkedFS) Link(from, to string) {
	f.dag.Link(from, to)
}

func (f *linkedFS) Defer(fn func() error) {
	f.closer.Closes = append(f.closer.Closes, fn)
}

func (f *FileSystem) GenerateFile(path string, fn func(ctx context.Context, fsys FS, file *File) error) {
	fsys := &linkedFS{f.closer, f.dag, f.syncfs, path}
	f.gen.GenerateFile(path, func(file *genfs.File) error {
		return fn(context.TODO(), fsys, file)
	})
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

func (f *FileSystem) GenerateDir(path string, fn func(ctx context.Context, fsys FS, dir *Dir) error) {
	fsys := &linkedFS{f.closer, f.dag, f.syncfs, path}
	f.gen.GenerateDir(path, func(dir *genfs.Dir) error {
		return fn(context.TODO(), fsys, &Dir{f.closer, f.dag, f.syncfs, dir})
	})
}

func (f *FileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

// ServeHTTP serves the filesystem. Served files are not cached.
func (f *FileSystem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.server.ServeHTTP(w, r)
}

// Sync the overlay to the filesystem
func (f *FileSystem) Sync(to string) error {
	// Clear the filesystem cache before syncing again
	f.cache.Clear()
	return dsync.To(f.syncfs, f.fsys, to)
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
