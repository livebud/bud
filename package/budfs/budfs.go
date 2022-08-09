package budfs

import (
	"fmt"
	"io/fs"

	"github.com/livebud/bud/package/merged"

	"github.com/livebud/bud/internal/virtual"
	"github.com/livebud/bud/package/log"
)

type Generator interface {
	Generate(target string) (fs.File, error)
}

type FS interface {
	fs.FS
	Link(from, to string)
}

func New(fsys fs.FS, log log.Interface) *FileSystem {
	filler := newFiller()
	return &FileSystem{
		// Merge the passed in filesystem with the filler filesystem
		fsys:   merged.Merge(fsys, filler),
		log:    log,
		radix:  newRadix(),
		filler: filler,
	}
}

type FileSystem struct {
	fsys   fs.FS
	log    log.Interface
	radix  *radix
	filler *filler
}

func (f *FileSystem) Open(name string) (fs.File, error) {
	// Get an exact match
	if generator, ok := f.radix.Get(name); ok {
		return generator.Generate(name)
	}
	// Get a generator that's a prefix of the name
	if _, generator, ok := f.radix.GetByPrefix(name); ok {
		return generator.Generate(name)
	}
	// Try opening the underlying filesystem
	return f.fsys.Open(name)
}

type File struct {
	path string
	Data []byte
	Mode fs.FileMode
}

func (f *File) Path() string {
	return f.path
}

type FileGenerator interface {
	GenerateFile(fsys FS, file *File) error
}

func (f *FileSystem) FileGenerator(filepath string, generator FileGenerator) {
	fileGenerator := &fileGenerator{f, filepath, generator}
	f.radix.Set(filepath, fileGenerator)
	f.filler.Add(filepath, fileGenerator)
}

type DirGenerator interface {
	GenerateDir(fsys FS, dir *Dir) error
}

func (f *FileSystem) DirGenerator(filepath string, generator DirGenerator) {
	f.radix.Set(filepath, &dirGenerator{f, filepath, generator})
}

type fileGenerator struct {
	fsys      *FileSystem
	path      string
	generator FileGenerator
}

var _ Generator = (*fileGenerator)(nil)

type scopedFS struct {
	fsys   *FileSystem
	target string
}

func (s *scopedFS) Link(from, to string) {
}

func (s *scopedFS) Open(name string) (fs.File, error) {
	return s.fsys.Open(name)
}

func (g *fileGenerator) Generate(target string) (fs.File, error) {
	// Prevents prefixes from matching files
	// (e.g. go.mod/go.mod from matching go.mod)
	if g.path != target {
		return nil, fs.ErrNotExist
	}
	file := &File{
		path: target,
	}
	// Scope FS for linking
	fsys := &scopedFS{
		fsys:   g.fsys,
		target: target,
	}
	if err := g.generator.GenerateFile(fsys, file); err != nil {
		return nil, fmt.Errorf("budfs: error generating file %q. %w", target, err)
	}
	// TODO: cache file based on scope key
	return &virtual.File{
		Name: file.path,
		Data: file.Data,
		Mode: file.Mode,
	}, nil
}

type Dir struct {
	path   string
	target string // equal to path unless called as a prefix
	// Entries []fs.Entry
}

func (d *Dir) Path() string {
	return d.path
}

func (d *Dir) Target() string {
	return d.target
}

func (d *Dir) FileGenerator(filepath string, generator FileGenerator) {
	fmt.Println(filepath)
}

type dirGenerator struct {
	fsys      *FileSystem
	path      string
	generator DirGenerator
}

func (g *dirGenerator) Generate(target string) (fs.File, error) {
	return nil, fmt.Errorf("not implemented yet")
}

// type dirEntry struct {
// 	path string
// }

// type dirFiller struct {
// 	path    string
// 	entries map[string]Generator
// }
