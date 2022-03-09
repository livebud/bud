package program

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/di"
)

var ErrCantWire = errors.New(`program: unable to wire`)

type Program struct {
	injector *di.Injector
	function *di.Function
}

func (p *Program) Parse(ctx context.Context) (*State, error) {
	// Default  imports
	imports := imports.New()
	imports.AddStd("os", "errors", "context", "path/filepath", "runtime")
	imports.AddNamed("console", "gitlab.com/mnm/bud/pkg/log/console")
	imports.AddNamed("trace", "gitlab.com/mnm/bud/package/trace")
	// Write up the dependencies
	provider, err := p.injector.Wire(p.function)
	if err != nil {
		// Don't wrap on purpose, this error gets swallowed up too easily
		return nil, fmt.Errorf("%w > %s", ErrCantWire, err)
	}
	// Add additional imports that we brought in
	for _, im := range provider.Imports {
		imports.AddNamed(im.Name, im.Path)
	}
	return &State{
		Imports:  imports.List(),
		Provider: provider,
	}, nil
}

func (p *Program) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	state, err := p.Parse(ctx)
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
