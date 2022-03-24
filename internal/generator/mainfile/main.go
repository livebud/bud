package mainfile

import (
	"context"
	_ "embed"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/package/overlay"
)

//go:embed main.gotext
var template string

var generator = gotemplate.MustParse("main.gotext", template)

// State for the generator
type State struct {
	Imports []*imports.Import
}

func New(module *gomod.Module) *Main {
	return &Main{module}
}

type Main struct {
	module *gomod.Module
}

func (m *Main) Parse(ctx context.Context) (*State, error) {
	imports := imports.New()
	imports.AddStd("os", "context")
	imports.AddNamed("program", m.module.Import("bud/.cli/program"))
	return &State{
		Imports: imports.List(),
	}, nil
}

// Generate a main file
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func (m *Main) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	state, err := m.Parse(ctx)
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
