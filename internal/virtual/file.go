package virtual

import (
	"io"
	"io/fs"
	"time"
)

// File struct
type File struct {
	Name    string
	Data    []byte
	Mode    fs.FileMode
	ModTime time.Time
	Sys     interface{}
	offset  int64
}

var _ fs.File = (*File)(nil)
var _ io.ReadSeeker = (*File)(nil)

// Reset the read data offset to 0
func (f *File) Reset() {
	f.offset = 0
}

func (f *File) Close() error {
	return nil
}

func (f *File) Read(b []byte) (int, error) {
	if f.offset >= int64(len(f.Data)) {
		return 0, io.EOF
	}
	if f.offset < 0 {
		return 0, &fs.PathError{Op: "read", Path: f.Name, Err: fs.ErrInvalid}
	}
	n := copy(b, f.Data[f.offset:])
	f.offset += int64(n)
	return n, nil
}

func (f *File) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		name:    f.Name,
		mode:    f.Mode &^ fs.ModeDir,
		modTime: f.ModTime,
		size:    int64(len(f.Data)),
		sys:     f.Sys,
	}, nil
}

func (f *File) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		// offset += 0
	case 1:
		offset += f.offset
	case 2:
		offset += int64(len(f.Data))
	}
	if offset < 0 || offset > int64(len(f.Data)) {
		return 0, &fs.PathError{Op: "seek", Path: f.Name, Err: fs.ErrInvalid}
	}
	f.offset = offset
	return offset, nil
}
