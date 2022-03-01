package conjure

import (
	"context"
	"fmt"
	"io/fs"
)

type FileServer interface {
	ServeFile(ctx context.Context, file *File) error
}

type serverg struct {
	path string // defined generator path
	fn   func(ctx context.Context, f *File) error
}

func (g *serverg) Generate(ctx context.Context, target string) (fs.File, error) {
	if target == g.path {
		return nil, fs.ErrInvalid
	}
	file := &File{
		path: target,
	}
	if err := g.fn(ctx, file); err != nil {
		return nil, fmt.Errorf("conjure: serve %q > %w", target, err)
	}
	return file.open()
}
