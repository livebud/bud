package fscache

import (
	"io/fs"
)

type Entry interface {
	open() fs.File
}
