package virtual

import (
	"io"
	"io/fs"
	"time"
)

// File struct
type File struct {
	Path    string
	Data    []byte
	Mode    fs.FileMode
	ModTime time.Time
	Sys     interface{}
	offset  int64
}

var _ io.ReadSeeker = (*File)(nil)
var _ fs.File = (*File)(nil)
var _ Entry = (*File)(nil)

func (f *File) Close() error {
	return nil
}

func (f *File) Read(b []byte) (int, error) {
	if f.offset >= int64(len(f.Data)) {
		return 0, io.EOF
	}
	if f.offset < 0 {
		return 0, &fs.PathError{Op: "read", Path: f.Path, Err: fs.ErrInvalid}
	}
	n := copy(b, f.Data[f.offset:])
	f.offset += int64(n)
	return n, nil
}

func (f *File) Stat() (fs.FileInfo, error) {
	return &fileInfo{
		path:    f.Path,
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
		return 0, &fs.PathError{Op: "seek", Path: f.Path, Err: fs.ErrInvalid}
	}
	f.offset = offset
	return offset, nil
}

func (f *File) Open() fs.File {
	return &File{
		Path:    f.Path,
		Data:    f.Data,
		Mode:    f.Mode,
		ModTime: f.ModTime,
		Sys:     f.Sys,
		offset:  0, // reset offset
	}
}
