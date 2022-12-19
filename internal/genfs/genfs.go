package genfs

import (
	"errors"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	"github.com/livebud/bud/internal/treefs"
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

type generator interface {
	Generate(target string) (fs.File, error)
}

// Modes for the generators
const (
	modeGenerator    = treefs.ModeGenerator
	modeGeneratorDir = treefs.ModeDir | treefs.ModeGenerator
)

func New(cache Cache, fsys fs.FS, linker Linker, log log.Log) FileSystem {
	tree := treefs.New(".")
	return &fileSystem{cache, fsys, linker, log, tree}
}

type fileSystem struct {
	cache  Cache        // File cache
	fsys   fs.FS        // Fallback filesystem
	linker Linker       // Link files into a DAG
	log    log.Log      // Log messages
	tree   *treefs.Tree // Filesystem tree containing generators and filler dirs
}

var _ FileSystem = (*fileSystem)(nil)

func (f *fileSystem) GenerateFile(path string, fn func(fsys FS, file *File) error) {
	fileg := &fileGenerator{f.cache, fn, f, f.linker, path}
	f.tree.Insert(path, modeGenerator, fileg)
}

func (f *fileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

func (f *fileSystem) GenerateDir(path string, fn func(fsys FS, dir *Dir) error) {
	dirg := &dirGenerator{f.cache, fn, f, f.linker, path, f.tree}
	f.tree.Insert(path, modeGeneratorDir, dirg)
}

func (f *fileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

func (f *fileSystem) ServeFile(dir string, fn func(fsys FS, file *File) error) {
	server := &fileServer{f.cache, fn, f, f.linker, dir}
	f.tree.Insert(dir, modeGeneratorDir, server)
}

func (f *fileSystem) FileServer(dir string, server FileServer) {
	f.ServeFile(dir, server.ServeFile)
}

func (f *fileSystem) Open(target string) (fs.File, error) {
	return f.open(".", target)
}

func (f *fileSystem) open(previous, target string) (fs.File, error) {
	// Check that target is valid
	if !fs.ValidPath(target) {
		return nil, formatError(fs.ErrInvalid, "genfs: open invalid path %q", target)
	}
	// 1. Look for an exact matching generator
	node, ok := f.tree.Get(target)
	if ok && node.Mode().IsGenerator() {
		file, err := node.Generate(target)
		if err != nil {
			return nil, formatError(err, "genfs: open generator %q", target)
		}
		return &wrapFile{file, f, target}, nil
	}
	// 2. Try opening the file from the fallback filesystem
	if file, err := f.fsys.Open(target); nil == err {
		return &wrapFile{file, f, target}, nil
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, formatError(err, "genfs: open filesystem %q", target)
	}
	// 3. If we did have a filler node above, return it now
	if ok {
		if file, err := node.Generate(target); nil == err {
			return file, nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, formatError(err, "genfs: open filler %q", target)
		}
	}
	// 4. Try finding a generator by its prefix
	node, ok = f.tree.FindByPrefix(target)
	if ok && isDifferentDirGenerator(previous, node) {
		if file, err := node.Generate(target); nil == err {
			return &wrapFile{file, f, target}, nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, formatError(err, "genfs: open by prefix %q", target)
		}
	}
	// Return a file not found error
	return nil, formatError(fs.ErrNotExist, "genfs: open %q", target)
}

// ReadDir implements fs.ReadDirFS
func (f *fileSystem) ReadDir(target string) (entries []fs.DirEntry, err error) {
	if !fs.ValidPath(target) {
		return nil, formatError(fs.ErrInvalid, "genfs: readdir invalid path %q", target)
	}

	// 1. Look for an exact matching directory generator
	node, ok := f.tree.Get(target)
	if ok && node.Mode().IsDir() {
		if file, err := node.Generate(target); nil == err {
			// wfile := &wrapFile{file, f, target}
			// dir, ok := wfile.ReadDir(-1)
			// if !ok {
			// 	return nil, formatError(fs.ErrInvalid, "genfs: expected dir to implement fs.ReadDirFile %q", target)
			// }
			wfile := &wrapFile{file, f, target}
			if des, err := wfile.ReadDir(-1); nil == err {
				entries = append(entries, des...)
			} else if !errors.Is(err, fs.ErrNotExist) {
				return nil, formatError(err, "genfs: readdir filesystem %q", target)
			}
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, formatError(err, "genfs: readdir generator %q", target)
		}
	}

	// 2. Try reading the directory from the fallback filesystem
	if des, err := fs.ReadDir(f.fsys, target); nil == err {
		entries = append(entries, des...)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, formatError(err, "genfs: readdir filesystem %q", target)
	}

	// Dedupe the set
	entries = dirEntrySet(entries)

	// Sort by name
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	return entries, nil
}

// Ensure the  node is different from the last one and is a directory generator.
func isDifferentDirGenerator(prev string, node treefs.Node) bool {
	return prev != node.Path() && node.Mode().IsDir() && node.Mode().IsGenerator()
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

func dirEntrySet(list []fs.DirEntry) []fs.DirEntry {
	seen := map[string]struct{}{}
	set := make([]fs.DirEntry, 0, len(list))
	for _, item := range list {
		if _, ok := seen[item.Name()]; ok {
			continue
		}
		seen[item.Name()] = struct{}{}
		set = append(set, item)
	}
	return set
}
