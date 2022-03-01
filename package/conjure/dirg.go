package conjure

import (
	"context"
	"io/fs"
	"path/filepath"
	"testing/fstest"
)

type DirGenerator interface {
	GenerateDir(ctx context.Context, dir *Dir) error
}

type dirg struct {
	path   string // defined generator path
	fn     func(ctx context.Context, dir *Dir) error
	filler fstest.MapFS
}

func (g *dirg) Generate(ctx context.Context, target string) (fs.File, error) {
	dir := &Dir{
		gpath:  g.path,
		tpath:  target,
		Mode:   fs.ModeDir,
		filler: g.filler,
		radix:  newRadix(),
	}
	if err := g.fn(ctx, dir); err != nil {
		return nil, err
	}
	// TODO: we shouldn't rely on filepath since paths should be agnostic
	// Unfortunately, there doesn't seem to be a path.Rel()
	rel, err := filepath.Rel(g.path, target)
	if err != nil {
		return nil, err
	}
	return dir.open(ctx, rel)
}
