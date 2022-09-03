package budfs

import (
	"context"
	"io/fs"
	"net/http"
	"path"
	"time"

	"github.com/livebud/bud/internal/virtual/vcache"

	"github.com/livebud/bud/package/vfs"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/package/budfs/genfs"
	"github.com/livebud/bud/package/budfs/mergefs"
	"github.com/livebud/bud/package/log"
)

func New(fsys vfs.ReadWritable, log log.Interface) *FileSystem {
	gen := genfs.New()
	merged := mergefs.Merge(gen, fsys)
	server := http.FileServer(http.FS(merged))
	return &FileSystem{
		cache:  vcache.New(),
		closer: new(once.Closer),
		dag:    dag.New(),
		gen:    gen,
		fsys:   fsys,
		log:    log,
		syncfs: merged,
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

type File = genfs.File

type Context struct {
	closer *once.Closer
	ctx    context.Context
	dag    *dag.Graph
	fsys   fs.FS
	path   string
}

var _ context.Context = (*Context)(nil)
var _ fs.FS = (*Context)(nil)
var _ fs.GlobFS = (*Context)(nil)

// Implement context.Context
func (c *Context) Deadline() (deadline time.Time, ok bool) { return c.ctx.Deadline() }
func (c *Context) Done() <-chan struct{}                   { return c.ctx.Done() }
func (c *Context) Err() error                              { return c.ctx.Err() }
func (c *Context) Value(key interface{}) interface{}       { return c.ctx.Value(key) }

// Open implement fs.FS
func (c *Context) Open(name string) (fs.File, error) {
	file, err := c.fsys.Open(name)
	if err != nil {
		return nil, err
	}
	c.dag.Link(c.path, name)
	return file, nil
}

// Glob implements fs.GlobFS
func (c *Context) Glob(name string) ([]string, error) {
	matches, err := fs.Glob(c.fsys, name)
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		c.dag.Link(c.path, match)
	}
	return matches, nil
}

func (c *Context) Link(from, to string) {
	c.dag.Link(from, to)
}

func (c *Context) Defer(fn func() error) {
	c.closer.Closes = append(c.closer.Closes, fn)
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

func (d *Dir) GenerateFile(path string, fn func(ctx *Context, file *File) error) {
	ctx := &Context{d.closer, context.TODO(), d.dag, d.fsys, path}
	d.dir.GenerateFile(path, func(file *genfs.File) error { return fn(ctx, file) })
}

func (d *Dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *Dir) GenerateDir(path string, fn func(ctx *Context, dir *Dir) error) {
	ctx := &Context{d.closer, context.TODO(), d.dag, d.fsys, path}
	d.dir.GenerateDir(path, func(dir *genfs.Dir) error {
		return fn(ctx, &Dir{d.closer, d.dag, d.fsys, dir})
	})
}

func (d *Dir) DirGenerator(path string, generator DirGenerator) {
	d.GenerateDir(path, generator.GenerateDir)
}

type FileGenerator interface {
	GenerateFile(ctx *Context, file *File) error
}

type DirGenerator interface {
	GenerateDir(ctx *Context, dir *Dir) error
}

type EmbedFile genfs.EmbedFile

var _ FileGenerator = (*EmbedFile)(nil)

func (e *EmbedFile) GenerateFile(ctx *Context, file *File) error {
	file.Data = e.Data
	return nil
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

func (f *FileSystem) GenerateFile(path string, fn func(ctx *Context, file *File) error) {
	ctx := &Context{f.closer, context.TODO(), f.dag, f.syncfs, path}
	f.gen.GenerateFile(path, func(file *genfs.File) error { return fn(ctx, file) })
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

func (f *FileSystem) GenerateDir(path string, fn func(ctx *Context, dir *Dir) error) {
	ctx := &Context{f.closer, context.TODO(), f.dag, f.syncfs, path}
	f.gen.GenerateDir(path, func(dir *genfs.Dir) error {
		return fn(ctx, &Dir{f.closer, f.dag, f.syncfs, dir})
	})
}

func (f *FileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

func (f *FileSystem) ServeFile(path string, fn func(ctx *Context, file *File) error) {
	ctx := &Context{f.closer, context.TODO(), f.dag, f.syncfs, path}
	f.gen.ServeFile(path, func(file *genfs.File) error { return fn(ctx, file) })
}

func (f *FileSystem) FileServer(path string, generator FileGenerator) {
	f.ServeFile(path, generator.GenerateFile)
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
