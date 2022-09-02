package virtual

import (
	"io/fs"
	"time"
)

type DirEntry struct {
	Path    string
	Mode    fs.FileMode
	ModTime time.Time
	Sys     interface{}
	Size    int64
}

var _ fs.DirEntry = (*DirEntry)(nil)

func (e *DirEntry) Name() string {
	return e.Path
}

func (e *DirEntry) IsDir() bool {
	return e.Mode&fs.ModeDir != 0
}

func (e *DirEntry) Type() fs.FileMode {
	return e.Mode
}

func (e *DirEntry) Info() (fs.FileInfo, error) {
	return &fileInfo{
		path:    e.Path,
		mode:    e.Mode,
		modTime: e.ModTime,
		sys:     e.Sys,
		size:    e.Size,
	}, nil
}
