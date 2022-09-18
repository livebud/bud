package virtual

import "io/fs"

type FS interface {
	fs.FS
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	RemoveAll(path string) error
	Sub(path string) (FS, error)
}

type Entry interface {
	open() fs.File
}

func New(entry Entry) fs.File {
	return entry.open()
}
