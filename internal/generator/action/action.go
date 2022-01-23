package action

import (

	// Embed templates
	_ "embed"

	"gitlab.com/mnm/bud/2/budfs"
	"gitlab.com/mnm/bud/2/di"
	"gitlab.com/mnm/bud/2/gen"
	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/2/parser"
	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed action.gotext
var template string

var generator = gotemplate.MustParse("action.gotext", template)

type Generator struct {
	BFS      budfs.FS
	Injector *di.Injector
	Module   *mod.Module
	Parser   *parser.Parser
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	state, err := Load(g.BFS, g.Injector, g.Module, g.Parser)
	if err != nil {
		return err
	}
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
