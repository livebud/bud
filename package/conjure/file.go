package conjure

import (
	"io"
	"io/fs"
	"path"
	"time"
)

type File struct {
	path    string
	Data    []byte
	Mode    fs.FileMode
	modTime time.Time
	sys     interface{}
}

func (f *File) Path() string {
	return f.path
}

// New fs.File
func (f *File) open() (fs.File, error) {
	return &file{f, 0}, nil
}

type file struct {
	*File
	offset int64
}

var _ fs.File = (*file)(nil)
var _ io.ReadSeeker = (*file)(nil)

func (f *file) Close() error {
	return nil
}

func (f *file) Read(b []byte) (int, error) {
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

func (f *file) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		name:    path.Base(f.path),
		mode:    f.Mode &^ fs.ModeDir,
		modTime: f.modTime,
		size:    int64(len(f.Data)),
		sys:     f.sys,
	}, nil
}

func (f *file) Seek(offset int64, whence int) (int64, error) {
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
