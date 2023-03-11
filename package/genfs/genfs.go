package genfs

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/virtual"
)

// Directory is the interface for adding file generators and directories
type Directory interface {
	GenerateFile(path string, fn func(fsys FS, file *File) error)
	FileGenerator(path string, generator FileGenerator)
	GenerateDir(path string, fn func(fsys FS, dir *Dir) error)
	DirGenerator(path string, generator DirGenerator)
	ServeFile(dir string, fn func(fsys FS, file *File) error)
	FileServer(dir string, server FileServer)
	GenerateExternal(path string, fn func(fsys FS, file *External) error)
	ExternalGenerator(path string, generator ExternalGenerator)
}

// // Extension allows you to extend the filesystem with more generators
// type Extension interface {
// 	Extend(g Directory)
// }

// FileSystem allows you to create and use a generator filesystem
type FileSystem interface {
	fs.FS
	fs.ReadDirFS
	Directory
	// Sub(dir string) FileSystem
	// Extend(extensions ...Extension) FileSystem
	// Sync(to virtual.FS, subdirs ...string) error
	// Copy(to virtual.FS, subdirs ...string) error
}

type Cache interface {
	Get(path string) (*virtual.File, error)
	Set(path string, file *virtual.File) error
	Link(from string, toPatterns ...string) error
}

type discardCache struct{}

func (discardCache) Get(path string) (*virtual.File, error)       { return nil, errors.New("not found") }
func (discardCache) Set(path string, file *virtual.File) error    { return nil }
func (discardCache) Link(from string, toPatterns ...string) error { return nil }

type FS interface {
	fs.FS
	fs.ReadDirFS
	fs.GlobFS
	Watch(patterns ...string) error
}

type generator interface {
	Generate(target string) (fs.File, error)
}

func New(cache Cache, fsys fs.FS, log log.Log) FileSystem {
	if cache == nil {
		cache = discardCache{}
	}
	return &fileSystem{cache, fsys, log, newTree()}
}

type fileSystem struct {
	cache Cache   // File cache that supports linking files together into a DAG
	fsys  fs.FS   // Merged external filesystem (local, remote, etc.) with filler
	log   log.Log // Log messages
	tree  *tree   // Tree for the generators and filler nodes
}

var _ FileSystem = (*fileSystem)(nil)

// // Extend the filesystem with more generators
// func (f *fileSystem) Extend(extensions ...Extension) FileSystem {
// 	for _, extension := range extensions {
// 		extension.Extend(f)
// 	}
// 	return f
// }

func (f *fileSystem) GenerateFile(path string, fn func(fsys FS, file *File) error) {
	fileg := &fileGenerator{f.cache, fn, f, path}
	f.tree.Insert(path, modeGen, fileg)
}

func (f *fileSystem) FileGenerator(path string, generator FileGenerator) {
	f.GenerateFile(path, generator.GenerateFile)
}

func (f *fileSystem) GenerateDir(path string, fn func(fsys FS, dir *Dir) error) {
	dirg := &dirGenerator{f.cache, fn, f, path, f.tree}
	f.tree.Insert(path, modeGenDir, dirg)
}

func (f *fileSystem) DirGenerator(path string, generator DirGenerator) {
	f.GenerateDir(path, generator.GenerateDir)
}

func (f *fileSystem) ServeFile(dir string, fn func(fsys FS, file *File) error) {
	server := &fileServer{f.cache, fn, f, dir}
	f.tree.Insert(dir, modeGenDir, server)
}

func (f *fileSystem) FileServer(dir string, server FileServer) {
	f.ServeFile(dir, server.ServeFile)
}

func (f *fileSystem) GenerateExternal(path string, fn func(fsys FS, file *External) error) {
	fileg := &externalGenerator{f.cache, fn, f, path}
	f.tree.Insert(path, modeGen, fileg)
}
func (f *fileSystem) ExternalGenerator(path string, generator ExternalGenerator) {
	f.GenerateExternal(path, generator.GenerateExternal)
}

func (f *fileSystem) Open(target string) (fs.File, error) {
	// Check that target is valid
	if !fs.ValidPath(target) {
		return nil, formatError(fs.ErrInvalid, "invalid target path %q", target)
	}
	return f.openFrom("", target)
}

