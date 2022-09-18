package virtual

import (
	"io"
	"io/fs"
	"path"
	"time"
)

type Dir struct {
	Path    string
	Mode    fs.FileMode
	ModTime time.Time
	Entries []fs.DirEntry
}

var _ fs.DirEntry = (*Dir)(nil)
var _ Entry = (*Dir)(nil)

// Name of the entry. Implements the fs.DirEntry interface.
func (d *Dir) Name() string {
	return path.Base(d.Path)
}

// Returns true if entry is a directory. Implements the fs.DirEntry interface.
func (d *Dir) IsDir() bool {
	return d.Mode.IsDir()
}

// Returns the type of entry. Implements the fs.DirEntry interface.
func (d *Dir) Type() fs.FileMode {
	return d.Mode.Type()
}

// Returns the file info. Implements the fs.DirEntry interface.
func (d *Dir) Info() (fs.FileInfo, error) {
	return &fileInfo{
		path:    d.Path,
		mode:    d.Mode | fs.ModeDir,
		modTime: d.ModTime,
	}, nil
}

func (d *Dir) open() fs.File {
	return &entryDir{d, 0}
}

type entryDir struct {
	*Dir
	offset int
}

var _ fs.File = (*entryDir)(nil)
var _ fs.ReadDirFile = (*entryDir)(nil)

func (d *entryDir) Close() error {
	return nil
}

func (d *entryDir) Stat() (fs.FileInfo, error) {
	return d.Info()
}

func (d *entryDir) Read(p []byte) (int, error) {
	return 0, &fs.PathError{Op: "read", Path: d.Path, Err: fs.ErrInvalid}
}

func (d *entryDir) ReadDir(count int) ([]fs.DirEntry, error) {
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
