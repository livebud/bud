package vfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

// Root filesystem
func Root() ReadWritable {
	return &root{}
}

type root struct {
}

// Support absolute
func toFS(name string) (fs.FS, string, error) {
	root := "/"
	if vol := filepath.VolumeName(name); vol != "" {
		root = vol
	}
	rel, err := filepath.Rel(root, name)
	if err != nil {
		return nil, "", err
	}
	rel = filepath.ToSlash(rel)
	return os.DirFS(root), rel, nil
}

func (r *root) Open(name string) (fs.File, error) {
	fsys, name, err := toFS(name)
	if err != nil {
		return nil, err
	}
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	return fsys.Open(name)
}

func (r *root) Stat(name string) (fs.FileInfo, error) {
	fsys, name, err := toFS(name)
	if err != nil {
		return nil, err
	}
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "stat", Path: name, Err: fs.ErrInvalid}
	}
	return fs.Stat(fsys, name)
}

func (r *root) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (r *root) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (r *root) RemoveAll(path string) error {
	return os.RemoveAll(path)
}
