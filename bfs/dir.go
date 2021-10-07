package bfs

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func newDir(key string) *Dir {
	return &Dir{
		path:       key,
		mode:       0755,
		generators: map[string]Generator{},
		watch:      map[string]Event{},
	}
}

type Dir struct {
	path       string
	mode       fs.FileMode
	generators map[string]Generator
	watch      map[string]Event
}

func (d *Dir) Path() string {
	return d.path
}

func (d *Dir) Mode(mode fs.FileMode) {
	d.mode = mode
}

func (d *Dir) Entry(path string, generator Generator) {
	d.generators[path] = generator
}

func (d *Dir) Watch(pattern string, event Event) {
	d.watch[pattern] |= event
}

func (d *Dir) open(f FS, key, relative, path string) (fs.File, error) {
	// Add all the generators to a radix tree
	radix := newRadix()
	for rel, generator := range d.generators {
		radix.Set(rel, generator)
	}
	// Exact submatch, open generator
	if generator, ok := radix.Get(relative); ok {
		file, err := generator.open(f, relative, ".", path)
		if err != nil {
			return nil, err
		}
		switch oe := file.(type) {
		case *openFile:
			oe.path = path
		case *openDir:
			oe.path = path
			d.mergeSynthetic(oe, path)
		default:
			return nil, fmt.Errorf("dfs open dir: unknown type %T", file)
		}
		return file, nil
	}
	// Get the generator with the longest matching prefix and open that.
	if prefix, generator, ok := radix.GetByPrefix(relative); ok {
		relative, err := filepath.Rel(prefix, relative)
		if err != nil {
			return nil, err
		}
		return generator.open(f, filepath.Join(key, prefix), relative, path)
	}
	// Filepath is within the key, create a synthetic file
	return d.synthesize(relative)
}

func (d *Dir) synthesize(path string) (fs.File, error) {
	var entries []fs.DirEntry
	var elem string
	var need = make(map[string]Generator)
	// Handle the root
	if path == "." {
		elem = "."
		for fname, generator := range d.generators {
			i := strings.Index(fname, "/")
			if i < 0 {
				fi := &fileInfo{
					name: fname,
				}
				switch generator.(type) {
				case GenerateDir, ServeFile:
					fi.mode = fs.ModeDir
				}
				entries = append(entries, fi)
			} else {
				need[fname[:i]] = generator
			}
		}
	} else {
		elem = path[strings.LastIndex(path, "/")+1:]
		prefix := path + "/"
		for fname, generator := range d.generators {
			if strings.HasPrefix(fname, prefix) {
				felem := fname[len(prefix):]
				i := strings.Index(felem, "/")
				if i < 0 {
					fi := &fileInfo{
						name: felem,
					}
					switch generator.(type) {
					case GenerateDir, ServeFile:
						fi.mode = fs.ModeDir
					}
					entries = append(entries, fi)
				} else {
					need[fname[len(prefix):len(prefix)+i]] = generator
				}
			}
		}
		if entries == nil && len(need) == 0 {
			return nil, &fs.PathError{Op: "open", Path: path, Err: fs.ErrNotExist}
		}
	}
	for _, entry := range entries {
		delete(need, entry.Name())
	}
	for fname := range need {
		fi := &fileInfo{
			name: fname,
			// everything in "need" is a parent node, so it's always a dir
			mode: fs.ModeDir,
		}
		entries = append(entries, fi)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
	return &openDir{
		path:    elem,
		entries: entries,
	}, nil
}

// Merge generated directories with generator paths
func (d *Dir) mergeSynthetic(dir *openDir, path string) {
	file, err := d.synthesize(path)
	if err != nil {
		return
	}
	sd, ok := file.(*openDir)
	if !ok {
		return
	}
	dir.entries = append(dir.entries, sd.entries...)
	sort.Slice(dir.entries, func(i, j int) bool {
		return dir.entries[i].Name() < dir.entries[j].Name()
	})
}

// openDir
type openDir struct {
	path    string
	entries []fs.DirEntry
	modTime time.Time
	offset  int
}

var _ fs.ReadDirFile = (*openDir)(nil)

func (d *openDir) Close() error {
	return nil
}

func (d *openDir) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		name:    filepath.Base(d.path),
		mode:    fs.ModeDir,
		modTime: d.modTime,
	}, nil
}

func (d *openDir) Read(p []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.path, Err: fs.ErrInvalid}
}

func (d *openDir) ReadDir(count int) ([]fs.DirEntry, error) {
	n := len(d.entries) - d.offset
	if count > 0 && n > count {
		n = count
	}
	if n == 0 && count > 0 {
		return nil, io.EOF
	}
	list := make([]fs.DirEntry, n)
	for i := range list {
		list[i] = d.entries[d.offset+i]
	}
	d.offset += n
	return list, nil
}

type GenerateDir func(f FS, dir *Dir) error

func (fn GenerateDir) open(f FS, key, relative, target string) (fs.File, error) {
	dir := newDir(key)
	if err := fn(f, dir); err != nil {
		return nil, err
	}
	for to, event := range dir.watch {
		f.link(key, to, event)
	}
	return dir.open(f, key, relative, target)
}

type dirGenerator interface {
	GenerateDir(f FS, dir *Dir) error
}

func DirGenerator(generator dirGenerator) Generator {
	return GenerateDir(generator.GenerateDir)
}
