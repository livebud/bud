package command

import (
	_ "embed"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/parser"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command.gotext", template)

type Generator struct {
	Module *mod.Module
	Parser *parser.Parser
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	// Load command state
	state, err := Load(g.Module, g.Parser)
	if err != nil {
		return err
	}
	// Generate our template
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	// fmt.Println(string(code))
	file.Write(code)
	return nil
}

// func (g *Generator) generateDI(f gen.F, file *gen.File) error {

// }
