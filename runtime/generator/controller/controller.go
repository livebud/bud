package controller

import (

	// Embed templates
	"context"
	_ "embed"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/parser"
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
