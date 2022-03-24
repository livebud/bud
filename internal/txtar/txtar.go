package txtar

import (
	"gitlab.com/mnm/bud/package/vfs"
	"golang.org/x/tools/txtar"
)

// ParseFile parse a txtar file into a virtual filesystem. Used for tests
func ParseFile(path string) (vfs.Memory, error) {
	archive, err := txtar.ParseFile(path)
	if err != nil {
		return nil, err
	}
	memory := vfs.Memory{}
	for _, file := range archive.Files {
		memory[file.Name] = &vfs.File{Data: file.Data}
	}
	return memory, nil
}
