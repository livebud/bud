package virtual

import "io/fs"

type Entry interface {
	Open() fs.File
}
