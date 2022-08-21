package virtual

import "io/fs"

type Entry interface {
	Open() fs.File
}

// WFS is a write-only filesystem.
type WritableFS interface {
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	RemoveAll(path string) error
}

// RWFS is a read-write filesystem.
type RWFS interface {
	fs.FS
	WritableFS
}
