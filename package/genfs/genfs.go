package genfs

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/livebud/bud/package/budfs/mergefs"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/virtual"
)

type FileSystem interface {
	fs.FS
	Generators
}

type Generators interface {
	GenerateFile(path string, fn func(fsys FS, file *File) error)
	FileGenerator(path string, generator FileGenerator)
	GenerateDir(path string, fn func(fsys FS, dir *Dir) error)
	DirGenerator(path string, generator DirGenerator)
	ServeFile(dir string, fn func(fsys FS, file *File) error)
	FileServer(dir string, server FileServer)
}

type Cache interface {
	Get(name string) (entry virtual.Entry, ok bool)
	Set(path string, entry virtual.Entry)
}

type Linker interface {
	Link(fromPattern string, toPatterns ...string) error
}

type FS interface {
	fs.FS
	fs.ReadDirFS
	fs.GlobFS
	Watch(patterns ...string) error
}

func New(cache Cache, fsys fs.FS, linker Linker, log log.Log) FileSystem {
	filler := newFiller()
	fsys = mergefs.Merge(fsys, filler)
	return &fileSystem{cache, fsys, linker, log, newRadix(), filler}
}

type fileSystem struct {
	cache   Cache   // File cache
	mergefs fs.FS   // Merged external filesystem (local, remote, etc.) with filler
	linker  Linker  // Link files into a DAG
	log     log.Log // Log messages
	radix   *radix  // Radix tree for matching generators
	filler  *filler // Fill in missing files and dirs between generators
}

var _ FileSystem = (*fileSystem)(nil)

func (f *fileSystem) GenerateFile(path string, fn func(fsys FS, file *File) error) {
	fileg := &fileGenerator{f.cache, fn, f, f.linker, path}
	f.radix.Insert(path, fileg)
	f.filler.Insert(path, fs.FileMode(0))
}

func (f *fileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

func (f *fileSystem) GenerateDir(path string, fn func(fsys FS, dir *Dir) error) {
	dirg := &dirGenerator{f.cache, fn, f, f.linker, path, f.radix, f.filler}
	f.radix.Insert(path, dirg)
	f.filler.Insert(path, fs.ModeDir)
}

func (f *fileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

func (f *fileSystem) ServeFile(dir string, fn func(fsys FS, file *File) error) {
	server := &fileServer{f.cache, fn, f, f.linker, dir}
	f.radix.Insert(dir, server)
	f.filler.Insert(dir, fs.ModeDir)
}

func (f *fileSystem) FileServer(dir string, server FileServer) {
	f.ServeFile(dir, server.ServeFile)
}

func (f *fileSystem) Open(target string) (fs.File, error) {
	// Check that target is valid
	if !fs.ValidPath(target) {
		return nil, formatError(fs.ErrInvalid, "invalid target path %q", target)
	}
	return f.openAs("", target)
}

func (f *fileSystem) openAs(callerPath string, target string) (fs.File, error) {
	if callerPath == target {
		return nil, formatError(fs.ErrInvalid, "genfs: cycle detected %q", target)
	}
	generator, ok := f.radix.Get(target)
	if ok {
		file, err := generator.Generate(target)
		if err != nil {
			return nil, formatError(err, "genfs: open %q", target)
		}
		return file, nil
	}
	// Try the merged filesystem
	if file, err := f.mergefs.Open(target); nil == err {
		return &wrapFile{file, f, target}, nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, formatError(err, "genfs: open %q", target)
	}
	// Lastly, try finding a generator by its prefix
	generator, prefix, ok := f.radix.FindByPrefix(target)
	if !ok {
		// We didn't find a generator with that prefix generator
		return nil, formatError(fs.ErrNotExist, "genfs: open %q", target)
	} else if prefix == callerPath {
		// Generator isn't making progress and we're stuck in a loop. This occurs
		// when we're trying to opening a file that matches a directory, but that
		// directory doesn't have the file.
		return nil, formatError(fs.ErrNotExist, "genfs: open %q", target)
	}
	file, err := generator.Generate(target)
	if err != nil {
		return nil, formatError(err, "genfs: open %q", target)
	}
	return file, nil
}

func relativePath(base, target string) string {
	rel := strings.TrimPrefix(target, base)
	if rel == "" {
		return "."
	} else if rel[0] == '/' {
		rel = rel[1:]
	}
	return rel
}

func formatError(err error, format string, args ...interface{}) error {
	return fmt.Errorf("genfs: %s. %w", fmt.Sprintf(format, args...), err)
}
