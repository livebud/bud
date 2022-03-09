package mainfile

import (
	"context"

	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/overlay"
)

type Main struct {
	programPath string
}

func (m *Main) Parse(ctx context.Context) (*State, error) {
	imports := imports.New()
	imports.AddStd("os", "context")
	imports.AddNamed("program", m.programPath)
	return &State{
		Imports: imports.List(),
	}, nil
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
