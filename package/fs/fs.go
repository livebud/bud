// Package fs is a re-implementation of io/fs but with contexts
package fs

import (
	"context"
	"io/fs"
)

type OpenFS interface {
	OpenContext(ctx context.Context, name string) (fs.File, error)
}

func Open(ctx context.Context, fsys FS, name string) (fs.File, error) {
	if fsys, ok := fsys.(OpenFS); ok {
		return fsys.OpenContext(ctx, name)
	}
	return fsys.Open(name)
}

type FS = fs.FS
type ReadDirFile = fs.ReadDirFile
type PathError = fs.PathError
type FileMode = fs.FileMode
type FileInfo = fs.FileInfo
type DirEntry = fs.DirEntry
type File = fs.File

var ModeDir = fs.ModeDir
var ValidPath = fs.ValidPath

var (
	ErrInvalid    = fs.ErrInvalid    // "invalid argument"
	ErrPermission = fs.ErrPermission // "permission denied"
	ErrExist      = fs.ErrExist      // "file already exists"
	ErrNotExist   = fs.ErrNotExist   // "file does not exist"
	ErrClosed     = fs.ErrClosed     // "file already closed"
)
