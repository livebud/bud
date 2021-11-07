package maingo

import (
	_ "embed"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed maingo.gotext
var template string

var generator = gotemplate.MustParse("maingo", template)

type Generator struct {
}

type State struct{}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	code, err := generator.Generate(State{})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
