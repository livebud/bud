package genfs

import (
	"fmt"
	"io/fs"
	gopath "path"

	"github.com/livebud/bud/internal/fsmode"
	"github.com/livebud/bud/package/virt"
)

type Dir struct {
	cache  Cache
	genfs  *FileSystem
	path   string // Current directory path
	target string // Final target path
	tree   *tree
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
	d.tree.Insert(fpath, fsmode.Gen, fileg)
}

func (d *Dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *Dir) GenerateDir(path string, fn func(fsys FS, dir *Dir) error) {
	fpath := gopath.Join(d.path, path)
	dirg := &dirGenerator{d.cache, fn, d.genfs, fpath, d.tree}
	d.tree.Insert(fpath, fsmode.GenDir, dirg)
}

func (d *Dir) DirGenerator(path string, generator DirGenerator) {
	d.GenerateDir(path, generator.GenerateDir)
}

func (d *Dir) ServeFile(path string, fn func(fsys FS, file *File) error) {
	fpath := gopath.Join(d.path, path)
	server := &fileServer{d.cache, fn, d.genfs, fpath}
	d.tree.Insert(fpath, fsmode.GenDir, server)
}

func (d *Dir) FileServer(path string, server FileServer) {
	d.ServeFile(path, server.ServeFile)
}

func (d *Dir) GenerateExternal(path string, fn func(fsys FS, file *External) error) {
	fpath := gopath.Join(d.path, path)
	fileg := &externalGenerator{d.cache, fn, d.genfs, fpath}
	d.tree.Insert(fpath, fsmode.Gen, fileg)
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
	cache Cache
	fn    func(fsys FS, dir *Dir) error
	genfs *FileSystem
	path  string
	tree  *tree
}

func (d *dirGenerator) Generate(target string) (fs.File, error) {
	if entry, err := d.cache.Get(d.path); err == nil {
		_ = entry
		// TODO: wrap the entry file in a virtualDir
		return nil, fmt.Errorf("cache get not implemented yet")
	}
	// Run the directory generator function
	scopedFS := &scopedFS{d.cache, d.genfs, d.path}
	dir := &Dir{d.cache, d.genfs, d.path, target, d.tree}
	if err := d.fn(scopedFS, dir); err != nil {
		return nil, err
	}
	// Traverse into the directory looking for the target
	if d.path != target {
		return d.genfs.openFrom(d.path, target)
	}
	entry := &virt.File{
		Path:    d.path,
		Mode:    fs.ModeDir,
		Entries: nil, // Entries get filled in on-demand.
	}
	// Cache the directory entry
	if err := d.cache.Set(d.path, entry); err != nil {
		return nil, err
	}
	// Return the virtual directory
	return wrapFile(virt.Open(entry), d.genfs, d.path), nil
}
