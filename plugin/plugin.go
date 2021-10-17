package plugin

import (
	"github.com/go-duo/bud/bfs"
	"github.com/go-duo/bud/go/mod"
)

func New(modfile mod.File) bfs.Generator {
	return bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
		return nil
	})
}
