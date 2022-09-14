package virtual

import (
	"io/fs"
	"path"
	"time"
)

type DirEntry struct {
	Path    string
	Size    int64
	Mode    fs.FileMode
	ModTime time.Time
}

var _ fs.DirEntry = (*DirEntry)(nil)

func (e *DirEntry) Name() string {
	return path.Base(e.Path)
}

func (e *DirEntry) IsDir() bool {
	return e.Mode&fs.ModeDir != 0
}

func (e *DirEntry) Type() fs.FileMode {
	return e.Mode.Type()
}

func (e *DirEntry) Info() (fs.FileInfo, error) {
	return &fileInfo{
		path:    e.Path,
		mode:    e.Mode,
		modTime: e.ModTime,
		size:    e.Size,
	}, nil
}
