package transform

import (
	"context"
	_ "embed"

	"github.com/livebud/bud/package/overlay"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
)

//go:embed transform.gotext
var template string

var generator = gotemplate.MustParse("transform.gotext", template)

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	Module *gomod.Module
}

type State struct {
	Imports    []*imports.Import
	Transforms []*Transform
}

type Transform struct {
	From     string
	To       string
	Platform string
	Variable string
	Type     string
	Function string
}

func (g *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	imports := imports.New()
	imports.AddNamed("transform", "github.com/livebud/bud/runtime/transform")
	imports.AddNamed("svelte", "github.com/livebud/bud/package/svelte")
	code, err := generator.Generate(State{
		Imports: imports.List(),
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
