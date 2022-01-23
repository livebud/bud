package command

import (
	_ "embed"

	"gitlab.com/mnm/bud/2/budfs"
	"gitlab.com/mnm/bud/2/gen"
	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/2/parser"
	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command.gotext", template)

type Generator struct {
	BFS    budfs.FS
	Module *mod.Module
	Parser *parser.Parser
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	// Load command state
	state, err := Load(g.BFS, g.Module, g.Parser)
	if err != nil {
		return err
	}
	// Generate our template
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}

// func (g *Generator) generateDI(f gen.F, file *gen.File) error {

// }
