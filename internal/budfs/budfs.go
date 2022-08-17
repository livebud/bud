package budfs

import (
	"io/fs"
	"net/http"

	"github.com/livebud/bud/internal/cachefs"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/genfs"
	"github.com/livebud/bud/internal/mergefs"
	"github.com/livebud/bud/package/log"
)

func New(fsys fs.FS, log log.Interface) *FileSystem {
	gen := genfs.New()
	cache := cachefs.New(log)
	fsys = cache.Wrap(mergefs.New(fsys, gen))
	return &FileSystem{
		dag:  dag.New(),
		gen:  gen,
		fsys: fsys,
		log:  log,
	}
}

type FileSystem struct {
	cache *cachefs.Cache
	dag   *dag.Graph
	gen   *genfs.FileSystem
	fsys  fs.FS
	log   log.Interface
}

type FS interface {
	fs.FS
	Link(from, to string)
}

type File = genfs.File
type Dir = genfs.Dir

type FileGenerator interface {
	GenerateFile(fsys FS, file *File) error
}

type DirGenerator interface {
	GenerateDir(fsys FS, dir *Dir) error
}

func (f *FileSystem) Open(name string) (fs.File, error) {
	file, err := f.fsys.Open(name)
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
	return f.fsys.Open(name)
}

func (f *linkedFS) Link(from, to string) {
	// f.fsys.Link(from, to)
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	fsys := &linkedFS{f.dag, f.fsys, path}
	f.gen.GenerateFile(path, func(file *genfs.File) error {
		return generator.GenerateFile(fsys, file)
	})
}

func (f *FileSystem) DirGenerator(path string, generator DirGenerator) {
	fsys := &linkedFS{f.dag, f.fsys, path}
	f.gen.GenerateDir(path, func(dir *genfs.Dir) error {
		return generator.GenerateDir(fsys, dir)
	})
}

func (f *FileSystem) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
