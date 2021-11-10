package controller

import (
	"io/fs"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
)

type Generator struct {
	Modfile mod.File
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	return fs.ErrNotExist
}

type Parser struct {
}
