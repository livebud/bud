package generator

import (
	"gitlab.com/mnm/bud/gen"
)

type Generator struct {
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	return file.Skip()
}
