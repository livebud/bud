package importfile

import (
	"context"
	_ "embed"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/package/overlay"
)

//go:embed import.gotext
var template string

var generator = gotemplate.MustParse("import.gotext", template)

// State for the generator
type State struct {
	Programs []*imports.Import
	Plugins  []*imports.Import // TODO
}

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
}

func (g *Generator) Parse(ctx context.Context) (*State, error) {
	return &State{
		Programs: []*imports.Import{
			{
				Name: "_",
				Path: g.module.Import("bud/.cli/program"),
			},
			{
				Name: "_",
				Path: g.module.Import("bud/.app/program"),
			},
		},
	}, nil
}

// Generate a main file
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func (g *Generator) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	state, err := g.Parse(ctx)
	if err != nil {
		return err
	}
	// Generate code from the state
	code, err := Generate(state)
	if err != nil {
		return err
	}
	// Write to the file
	file.Data = code
	return nil
}
