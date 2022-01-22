package virtual

import (
	"io/fs"
)

type Entry interface {
	open() fs.File
}
