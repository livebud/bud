package virtual

import (
	"io"
	"io/fs"
)

type openDir struct {
	*File
	offset int
}

var _ fs.File = (*openDir)(nil)
var _ fs.ReadDirFile = (*openDir)(nil)
var _ fs.DirEntry = (*openDir)(nil)

var _ fs.File = (*openDir)(nil)
var _ fs.ReadDirFile = (*openDir)(nil)

func (d *openDir) Close() error {
	return nil
}

func (d *openDir) Stat() (fs.FileInfo, error) {
	return d.Info()
}

func (d *openDir) Read(p []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.Path, Err: fs.ErrInvalid}
}

func (d *openDir) ReadDir(count int) ([]fs.DirEntry, error) {
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
