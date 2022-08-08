package virtual

import (
	"io"
	"io/fs"
	"time"
)

type Dir struct {
	Name    string
	Mode    fs.FileMode
	ModTime time.Time
	Sys     interface{}
	Entries []fs.DirEntry
	offset  int
}

var _ fs.ReadDirFile = (*Dir)(nil)
var _ fs.File = (*Dir)(nil)

func (d *Dir) Close() error {
	return nil
}

func (d *Dir) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		name:    d.Name,
		mode:    d.Mode | fs.ModeDir,
		modTime: d.ModTime,
		sys:     d.Sys,
	}, nil
}

func (d *Dir) Read(p []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.Name, Err: fs.ErrInvalid}
}

func (d *Dir) ReadDir(count int) ([]fs.DirEntry, error) {
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
