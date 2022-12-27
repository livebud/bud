package virt

import (
	"io/fs"
)

type FS interface {
	fs.FS
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	RemoveAll(path string) error
	Sub(path string) (FS, error)
}

func Open(f *File) fs.File {
	if f.Mode.IsDir() {
		return &openDir{f, 0}
	}
	return &openFile{f, 0}
}
