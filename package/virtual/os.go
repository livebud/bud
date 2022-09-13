package virtual

import (
	"io/fs"
	"os"
	"path/filepath"
)

// OS creates a new OS filesystem rooted at the given directory.
// TODO: create an os_windows for opening on multiple drives
// with the same API:
// https://github.com/golang/go/issues/44279#issuecomment-955766528
type OS string

var _ FS = (OS)("")

func (dir OS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	return os.Open(filepath.Join(string(dir), name))
}

func (dir OS) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(filepath.Join(string(dir), path), perm)
}

func (dir OS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filepath.Join(string(dir), name), data, perm)
}

func (dir OS) RemoveAll(path string) error {
	return os.RemoveAll(filepath.Join(string(dir), path))
}

func (dir OS) Sub(subdir string) (FS, error) {
	return OS(filepath.Join(string(dir), subdir)), nil
}
