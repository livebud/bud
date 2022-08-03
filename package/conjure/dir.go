package conjure

import (
	"io/fs"
	"path/filepath"
	"testing/fstest"
)

type Dir struct {
	gpath string // generator path
	Mode  fs.FileMode

	radix  *radix
	filler fstest.MapFS
}

func (d *Dir) Path() string {
	return d.gpath
}

func (d *Dir) GenerateFile(path string, fn func(f *File) error) {
	fullpath := filepath.Join(d.gpath, path)
	d.radix.Set(path, &fileGenerator{fullpath, fn})
	d.filler[fullpath] = &fstest.MapFile{}
}

func (d *Dir) FileGenerator(path string, generator FileGenerator) {
	d.GenerateFile(path, generator.GenerateFile)
}

func (d *Dir) GenerateDir(path string, fn func(d *Dir) error) {
	fullpath := filepath.Join(d.gpath, path)
	d.filler[fullpath] = &fstest.MapFile{Mode: fs.ModeDir}
	d.radix.Set(path, &dirg{
		path:   fullpath,
		fn:     fn,
		filler: d.filler,
	})
}

func (d *Dir) DirGenerator(path string, generator DirGenerator) {
	d.GenerateDir(path, generator.GenerateDir)
}

func (d *Dir) open(target string) (fs.File, error) {
	// TODO: we shouldn't rely on filepath since paths should be agnostic
	// Unfortunately, there doesn't seem to be a path.Rel()
	rel, err := filepath.Rel(d.gpath, target)
	if err != nil {
		return nil, err
	}
	// Exact submatch, open generator
	if generator, ok := d.radix.Get(rel); ok {
		file, err := generator.Generate(target)
		if err != nil {
			return nil, err
		}
		return file, nil
	}
	// Get the generator with the longest matching prefix and open that.
	if _, generator, ok := d.radix.GetByPrefix(rel); ok {
		return generator.Generate(target)
	}
	// Try the filler filesystem
	return d.filler.Open(target)
}
