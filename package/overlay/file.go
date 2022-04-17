package overlay

import (
	"context"

	"github.com/livebud/bud/package/conjure"
)

type FileGenerator interface {
	GenerateFile(ctx context.Context, fsys F, file *File) error
}

// TODO: don't wrap, just extend
type File struct {
	*conjure.File
}

// Link a path
func (f *File) Link(path string) {

}
