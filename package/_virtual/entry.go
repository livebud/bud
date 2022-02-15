package virtual

import (
	"io/fs"
)

type Entry interface {
	New() fs.File
}
