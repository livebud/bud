package conjure

import (
	"io/fs"
	"path/filepath"
	"testing/fstest"
)

type GenerateDir func(dir *Dir) error

func (fn GenerateDir) GenerateDir(dir *Dir) error {
	return fn(dir)
}

type DirGenerator interface {
	GenerateDir(dir *Dir) error
}

type dirg struct {
	path   string // defined generator path
	gen    DirGenerator
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
	if err := g.gen.GenerateDir(dir); err != nil {
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
