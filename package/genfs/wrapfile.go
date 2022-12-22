package genfs

import (
	"errors"
	"io"
	"io/fs"

	"github.com/livebud/bud/internal/once"
)

// errNotImplemented mirrors what fs.ReadDir returns when called on a file
var errNotImplemented = errors.New("not implemented")

// Open a wrapped file. This is similar to a virtual file except that it calls
// back to genfs as the filesystem for reading directories and getting info.
func wrapFile(file fs.File, genfs fs.FS, path string) *wrappedFile {
	return &wrappedFile{
		File:  file,
		genfs: genfs,
		path:  path,

		offset:      0,
		readOnce:    new(once.Bytes),
		readDirOnce: new(once.DirEntries),
		statOnce:    new(once.FileInfo),
	}
}

// Wrap the file to override ReadDir so that ReadDir reads from the generated
// files
type wrappedFile struct {
	fs.File
	genfs fs.FS
	path  string

	// stateful
	offset      int64
	readOnce    *once.Bytes
	readDirOnce *once.DirEntries
	statOnce    *once.FileInfo
}

var _ fs.File = (*wrappedFile)(nil)
var _ fs.ReadDirFile = (*wrappedFile)(nil)
var _ io.ReadSeeker = (*wrappedFile)(nil)

func (w *wrappedFile) Read(b []byte) (int, error) {
	data, err := w.readOnce.Do(func() ([]byte, error) { return io.ReadAll(w.File) })
	if err != nil {
		return 0, err
	}
	if w.offset >= int64(len(data)) {
		return 0, io.EOF
	}
	if w.offset < 0 {
		return 0, &fs.PathError{Op: "read", Path: w.path, Err: fs.ErrInvalid}
	}
	n := copy(b, data[w.offset:])
	w.offset += int64(n)
	return n, nil
}

func (w *wrappedFile) readDir(count int) (des []fs.DirEntry, err error) {
	if _, ok := w.File.(fs.ReadDirFile); !ok {
		return nil, formatError(errNotImplemented, "cannot readdir %q", w.path)
	}
	des, err = fs.ReadDir(w.genfs, w.path)
	if err != nil {
		return nil, err
	}
	return des, nil
}

// Override the default ReadDir so that file stat's can use the generated files
func (w *wrappedFile) ReadDir(count int) (des []fs.DirEntry, err error) {
	// Read directory entries at most once
	des, err = w.readDirOnce.Do(func() ([]fs.DirEntry, error) { return w.readDir(count) })
	if err != nil {
		return nil, err
	}
	offset := int(w.offset)
	n := len(des) - offset
	if count > 0 && n > count {
		n = count
	}
	if n == 0 && count > 0 {
		return nil, io.EOF
	}
	entries := make([]fs.DirEntry, n)
	for i := range entries {
		entries[i] = des[offset+i]
	}
	w.offset += int64(n)
	return entries, nil
}

func (w *wrappedFile) Seek(offset int64, whence int) (int64, error) {
	stat, err := w.statOnce.Do(func() (fs.FileInfo, error) { return w.File.Stat() })
	if err != nil {
		return 0, err
	}
	switch whence {
	case 0:
		// offset += 0
	case 1:
		offset += w.offset
	case 2:
		offset += stat.Size()
	}
	if offset < 0 || offset > stat.Size() {
		return 0, &fs.PathError{Op: "seek", Path: w.path, Err: fs.ErrInvalid}
	}
	w.offset = offset
	return offset, nil
}
