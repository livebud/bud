package tailwind

import (
	"context"

	"github.com/livebud/bud/package/overlay"
)

type Generator struct {
}

func (g *Generator) GenerateDir(ctx context.Context, fsys overlay.F, dir *overlay.Dir) error {
	dir.GenerateFile("tailwind.css", func(ctx context.Context, fsys overlay.F, file *overlay.File) error {
		file.Data = []byte(`/** tailwind **/`)
		return nil
	})
	return nil
}
