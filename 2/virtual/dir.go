package virtual

import (
	"io"
	"io/fs"
	"time"
)

type Dir struct {
	Name    string
	Entries []fs.DirEntry
	Mode    fs.FileMode
	ModTime time.Time
	Sys     interface{}
}

var _ Entry = (*Dir)(nil)

func (d *Dir) open() fs.File {
	return &dir{d, 0}
}

type dir struct {
	*Dir
	offset int
}

var _ fs.ReadDirFile = (*dir)(nil)
var _ fs.File = (*dir)(nil)

func (d *dir) Close() error {
	return nil
}

func (d *dir) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		name:    d.Name,
		mode:    d.Mode | fs.ModeDir,
		modTime: d.ModTime,
		sys:     d.Sys,
	}, nil
}

func (d *dir) Read(p []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.Name, Err: fs.ErrInvalid}
}

func (d *dir) ReadDir(count int) ([]fs.DirEntry, error) {
	n := len(d.Entries) - d.offset
	if count > 0 && n > count {
		n = count
	}
	if n == 0 && count > 0 {
		return nil, io.EOF
	}
	list := make([]fs.DirEntry, n)
	for i := range list {
		list[i] = d.Entries[d.offset+i]
	}
	d.offset += n
	return list, nil
}
