package budfs

import (
	"io/fs"
	"path"
	"sort"

	"github.com/livebud/bud/internal/virtual"
)

func newFiller() *filler {
	return &filler{
		fsys: map[string]map[string]fs.DirEntry{},
	}
}

type filler struct {
	fsys map[string]map[string]fs.DirEntry
}

var _ fs.FS = (*filler)(nil)

func (f *filler) Open(name string) (fs.File, error) {
	dirMap := f.fsys[name]
	entries := make([]fs.DirEntry, len(dirMap))
	i := 0
	for _, de := range f.fsys[name] {
		entries[i] = de
		i++
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	return &virtual.Dir{
		Name:    name,
		Mode:    dirMode,
		Entries: entries,
	}, nil
}

func (f *filler) Add(target string, generator Generator) {
	if target == "." {
		return
	}
	dirpath := path.Dir(target)
	if _, ok := f.fsys[dirpath]; !ok {
		f.fsys[dirpath] = map[string]fs.DirEntry{}
	}
	basename := path.Base(target)
	f.fsys[dirpath][basename] = &fillerEntry{
		target:    target,
		basename:  basename,
		generator: generator,
	}
	// Recurse until we reach the root
	f.Add(dirpath, &fillerDir{
		Mode: dirMode,
		// TODO: add entries somehow
	})
}

type fillerDir struct {
	Mode fs.FileMode
	// Entries map[string]
}

var _ Generator = (*Embed)(nil)
var _ FileGenerator = (*Embed)(nil)

func (d *fillerDir) Generate(target string) (fs.File, error) {
	return &virtual.Dir{
		Name: target,
		Mode: d.Mode,
	}, nil
}

type fillerEntry struct {
	target    string
	basename  string
	generator Generator
}

var _ fs.DirEntry = (*fillerEntry)(nil)

func (d *fillerEntry) Name() string {
	return d.basename
}

func (d *fillerEntry) IsDir() bool {
	return true
}

var dirMode = 0755 | fs.ModeDir

func (d *fillerEntry) Type() fs.FileMode {
	return dirMode
}

func (d *fillerEntry) Info() (fs.FileInfo, error) {
	file, err := d.generator.Generate(d.target)
	if err != nil {
		return nil, err
	}
	return file.Stat()
}
