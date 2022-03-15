package overlay

import (
	"context"

	"gitlab.com/mnm/bud/package/conjure"
)

type FileGenerator interface {
	GenerateFile(ctx context.Context, fsys F, file *File) error
}

type File struct {
	*conjure.File
}

// Link a path
func (f *File) Link(path string) {

}

type Embed = conjure.Embed
