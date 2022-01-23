package web

import (
	_ "embed"

	"gitlab.com/mnm/bud/2/budfs"
	"gitlab.com/mnm/bud/2/gen"
	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/2/parser"
	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed web.gotext
var template string

var generator = gotemplate.MustParse("web", template)

type Generator struct {
	BFS    budfs.FS
	Module *mod.Module
	Parser *parser.Parser
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	state, err := Load(g.BFS, g.Module, g.Parser)
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
