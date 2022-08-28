package virtual

import (
	"io/fs"
)

type DirEntry struct {
	Path    string
	ModeDir bool
}

var _ fs.DirEntry = (*DirEntry)(nil)

func (e *DirEntry) Name() string {
	return e.Path
}

func (e *DirEntry) IsDir() bool {
	return e.ModeDir
}

func (e *DirEntry) Type() fs.FileMode {
	if e.ModeDir {
		return fs.ModeDir
	}
	return 0
}

func (e *DirEntry) Info() (fs.FileInfo, error) {
	return &FileInfo{
		Path:    e.Path,
		ModeDir: e.ModeDir,
	}, nil
}
