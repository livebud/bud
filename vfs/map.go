package vfs

import (
	"io/fs"
)

type Map map[string][]byte

// TODO: support the vfs.ReadWritable interface
var _ fs.FS = (Map)(nil)

func toMemory(m Map) Memory {
	memory := Memory{}
	for path, data := range m {
		memory[path] = &File{Data: []byte(data)}
	}
	return memory
}

func (m Map) Open(name string) (fs.File, error) {
	return toMemory(m).Open(name)
}

func (m Map) MkdirAll(path string, perm fs.FileMode) error {
	return toMemory(m).MkdirAll(path, perm)
}

func (m Map) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return toMemory(m).WriteFile(name, data, perm)
}

func (m Map) RemoveAll(path string) error {
	return toMemory(m).RemoveAll(path)
}
