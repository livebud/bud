package conjure

import (
	"io/fs"
	"sync"
	"testing/fstest"
)

type DirGenerator interface {
	GenerateDir(dir *Dir) error
}

type dirg struct {
	path   string // defined generator path
	filler fstest.MapFS

	fn   func(dir *Dir) error
	once once
}

type once struct {
	o   sync.Once
	dir *Dir
	err error
}

func (o *once) Do(fn func() (*Dir, error)) (dir *Dir, err error) {
	o.o.Do(func() { o.dir, o.err = fn() })
	return o.dir, o.err
}

func (g *dirg) generateDir(target string) (*Dir, error) {
	return g.once.Do(func() (*Dir, error) {
		dir := &Dir{
			gpath:  g.path,
			Mode:   fs.ModeDir,
			filler: g.filler,
			radix:  newRadix(),
		}
		if err := g.fn(dir); err != nil {
			return nil, err
		}
		return dir, nil
	})
}

func (g *dirg) Generate(target string) (fs.File, error) {
	dir, err := g.generateDir(target)
	if err != nil {
		return nil, err
	}
	return dir.open(target)
}
