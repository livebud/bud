package vfs

import (
	"io"
	"io/fs"
	"path"
	"time"
)

type Entry struct {
	Data    []byte
	Mode    fs.FileMode
	ModTime time.Time
}

type dirEntry struct {
	*Entry
	path string
}

var _ fs.DirEntry = (*dirEntry)(nil)

// Name of the entry. Implements the fs.DirEntry interface.
func (de *dirEntry) Name() string {
	return path.Base(de.path)
}

// Returns true if entry is a directory. Implements the fs.DirEntry interface.
func (de *dirEntry) IsDir() bool {
	return de.Mode.IsDir()
}

// Returns the type of entry. Implements the fs.DirEntry interface.
func (de *dirEntry) Type() fs.FileMode {
	return de.Mode.Type()
}

// Returns the file info. Implements the fs.DirEntry interface.
func (de *dirEntry) Info() (fs.FileInfo, error) {
	return &fileInfo{
		path:    de.path,
		mode:    de.Mode,
		modTime: de.ModTime,
		size:    int64(len(de.Data)),
		sys:     nil,
	}, nil
}

type entryFile struct {
	*Entry
	path   string
	offset int64
}

var _ fs.File = (*entryFile)(nil)

func (f *entryFile) Close() error {
	return nil
}

func (f *entryFile) Read(b []byte) (int, error) {
	if f.offset >= int64(len(f.Data)) {
		return 0, io.EOF
	}
	if f.offset < 0 {
		return 0, &fs.PathError{Op: "read", Path: f.path, Err: fs.ErrInvalid}
	}
	n := copy(b, f.Data[f.offset:])
	f.offset += int64(n)
	return n, nil
}

func (f *entryFile) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		path:    f.path,
		mode:    f.Mode &^ fs.ModeDir,
		modTime: f.ModTime,
		size:    int64(len(f.Data)),
		sys:     nil,
	}, nil
}

func (f *entryFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		// offset += 0
	case 1:
		offset += f.offset
	case 2:
		offset += int64(len(f.Data))
	}
	if offset < 0 || offset > int64(len(f.Data)) {
		return 0, &fs.PathError{Op: "seek", Path: f.path, Err: fs.ErrInvalid}
	}
	f.offset = offset
	return offset, nil
}

type entryDir struct {
	*Entry
	entries []fs.DirEntry
	path    string
	offset  int
}

var _ fs.File = (*entryDir)(nil)
var _ fs.ReadDirFile = (*entryDir)(nil)

func (d *entryDir) Close() error {
	return nil
}

func (d *entryDir) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		path:    d.path,
		mode:    d.Mode | fs.ModeDir,
		modTime: d.ModTime,
		sys:     nil,
	}, nil
}

func (d *entryDir) Read(p []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.path, Err: fs.ErrInvalid}
}

func (d *entryDir) ReadDir(count int) ([]fs.DirEntry, error) {
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

// A fileInfo implements fs.FileInfo and fs.DirEntry for a given map file.
// This is a copy of the virtual package. We copy this because there's no
// way to make the mode field public without breaking the fs.FileInfo interface.
type fileInfo struct {
	path    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	sys     interface{}
}

var _ fs.FileInfo = (*fileInfo)(nil)
var _ fs.DirEntry = (*fileInfo)(nil)

func (i *fileInfo) Name() string               { return path.Base(i.path) }
func (i *fileInfo) Mode() fs.FileMode          { return fs.FileMode(i.mode) }
func (i *fileInfo) Type() fs.FileMode          { return i.mode.Type() }
func (i *fileInfo) ModTime() time.Time         { return i.modTime }
func (i *fileInfo) IsDir() bool                { return i.mode&fs.ModeDir != 0 }
func (i *fileInfo) Sys() interface{}           { return i.sys }
func (i *fileInfo) Info() (fs.FileInfo, error) { return i, nil }
func (i *fileInfo) Size() int64                { return i.size }
