package vfs

import (
	"io/fs"
)

type Map map[string]string

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
