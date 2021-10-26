package vfs

import (
	"errors"
	"io/fs"
	"os"
	"strings"
	"testing/fstest"
)

type Memory fstest.MapFS
type File = fstest.MapFile

var _ ReadWritable = (Memory)(nil)

func (m Memory) Open(name string) (fs.File, error) {
	return fstest.MapFS(m).Open(name)
}

func (m Memory) MkdirAll(path string, perm fs.FileMode) error {
	// Don't create a directory unless we have to
	if _, err := fs.Stat(m, path); nil == err {
		return nil
	}
	m[path] = &fstest.MapFile{ModTime: Now(), Mode: perm | os.ModeDir}
	return nil
}

func (m Memory) WriteFile(name string, data []byte, perm fs.FileMode) error {
	m[name] = &fstest.MapFile{Data: data, ModTime: Now(), Mode: perm}
	return nil
}

func (m Memory) RemoveAll(path string) error {
	stat, err := fs.Stat(m, path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	// Delete the path
	delete(m, path)
	// Only delete the file
	if !stat.IsDir() {
		return nil
	}
	// Need to delete the rest of the files
	dirpath := path + "/"
	for fpath := range m {
		if strings.HasPrefix(fpath, dirpath) {
			delete(m, fpath)
		}
	}
	return nil
}
