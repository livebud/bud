package budfs

import (
	"io/fs"
	"net/http"
	"path"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/package/budfs/cachefs"
	"github.com/livebud/bud/package/budfs/genfs"
	"github.com/livebud/bud/package/budfs/mergefs"
	"github.com/livebud/bud/package/log"
)

func New(fsys fs.FS, log log.Interface) *FileSystem {
	gen := genfs.New()
	cache := cachefs.New(log)
	merged := mergefs.New(gen, fsys)
	syncfs := cache.Wrap(merged)
	server := http.FileServer(http.FS(merged))
	return &FileSystem{
		cache:  cache,
		dag:    dag.New(),
		gen:    gen,
		log:    log,
		syncfs: syncfs,
		server: server,
	}
}

type FileSystem struct {
	cache  *cachefs.Cache
	dag    *dag.Graph
	gen    *genfs.FileSystem
	log    log.Interface
	syncfs fs.FS
	server http.Handler
}

type FS interface {
	fs.FS
	Link(from, to string)
}

type File = genfs.File

type Dir struct {
	dag  *dag.Graph
	fsys fs.FS
	dir  *genfs.Dir
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

func (d *Dir) GenerateFile(path string, fn func(fsys FS, file *File) error) {
	fsys := &linkedFS{d.dag, d.fsys, path}
	d.dir.GenerateFile(path, func(file *genfs.File) error {
		return fn(fsys, file)
	})
}

func (d *Dir) GenerateDir(path string, fn func(fsys FS, dir *Dir) error) {
	fsys := &linkedFS{d.dag, d.fsys, path}
	d.dir.GenerateDir(path, func(dir *genfs.Dir) error {
		return fn(fsys, &Dir{d.dag, d.fsys, dir})
	})
}

type FileGenerator interface {
	GenerateFile(fsys FS, file *File) error
}

type DirGenerator interface {
	GenerateDir(fsys FS, dir *Dir) error
}

func (f *FileSystem) Open(name string) (fs.File, error) {
	file, err := f.syncfs.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

type linkedFS struct {
	dag  *dag.Graph
	fsys fs.FS
	path string
}

func (f *linkedFS) Open(name string) (fs.File, error) {
	f.dag.Link(f.path, name)
	return f.fsys.Open(name)
}

// Link allows for arbitrary links
func (f *linkedFS) Link(from, to string) {
	f.dag.Link(from, to)
}

type GenerateFile func(fsys FS, file *File) error

func (fn GenerateFile) GenerateFile(fsys FS, file *File) error {
	return fn(fsys, file)
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	fsys := &linkedFS{f.dag, f.syncfs, path}
	f.gen.GenerateFile(path, func(file *genfs.File) error {
		return generator.GenerateFile(fsys, file)
	})
}

type GenerateDir func(fsys FS, dir *Dir) error

func (fn GenerateDir) GenerateDir(fsys FS, dir *Dir) error {
	return fn(fsys, dir)
}

func (f *FileSystem) DirGenerator(path string, generator DirGenerator) {
	fsys := &linkedFS{f.dag, f.syncfs, path}
	f.gen.GenerateDir(path, func(dir *genfs.Dir) error {
		return generator.GenerateDir(fsys, &Dir{f.dag, f.syncfs, dir})
	})
}

// ServeHTTP serves the filesystem. Served files are not cached.
func (f *FileSystem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.server.ServeHTTP(w, r)
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
