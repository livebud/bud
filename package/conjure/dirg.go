package conjure

import (
	"io/fs"
	"path/filepath"
	"testing/fstest"
)

type DirGenerator interface {
	GenerateDir(dir *Dir) error
}

type dirg struct {
	path   string // defined generator path
	fn     func(dir *Dir) error
	filler fstest.MapFS
}

func (g *dirg) Generate(target string) (fs.File, error) {
	dir := &Dir{
		gpath:  g.path,
		tpath:  target,
		Mode:   fs.ModeDir,
		filler: g.filler,
		radix:  newRadix(),
	}
	if err := g.fn(dir); err != nil {
		return nil, err
	}
	// TODO: we shouldn't rely on filepath since paths should be agnostic
	// Unfortunately, there doesn't seem to be a path.Rel()
	rel, err := filepath.Rel(g.path, target)
	if err != nil {
		return nil, err
	}
	return dir.open(rel)
}
