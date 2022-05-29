package program

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"io/fs"

	"github.com/livebud/bud/package/vfs"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"
)

//go:embed program.gotext
var template string

var generator = gotemplate.MustParse("program.gotext", template)

var ErrCantWire = errors.New(`program: unable to wire`)

// State of the program code
type State struct {
	Imports  []*imports.Import
	Flags    map[string]string
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

func (p *Program) Parse(ctx context.Context, fsys fs.FS) (*State, error) {
	// Program depends on bud/.cli/command existing
	if err := vfs.Exist(fsys, "bud/.cli/command"); err != nil {
		return nil, err
	}
	// Default  imports
	imports := imports.New()
	imports.AddStd("errors", "context")
	imports.AddNamed("console", "github.com/livebud/bud/package/log/console")
	imports.AddNamed("command", p.module.Import("bud/.cli/command"))
	// Write up the dependencies
	jsVM := di.ToType("github.com/livebud/bud/package/js", "VM")
	loadCLI := &di.Function{
		Name:    "loadCLI",
		Imports: imports,
		Target:  p.module.Import("bud/.cli/program"),
		Params: []di.Dependency{
			di.ToType("github.com/livebud/bud/package/gomod", "*Module"),
		},
		Aliases: di.Aliases{
			jsVM:                     di.ToType("github.com/livebud/bud/package/js/v8", "*VM"),
			di.ToType("io/fs", "FS"): di.ToType("github.com/livebud/bud/package/overlay", "*FileSystem"),
			di.ToType("github.com/livebud/bud/runtime/transform", "*Map"): di.ToType(p.module.Import("bud/.cli/transform"), "*Map"),
		},
		Results: []di.Dependency{
			di.ToType(p.module.Import("bud/.cli/command"), "*CLI"),
			&di.Error{},
		},
	}
	provider, err := p.injector.Wire(loadCLI)
	if err != nil {
		// Don't wrap on purpose, this error gets swallowed up too easily
		return nil, fmt.Errorf("%w > %s", ErrCantWire, err)
	}
	// Add additional imports that we brought in
	for i, im := range provider.Imports {
		provider.Imports[i].Name = imports.Add(im.Path)
	}
	return &State{
		Imports:  imports.List(),
		Provider: provider,
	}, nil
}

func (p *Program) GenerateFile(ctx context.Context, fsys overlay.F, file *overlay.File) error {
	state, err := p.Parse(ctx, fsys)
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
