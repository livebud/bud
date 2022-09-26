package tailwind

import (
	"fmt"

	"github.com/livebud/bud/package/budfs"
)

type Generator struct {
}

func (g *Generator) GenerateDir(fsys budfs.FS, dir *budfs.Dir) error {
	fmt.Println("got dir...")
	return nil
}
