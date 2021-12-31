package action

import (
	// Embed templates
	_ "embed"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/parser"
)

//go:embed action.gotext
var template string

var generator = gotemplate.MustParse("action.gotext", template)

type Generator struct {
	Injector *di.Injector
	Module   *mod.Module
	Parser   *parser.Parser
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	state, err := Load(g.Injector, g.Module, g.Parser)
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
