package virtual

import (
	"io/fs"
	"path"
	"time"
)

// FileInfo implements fs.FileInfo and fs.DirEntry for a given map file.
type FileInfo struct {
	Path    string
	ModeDir bool
}

var _ fs.FileInfo = (*FileInfo)(nil)
var _ fs.DirEntry = (*FileInfo)(nil)

var zero time.Time

const unknownSize = -1

func (i *FileInfo) Name() string               { return path.Base(i.Path) }
func (i *FileInfo) Type() fs.FileMode          { return i.Mode() }
func (i *FileInfo) ModTime() time.Time         { return zero }
func (i *FileInfo) IsDir() bool                { return i.ModeDir }
func (i *FileInfo) Sys() interface{}           { return nil }
func (i *FileInfo) Info() (fs.FileInfo, error) { return i, nil }
func (i *FileInfo) Size() int64                { return unknownSize }

func (i *FileInfo) Mode() fs.FileMode {
	if i.ModeDir {
		return fs.ModeDir
	}
	return 0
}
