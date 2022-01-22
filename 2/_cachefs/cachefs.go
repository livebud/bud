package cachefs

import (
	"io"
	"io/fs"
	"testing/fstest"

	"gitlab.com/mnm/bud/2/singleflight"
)

func New(fsys fs.FS, loader *singleflight.Loader, cache *Store) *FS {
	return &FS{fsys: fsys, loader: loader, cache: cache}
}

type FS struct {
	fsys   fs.FS
	loader *singleflight.Loader
	cache  *Store
}

func (f *FS) Open(name string) (fs.File, error) {
	// Try reading from the cache
	if f.cache.Has(name) {
		return f.cache.Open(name)
	}
	// Load the resource
	file, err := f.loader.Load(f.fsys, name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// Get the stats
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.IsDir() {
		f.cache.Set(name, &fstest.MapFile{
			ModTime: stat.ModTime(),
			Mode:    stat.Mode(),
			Sys:     stat.Sys(),
		})
		return f.cache.Open(name)
	}
	// Read the data fully
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	f.cache.Set(name, &fstest.MapFile{
		Data:    data,
		ModTime: stat.ModTime(),
		Mode:    stat.Mode(),
		Sys:     stat.Sys(),
	})
	// Seek back to the beginning
	if s, ok := file.(io.Seeker); ok {
		if _, err := s.Seek(0, io.SeekStart); err != nil {
			return nil, err
		}
	}
	// Open from the cache
	cached, err := f.cache.Open(name)
	if err != nil {
		return nil, err
	}
	return cached, nil
}

// Events

// type Dir struct {
// 	name    string
// 	mode    fs.FileMode
// 	modTime time.Time
// }

// func (d *Dir) Open(name string) (fs.File, error) {
// 	return &openDir{}
// }

// type Entry struct {
// 	Name    string
// 	Mode    fs.FileMode
// 	ModTime time.Time
// }

// // openDir
// type openDir struct {
// 	path    string
// 	entries []fs.DirEntry
// 	mode    fs.FileMode
// 	modTime time.Time
// 	size    int64
// 	offset  int
// }

// var _ fs.ReadDirFile = (*openDir)(nil)

// func (d *openDir) Close() error {
// 	return nil
// }

// func (d *openDir) Stat() (fs.FileInfo, error) {
// 	return &fileInfo{
// 		name:    filepath.Base(d.path),
// 		mode:    d.mode | fs.ModeDir,
// 		modTime: d.modTime,
// 		size:    d.size,
// 	}, nil
// }

// func (d *openDir) Read(p []byte) (int, error) {
// 	return 0, &fs.PathError{Op: "read", Path: d.path, Err: fs.ErrInvalid}
// }

// func (d *openDir) ReadDir(count int) ([]fs.DirEntry, error) {
// 	n := len(d.entries) - d.offset
// 	if count > 0 && n > count {
// 		n = count
// 	}
// 	if n == 0 && count > 0 {
// 		return nil, io.EOF
// 	}
// 	list := make([]fs.DirEntry, n)
// 	for i := range list {
// 		list[i] = d.entries[d.offset+i]
// 	}
// 	d.offset += n
// 	return list, nil
// }

// // A fileInfo implements fs.FileInfo and fs.DirEntry for a given map file.
// type fileInfo struct {
// 	name    string
// 	data    []byte
// 	size    int64
// 	mode    fs.FileMode
// 	modTime time.Time
// 	sys     interface{}
// }

// func (i *fileInfo) Name() string               { return i.name }
// func (i *fileInfo) Mode() fs.FileMode          { return i.mode }
// func (i *fileInfo) Type() fs.FileMode          { return i.mode.Type() }
// func (i *fileInfo) ModTime() time.Time         { return i.modTime }
// func (i *fileInfo) IsDir() bool                { return i.mode&fs.ModeDir != 0 }
// func (i *fileInfo) Sys() interface{}           { return i.sys }
// func (i *fileInfo) Info() (fs.FileInfo, error) { return i, nil }
// func (i *fileInfo) Size() int64                { return i.size }
