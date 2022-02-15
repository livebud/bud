package conjure

import (
	"io/fs"
	"path/filepath"
	"testing/fstest"
	"time"
)

type Dir struct {
	gpath   string // generator path
	tpath   string // target path
	Mode    fs.FileMode
	modTime time.Time
	sys     interface{}

	radix  *radix
	filler fstest.MapFS
}

func (d *Dir) Path() string {
	return d.tpath
}

func (d *Dir) File(path string, generator FileGenerator) {
	fullpath := filepath.Join(d.gpath, path)
	d.radix.Set(path, &fileGenerator{fullpath, generator})
	d.filler[fullpath] = &fstest.MapFile{}
}

func (d *Dir) Dir(path string, generator DirGenerator) {
	fullpath := filepath.Join(d.gpath, path)
	d.filler[fullpath] = &fstest.MapFile{Mode: fs.ModeDir}
	d.radix.Set(path, &dirg{
		path:   fullpath,
		gen:    generator,
		filler: d.filler,
	})
}

func (d *Dir) open(rel string) (fs.File, error) {
	// Exact submatch, open generator
	if generator, ok := d.radix.Get(rel); ok {
		return generator.Generate(d.tpath)
	}
	// Get the generator with the longest matching prefix and open that.
	if _, generator, ok := d.radix.GetByPrefix(rel); ok {
		return generator.Generate(d.tpath)
	}
	// Try the filler filesystem
	return d.filler.Open(d.tpath)
}
