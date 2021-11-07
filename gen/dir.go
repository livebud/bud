package gen

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

func (d *Dir) open(f F, key, relative, path string) (fs.File, error) {
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
				case GenerateDir, *serveFS:
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
					case GenerateDir, *serveFS:
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
	mode    fs.FileMode
	modTime time.Time
	size    int64
	offset  int
}

var _ fs.ReadDirFile = (*openDir)(nil)

func (d *openDir) Close() error {
	return nil
}

func (d *openDir) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		name:    filepath.Base(d.path),
		mode:    d.mode | fs.ModeDir,
		modTime: d.modTime,
		size:    d.size,
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

type GenerateDir func(f F, dir *Dir) error

func (fn GenerateDir) open(f F, key, relative, target string) (fs.File, error) {
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
	GenerateDir(f F, dir *Dir) error
}

func DirGenerator(generator dirGenerator) Generator {
	return GenerateDir(generator.GenerateDir)
}

func ServeFS(fsys fs.FS) Generator {
	return &serveFS{fsys}
}

type serveFS struct{ fsys fs.FS }

func (s *serveFS) open(f F, key, relative, target string) (fs.File, error) {
	stat, err := fs.Stat(s.fsys, relative)
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		des, err := fs.ReadDir(s.fsys, relative)
		if err != nil {
			return nil, err
		}
		return &openDir{
			path:    target,
			modTime: stat.ModTime(),
			mode:    stat.Mode(),
			size:    stat.Size(),
			entries: des,
		}, nil
	}
	data, err := fs.ReadFile(s.fsys, relative)
	if err != nil {
		return nil, err
	}
	return &openFile{
		path:    target,
		data:    data,
		modTime: stat.ModTime(),
		mode:    stat.Mode(),
		size:    stat.Size(),
	}, nil
}

// type Entry struct {
// 	path    string
// 	mode    fs.FileMode
// 	entries []fs.DirEntry
// 	modTime time.Time
// 	data    []byte
// }

// func (e *Entry) Path() string {
// 	return e.path
// }

// func (e *Entry) Mode(mode fs.FileMode) {
// 	e.mode = mode
// }

// func (e *Entry) Entry(entries ...fs.DirEntry) {
// 	e.mode = e.mode & fs.ModeDir
// 	e.entries = append(e.entries, entries...)
// }

// func (e *Entry) Write(data []byte) {
// 	e.data = append(e.data, data...)
// }

// func (e *Entry) open(fsys FS, key, relative, path string) (fs.File, error) {
// 	sort.Slice(e.entries, func(i, j int) bool {
// 		return e.entries[i].Name() < e.entries[j].Name()
// 	})
// 	if e.mode&fs.ModeDir != 0 {
// 		return &openDir{
// 			path:    path,
// 			modTime: e.modTime,
// 			entries: e.entries,
// 		}, nil
// 	}
// 	return &openFile{
// 		path:    path,
// 		modTime: e.modTime,
// 		mode:    e.mode,
// 		data:    e.data,
// 	}, nil
// }

// type ServeDir func(f FS, entry *Entry) error

// func (fn ServeDir) open(f FS, key, relative, target string) (fs.File, error) {
// 	entry := &Entry{path: target}
// 	if err := fn(f, entry); err != nil {
// 		return nil, err
// 	}
// 	return entry.open(f, key, relative, target)
// }

// // type ServeDir string

// // var _ Generator = ServeDir("")

// // func (dir ServeDir) open(f FS, key, relative, target string) (fs.File, error) {
// // 	path := filepath.Join(string(dir), relative)
// // 	file, err := os.Open(path)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	defer file.Close()
// // 	stat, err := file.Stat()
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	if stat.IsDir() {
// // 		vdir := newDir(target)
// // 		return vdir.open(f, key, relative, target)
// // 	}
// // 	data, err := ioutil.ReadAll(file)
// // 	if err != nil {
// // 		return nil, err
// // 	}
// // 	vfile := newFile(target)
// // 	vfile.modTime = time.Now()
// // 	vfile.mode = 0644
// // 	vfile.data = data
// // 	return vfile.open(f, key, relative, target)
// // }
