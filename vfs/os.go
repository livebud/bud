package vfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

type OS string

var _ ReadWritable = (OS)("")

func (dir OS) Open(name string) (fs.File, error) {
	return os.DirFS(string(dir)).Open(name)
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
