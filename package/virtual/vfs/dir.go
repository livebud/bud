package vfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

// Dir creates a new OS filesystem rooted at the given directory.
// TODO: create an os_windows for opening on multiple drives
// with the same API:
// https://github.com/golang/go/issues/44279#issuecomment-955766528
type Dir string

var _ FS = (Dir)("")

func (dir Dir) Open(name string) (fs.File, error) {
	return os.Open(filepath.Join(string(dir), name))
}

func (dir Dir) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(filepath.Join(string(dir), path), perm)
}

func (dir Dir) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filepath.Join(string(dir), name), data, perm)
}

func (dir Dir) RemoveAll(path string) error {
	return os.RemoveAll(filepath.Join(string(dir), path))
}

func (dir Dir) Sub(subdir string) (FS, error) {
	return Dir(filepath.Join(string(dir), subdir)), nil
}
