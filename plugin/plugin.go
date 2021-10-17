package plugin

import (
	"gitlab.com/mnm/bud/bfs"
	"gitlab.com/mnm/bud/go/mod"
)

func New(modfile mod.File) bfs.Generator {
	return bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
		return nil
	})
}