func (f *fileSystem) openFrom(previous string, target string) (fs.File, error) {
	// First look for an exact matching generator
	node, found := f.tree.Find(target)
	if found && node.Generator != nil {
		file, err := node.Generator.Generate(target)
		if err != nil {
			return nil, formatError(err, "open %q", target)
		}
		return wrapFile(file, f, node.Path), nil
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
		dir := virtual.Open(&virtual.File{
			Path: target,
			Mode: node.Mode.FileMode(),
		})
		return wrapFile(dir, f, node.Path), nil
	}
	// Lastly, try finding a node by its prefix
	node, found = f.tree.FindPrefix(target)
	if found && node.Path != previous && node.Mode.IsDir() && node.Generator != nil {
		if file, err := node.Generator.Generate(target); nil == err {
			return wrapFile(file, f, node.Path), nil
		} else if !errors.Is(err, fs.ErrNotExist) {
			return nil, formatError(err, "open by prefix %q", target)
		}
	}
	// Return a file not found error
	return nil, formatError(fs.ErrNotExist, "open %q", target)
}

func (f *fileSystem) ReadDir(target string) ([]fs.DirEntry, error) {
	deset := newDirEntrySet()
	node, ok := f.tree.Find(target)
	if ok {
		if !node.Mode.IsDir() {
			return nil, formatError(errNotImplemented, "tree readdir %q", target)
		}
		// Run the directory generator
		if node.Mode.IsGen() {
			// Generate is expected to update the tree, that's why we don't use the
			// returned file
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

// func (f *fileSystem) Sync(to virtual.FS, subdirs ...string) error {
// 	return virtual.Sync(f.log, f, to, subdirs...)
// }

// func (f *fileSystem) Copy(to virtual.FS, subdirs ...string) error {
// 	return virtual.Copy(f.log, f, to, subdirs...)
// }

// func (f *fileSystem) Sub(dir string) FileSystem {
// 	return &subSystem{dir, f, f.log}
// }

// type subSystem struct {
// 	dir  string
// 	fsys FileSystem
// 	log  log.Log
// }

// var _ FileSystem = (*subSystem)(nil)

// func (s *subSystem) Open(name string) (fs.File, error) {
// 	return s.fsys.Open(path.Join(s.dir, name))
// }

// func (s *subSystem) ReadDir(name string) ([]fs.DirEntry, error) {
// 	return s.fsys.ReadDir(path.Join(s.dir, name))
// }

// func (s *subSystem) GenerateFile(subpath string, fn func(fsys FS, file *File) error) {
// 	s.fsys.GenerateFile(path.Join(s.dir, subpath), fn)
// }

// func (s *subSystem) FileGenerator(subpath string, generator FileGenerator) {
// 	s.fsys.FileGenerator(path.Join(s.dir, subpath), generator)
// }

// func (s *subSystem) GenerateDir(subdir string, fn func(fsys FS, dir *Dir) error) {
// 	s.fsys.GenerateDir(path.Join(s.dir, subdir), fn)
// }

// func (s *subSystem) DirGenerator(subdir string, generator DirGenerator) {
// 	s.fsys.DirGenerator(path.Join(s.dir, subdir), generator)
// }

// func (s *subSystem) ServeFile(subdir string, fn func(fsys FS, file *File) error) {
// 	s.fsys.ServeFile(path.Join(s.dir, subdir), fn)
// }

// func (s *subSystem) FileServer(subdir string, server FileServer) {
// 	s.fsys.FileServer(path.Join(s.dir, subdir), server)
// }

// func (s *subSystem) GenerateExternal(subpath string, fn func(fsys FS, file *External) error) {
// 	s.fsys.GenerateExternal(path.Join(s.dir, subpath), fn)
// }

// func (s *subSystem) ExternalGenerator(subpath string, generator ExternalGenerator) {
// 	s.fsys.ExternalGenerator(path.Join(s.dir, subpath), generator)
// }

// func (s *subSystem) Sub(dir string) FileSystem {
// 	return &subSystem{path.Join(s.dir, dir), s.fsys, s.log}
// }

// func (s *subSystem) Extend(extensions ...Extension) FileSystem {
// 	return s.fsys.Extend(extensions...)
// }

// func (s *subSystem) Sync(to virtual.FS, subdirs ...string) error {
// 	return virtual.Sync(s.log, s, to, subdirs...)
// }

// func (s *subSystem) Copy(to virtual.FS, subdirs ...string) error {
// 	return virtual.Copy(s.log, s, to, subdirs...)
// }

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
