package program

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
)

//go:embed program.gotext
var template string

var generator = gotemplate.MustParse("program.gotext", template)

var ErrCantWire = errors.New(`program: unable to wire`)

// State of the program code
type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
}

// Generate the program
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New(injector *di.Injector, module *gomod.Module) *Program {
	return &Program{injector, module}
}

type Program struct {
	injector *di.Injector
	module   *gomod.Module
}

func (p *Program) Parse(ctx context.Context) (*State, error) {
	// Default  imports
	imports := imports.New()
	imports.AddStd("errors", "context", "path/filepath", "runtime")
	imports.AddNamed("console", "gitlab.com/mnm/bud/pkg/log/console")
	imports.AddNamed("command", p.module.Import("bud/.cli/command"))
	// imports.AddNamed("gomod", "gitlab.com/mnm/bud/pkg/gomod")
	// imports.AddNamed("trace", "gitlab.com/mnm/bud/package/trace")
	// Write up the dependencies
	provider, err := p.injector.Wire(&di.Function{
		Name:    "loadCLI",
		Imports: imports,
		Target:  p.module.Import("bud/.cli/program"),
		Params: []di.Dependency{
			di.ToType("gitlab.com/mnm/bud/pkg/gomod", "*Module"),
		},
		Aliases: di.Aliases{
			di.ToType("io/fs", "FS"): di.ToType("gitlab.com/mnm/bud/package/overlay", "*FileSystem"),
		},
		Results: []di.Dependency{
			di.ToType(p.module.Import("bud/.cli/command"), "*CLI"),
			&di.Error{},
		},
	})
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
