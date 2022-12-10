package genfs

import (
	"fmt"
	"io/fs"
	gopath "path"

	"github.com/livebud/bud/package/virtual"
)

type Dir struct {
	cache  Cache
	genfs  *FileSystem
	path   string  // Current directory path
	target string  // Final target path
	radix  *radix  // Radix tree for matching generators
	filler *filler // Fill in missing files and dirs between generators
}

var _ Generators = (*Dir)(nil)

func (d *Dir) Target() string {
	return d.target
}

func (d *Dir) Relative() string {
	return relativePath(d.path, d.target)
}

func (d *Dir) Path() string {
	return d.path
}

func (d *Dir) Mode() fs.FileMode {
	return fs.ModeDir
}

func (d *Dir) GenerateFile(path string, fn func(fsys FS, file *File) error) {
	fpath := gopath.Join(d.path, path)
	fileg := &fileGenerator{d.cache, fn, d.genfs, fpath}
	d.radix.Insert(fpath, fileg)
	d.filler.Insert(fpath, fs.FileMode(0))
}

func (d *Dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *Dir) GenerateDir(path string, fn func(fsys FS, dir *Dir) error) {
	fpath := gopath.Join(d.path, path)
	dirg := &dirGenerator{d.cache, fn, d.genfs, fpath, d.radix, d.filler}
	d.radix.Insert(fpath, dirg)
	d.filler.Insert(fpath, fs.ModeDir)
}

func (d *Dir) DirGenerator(path string, generator DirGenerator) {
	d.GenerateDir(path, generator.GenerateDir)
}

func (d *Dir) ServeFile(path string, fn func(fsys FS, file *File) error) {
	fpath := gopath.Join(d.path, path)
	server := &fileServer{d.cache, fn, d.genfs, fpath}
	d.radix.Insert(fpath, server)
	d.filler.Insert(fpath, fs.ModeDir)
}

func (d *Dir) FileServer(path string, server FileServer) {
	d.ServeFile(path, server.ServeFile)
}

func (d *Dir) GenerateExternal(path string, fn func(fsys FS, file *ExternalFile) error) {
	fpath := gopath.Join(d.path, path)
	external := &externalGenerator{d.cache, fn, d.genfs, fpath}
	d.radix.Insert(fpath, external)
	d.filler.Insert(fpath, fs.FileMode(0))
}

func (d *Dir) ExternalGenerator(path string, generator ExternalGenerator) {
	d.GenerateExternal(path, generator.GenerateExternal)
}

type DirGenerator interface {
	GenerateDir(fsys FS, dir *Dir) error
}

type GenerateDir func(fsys FS, dir *Dir) error

func (fn GenerateDir) GenerateDir(fsys FS, dir *Dir) error {
	return fn(fsys, dir)
}

type dirGenerator struct {
	cache  Cache
	fn     func(fsys FS, dir *Dir) error
	genfs  *FileSystem
	path   string
	radix  *radix  // Radix tree for matching generators
	filler *filler // Fill in missing files and dirs between generators
}

func (d *dirGenerator) Generate(target string) (fs.File, error) {
	if entry, ok := d.cache.Get(d.path); ok {
		_ = entry
		// TODO: wrap the entry file in a virtualDir
		return nil, fmt.Errorf("cache get not implemented yet")
	}
	scopedFS := &scopedFS{d.cache, d.genfs, d.path}
	dir := &Dir{d.cache, d.genfs, d.path, target, d.radix, d.filler}
	if err := d.fn(scopedFS, dir); err != nil {
		return nil, err
	}
	// Traverse into the directory looking for the target
	if d.path != target {
		return d.genfs.openAs(d.path, target)
	}
	des, err := fs.ReadDir(d.filler, target)
	if err != nil {
		return nil, err
	}
	// Create the virtual directory
	dirEntries := make([]fs.DirEntry, len(des))
	for i := range des {
		dirEntries[i] = &dirEntry{d.genfs, des[i].Name(), des[i].Type(), gopath.Join(d.path, des[i].Name())}
	}
	entry := &virtual.Dir{
		Path:    d.path,
		Mode:    fs.ModeDir,
		Entries: dirEntries,
	}
	// Cache the directory entry
	d.cache.Set(d.path, entry)
	// Return the virtual directory
	return &wrapFile{
		File:  virtual.New(entry),
		genfs: d.genfs,
		path:  d.path,
	}, nil
}
