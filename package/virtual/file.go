package virtual

import (
	"io"
	"io/fs"
	"path"
	"time"
)

// File struct
type File struct {
	Path    string
	Data    []byte
	Mode    fs.FileMode
	ModTime time.Time
}

var _ fs.DirEntry = (*File)(nil)
var _ Entry = (*File)(nil)

// Name of the entry. Implements the fs.DirEntry interface.
func (f *File) Name() string {
	return path.Base(f.Path)
}

// Returns true if entry is a directory. Implements the fs.DirEntry interface.
func (f *File) IsDir() bool {
	return f.Mode.IsDir()
}

// Returns the type of entry. Implements the fs.DirEntry interface.
func (f *File) Type() fs.FileMode {
	return f.Mode.Type()
}

// Returns the file info. Implements the fs.DirEntry interface.
func (f *File) Info() (fs.FileInfo, error) {
	return &fileInfo{
		path:    f.Path,
		mode:    f.Mode &^ fs.ModeDir,
		modTime: f.ModTime,
		size:    int64(len(f.Data)),
	}, nil
}

func (f *File) open() fs.File {
	return &entryFile{f, 0}
}

type entryFile struct {
	*File
	offset int64
}

var _ io.ReadSeeker = (*entryFile)(nil)
var _ fs.File = (*entryFile)(nil)

func (f *entryFile) Close() error {
	return nil
}

func (f *entryFile) Read(b []byte) (int, error) {
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

func (f *entryFile) Stat() (fs.FileInfo, error) {
	return f.Info()
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
		return 0, &fs.PathError{Op: "seek", Path: f.Path, Err: fs.ErrInvalid}
	}
	f.offset = offset
	return offset, nil
}
