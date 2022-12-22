package genfs

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/virtual"
)

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
	Link(from, to string)
	Check(from string, checker func(path string) (changed bool))
}

type FS interface {
	fs.FS
	fs.ReadDirFS
	fs.GlobFS
	Watch(patterns ...string) error
}

type generator interface {
	Generate(target string) (fs.File, error)
}

func New(cache Cache, fsys fs.FS, log log.Log) *FileSystem {
	return &FileSystem{cache, fsys, log, newTree()}
}

type FileSystem struct {
	cache Cache   // File cache that supports linking files together into a DAG
	fsys  fs.FS   // Merged external filesystem (local, remote, etc.) with filler
	log   log.Log // Log messages
	tree  *tree   // Tree for the generators and filler nodes
}

var _ Generators = (*FileSystem)(nil)

func (f *FileSystem) GenerateFile(path string, fn func(fsys FS, file *File) error) {
	fileg := &fileGenerator{f.cache, fn, f, path}
	f.tree.Insert(path, modeGen, fileg)
}

func (f *FileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

func (f *FileSystem) GenerateDir(path string, fn func(fsys FS, dir *Dir) error) {
	dirg := &dirGenerator{f.cache, fn, f, path, f.tree}
	f.tree.Insert(path, modeDir|modeGen, dirg)
}

func (f *FileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

func (f *FileSystem) ServeFile(dir string, fn func(fsys FS, file *File) error) {
	server := &fileServer{f.cache, fn, f, dir}
	f.tree.Insert(dir, modeDir|modeGen, server)
}

func (f *FileSystem) FileServer(dir string, server FileServer) {
	f.ServeFile(dir, server.ServeFile)
}

func (f *FileSystem) Open(target string) (fs.File, error) {
	// Check that target is valid
	if !fs.ValidPath(target) {
		return nil, formatError(fs.ErrInvalid, "invalid target path %q", target)
	}
	return f.openFrom("", target)
}

func (f *FileSystem) openFrom(previous string, target string) (fs.File, error) {
	// First look for an exact matching generator
	node, found := f.tree.Find(target)
	if found && node.Generator != nil {
		file, err := node.Generator.Generate(target)
		if err != nil {
			return nil, formatError(err, "open %q", target)
		}
		return wrapFile(file, f, target), nil
	}
	// Next try opening the file from the fallback filesystem
	if file, err := f.fsys.Open(target); nil == err {
		return wrapFile(file, f, target), nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, formatError(err, "open %q", target)
	}
	// Next, if we did find a generator node above, return it now. It'll be a
	// filler directory, not a generator.
	if found && node.Mode.IsDir() {
		dir := virtual.New(&virtual.Dir{
			Path: target,
			Mode: node.Mode.FileMode(),
		})
		return wrapFile(dir, f, target), nil
	}
	// Lastly, try finding a node by its prefix
	node, found = f.tree.FindPrefix(target)
	if found && node.Path != previous && node.Mode.IsDir() && node.Generator != nil {
		if file, err := node.Generator.Generate(target); nil == err {
			return wrapFile(file, f, target), nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, formatError(err, "open by prefix %q", target)
		}
	}
	// Return a file not found error
	return nil, formatError(fs.ErrNotExist, "open %q", target)
}

func (f *FileSystem) ReadDir(target string) ([]fs.DirEntry, error) {
	deset := newDirEntrySet()
	node, ok := f.tree.Find(target)
	if ok {
		if !node.Mode.IsDir() {
			return nil, formatError(errNotImplemented, "tree readdir %q", target)
		}
		// Run the directory generator
		if node.Mode.IsGen() {
			if _, err := node.Generator.Generate(target); err != nil {
				return nil, err
			}
		}
		for _, child := range node.children {
			deset.Add(newDirEntry(f, child.Name, child.Mode.FileMode(), child.Path))
		}
	}
	des, err := fs.ReadDir(f.fsys, target)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, formatError(err, "fallback readdir %q", target)
	}
	for _, de := range des {
		deset.Add(de)
	}
	return deset.List(), nil
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
