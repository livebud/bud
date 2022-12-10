package genfs

import (
	"errors"
	"io/fs"
	gopath "path"
)

// Wrap the file to override ReadDir so that ReadDir reads from the generated
// files
type wrapFile struct {
	fs.File
	genfs fs.FS
	path  string
}

var _ fs.ReadDirFile = (*wrapFile)(nil)

// errNotImplemented mirrors what fs.ReadDir returns when called on a file
var errNotImplemented = errors.New("not implemented")

// Override the default ReadDir so that file stat's can use the generated files
func (d *wrapFile) ReadDir(n int) (des []fs.DirEntry, err error) {
	dir, ok := d.File.(fs.ReadDirFile)
	if !ok {
		return nil, formatError(errNotImplemented, "cannot readdir %q", d.path)
	}
	des, err = dir.ReadDir(n)
	if err != nil {
		return nil, err
	}
	dirEntries := make([]fs.DirEntry, len(des))
	for i := range des {
		dirEntries[i] = &dirEntry{d.genfs, des[i].Name(), des[i].Type(), gopath.Join(d.path, des[i].Name())}
	}
	return dirEntries, nil
}

type dirEntry struct {
	genfs fs.FS
	name  string
	mode  fs.FileMode
	path  string
}

var _ fs.DirEntry = (*dirEntry)(nil)

func (d *dirEntry) Name() string {
	return d.name
}

func (d *dirEntry) Type() fs.FileMode {
	return d.mode
}

func (d *dirEntry) IsDir() bool {
	return d.mode.IsDir()
}

func (d *dirEntry) Info() (fs.FileInfo, error) {
	return fs.Stat(d.genfs, d.path)
}
