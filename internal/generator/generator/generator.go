package generator

import (
	"io/fs"

	"gitlab.com/mnm/bud/gen"
)

type Generator struct {
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	return fs.ErrNotExist
}
