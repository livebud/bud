package action

import (

	// Embed templates
	_ "embed"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/pkg/budfs"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gen"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

//go:embed action.gotext
var template string

var generator = gotemplate.MustParse("action.gotext", template)

type Generator struct {
	BFS      budfs.FS
	Injector *di.Injector
	Module   *gomod.Module
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
