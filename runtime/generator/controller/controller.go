package controller

import (

	// Embed templates
	"context"
	_ "embed"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/package/di"
	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/package/parser"
)

//go:embed controller.gotext
var template string

var generator = gotemplate.MustParse("controller.gotext", template)

type Generator struct {
	Injector *di.Injector
	Module   *gomod.Module
	Parser   *parser.Parser
}

func (g *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	state, err := Load(fsys, g.Injector, g.Module, g.Parser)
	if err != nil {
		return err
	}
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
