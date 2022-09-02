package virtual

import "io/fs"

type Entry interface {
	Open() fs.File
}

// Opener is a utility function that implements fs.FS
type Opener func(name string) (fs.File, error)

func (fn Opener) Open(name string) (fs.File, error) {
	return fn(name)
}
