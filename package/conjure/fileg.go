package conjure

import (
	"context"
	"fmt"
	"io/fs"
)

type FileGenerator interface {
	GenerateFile(ctx context.Context, file *File) error
}

type fileGenerator struct {
	path string
	fn   func(ctx context.Context, file *File) error
}

func (g *fileGenerator) Generate(ctx context.Context, target string) (fs.File, error) {
	// Prevents prefixes from matching files
	// (e.g. go.mod/go.mod from matching go.mod)
	if g.path != target {
		return nil, fs.ErrNotExist
	}
	file := &File{
		path: target,
	}
	if err := g.fn(ctx, file); err != nil {
		return nil, fmt.Errorf("conjure: generate %q > %w", target, err)
	}
	return file.open()
}
