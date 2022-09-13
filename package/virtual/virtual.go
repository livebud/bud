package virtual

import "io/fs"

type Entry interface {
	Open() fs.File
}

// Open is a utility function that implements Entry
type Open func() fs.File

func (fn Open) Open() fs.File {
	return fn()
}
