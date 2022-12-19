package genfs

import (
	"fmt"
	"io/fs"
	gopath "path"

	"github.com/livebud/bud/internal/treefs"
	"github.com/livebud/bud/package/virtual"
)

type Dir struct {
	cache  Cache
	genfs  *fileSystem
	linker Linker
	path   string       // Current directory path
	mode   treefs.Mode  // Current directory mode
	target string       // Final target path
	tree   *treefs.Tree // Filesystem tree containing generators and filler dirs
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

func (d *Dir) Mode() treefs.Mode {
	return d.mode
}

func (d *Dir) GenerateFile(path string, fn func(fsys FS, file *File) error) {
	fpath := gopath.Join(d.path, path)
	fileg := &fileGenerator{d.cache, fn, d.genfs, d.linker, fpath}
	d.tree.Insert(fpath, modeGenerator, fileg)
}

func (d *Dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *Dir) GenerateDir(path string, fn func(fsys FS, dir *Dir) error) {
	fpath := gopath.Join(d.path, path)
	dirg := &dirGenerator{d.cache, fn, d.genfs, d.linker, fpath, d.tree}
	d.tree.Insert(fpath, modeGeneratorDir, dirg)
}

func (d *Dir) DirGenerator(path string, generator DirGenerator) {
	d.GenerateDir(path, generator.GenerateDir)
}

func (d *Dir) ServeFile(path string, fn func(fsys FS, file *File) error) {
	fpath := gopath.Join(d.path, path)
	server := &fileServer{d.cache, fn, d.genfs, d.linker, fpath}
	d.tree.Insert(fpath, modeGeneratorDir, server)
}

func (d *Dir) FileServer(path string, server FileServer) {
	d.ServeFile(path, server.ServeFile)
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
	genfs  *fileSystem
	linker Linker
	path   string
	tree   *treefs.Tree // Filesystem tree containing generators and filler dirs
}

func (d *dirGenerator) Generate(target string) (fs.File, error) {
	if entry, ok := d.cache.Get(d.path); ok {
		_ = entry
		// TODO: wrap the entry file in a virtualDir
		return nil, fmt.Errorf("cache get not implemented yet")
	}
	scopedFS := &scopedFS{d.cache, d.genfs, d.path, d.linker}
	dir := &Dir{d.cache, d.genfs, d.linker, d.path, modeGeneratorDir, target, d.tree}
	if err := d.fn(scopedFS, dir); err != nil {
		return nil, err
	}
	// Traverse into the directory looking for the target
	if d.path != target {
		file, err := d.genfs.open(d.path, target)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	node, ok := d.tree.Get(target)
	if !ok {
		return nil, fmt.Errorf("genfs: unexpected missing tree node for %q", target)
	}
	children := node.Children()
	// Create the virtual directory
	des := make([]fs.DirEntry, len(children))
	for i := range des {
		des[i] = &dirEntry{
			d.genfs,
			children[i].Name(),
			children[i].Mode().FileMode(),
			gopath.Join(d.path, children[i].Name()),
		}
	}
	entry := &virtual.Dir{
		Path:    d.path,
		Mode:    fs.ModeDir,
		Entries: des,
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
