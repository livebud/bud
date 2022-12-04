// Package vfs is deprecated and will be deleted once remaining code is migrated
// All new functionality should go in the virtual package.
package vfs

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

type Writable interface {
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	RemoveAll(path string) error
}

type ReadWritable interface {
	fs.FS
	Writable
}

// Now may be overridden for testing purposes
var Now = func() time.Time {
	return time.Now()
}

// WriteAll the filesystem at "from" to "to"
func WriteAll(from, to string, fsys fs.FS) error {
	return fs.WalkDir(fsys, from, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		toPath := filepath.Join(to, path)
		if de.IsDir() {
			mode := de.Type()
			if mode == fs.ModeDir {
				mode = fs.FileMode(0755)
			}
			return os.MkdirAll(toPath, mode)
		}
		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		mode := de.Type()
		if mode == 0 {
			mode = fs.FileMode(0644)
		}
		return os.WriteFile(toPath, data, mode)
	})
}

func Write(to string, fsys fs.FS) error {
	return WriteAll(".", to, fsys)
}
