package overlay

import (
	"context"

	"github.com/livebud/bud/package/budfs/genfs"
)

type FileGenerator interface {
	GenerateFile(ctx context.Context, fsys F, file *File) error
}

// TODO: don't wrap, just extend
type File struct {
	*genfs.File
}

// Link a path
func (f *File) Link(path string) {

}
