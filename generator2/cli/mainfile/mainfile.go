package mainfile

import (
	"context"
	_ "embed"

	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"
)

//go:embed mainfile.gotext
var template string

var generator = gotemplate.MustParse("mainfile.gotext", template)

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
}

func (g *Generator) GenerateFile(ctx context.Context, f overlay.F, file *overlay.File) error {
	// Load command state
	state, err := g.Load()
	if err != nil {
		return err
	}
	// Generate our template
	file.Data, err = generator.Generate(state)
	if err != nil {
		return err
	}
	return nil
}

func (g *Generator) Load() (*State, error) {
	loader := &loader{
		Generator: g,
		imports:   imports.New(),
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	*Generator
	imports *imports.Set
}

func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	l.imports.AddStd("os", "context")
	l.imports.AddNamed("program", l.module.Import("bud/.cli/program"))
	return &State{
		Imports: l.imports.List(),
	}, nil
}
