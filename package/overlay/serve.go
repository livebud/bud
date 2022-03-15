package overlay

import "context"

type FileServer interface {
	ServeFile(ctx context.Context, fsys F, file *File) error
}

type ServeFile func(ctx context.Context, fsys F, file *File) error

func (fn ServeFile) ServeFile(ctx context.Context, fsys F, file *File) error {
	return fn(ctx, fsys, file)
}
