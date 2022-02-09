package program

import (
	_ "embed"
	"fmt"

	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gen"
	"gitlab.com/mnm/bud/pkg/gomod"
)

//go:embed program.gotext
var template string

var generator = gotemplate.MustParse("program.gotext", template)

func New(genFS gen.FS, injector *di.Injector, module *gomod.Module) *Generator {
	return &Generator{genFS, injector, module}
}

type Generator struct {
	genFS    gen.FS
	injector *di.Injector
	module   *gomod.Module
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	// Load command state
	state, err := g.Load()
	if err != nil {
		return err
	}
	// Generate our template
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}

func (g *Generator) Load() (*State, error) {
	loader := &loader{Generator: g, imports: imports.New()}
	return loader.Load()
}

type loader struct {
	bail.Struct
	*Generator
	imports *imports.Set
}

func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// Add imports
	l.imports.AddStd("errors", "context", "path/filepath", "runtime")
	l.imports.AddNamed("console", "gitlab.com/mnm/bud/pkg/log/console")
	l.imports.AddNamed("buddy", "gitlab.com/mnm/bud/pkg/buddy")
	// Inject the provider
	state.Provider, err = l.injector.Wire(&di.Function{
		Name:   "loadCLI",
		Target: l.module.Import("bud/.cli/program"),
		Params: []di.Dependency{
			di.ToType("gitlab.com/mnm/bud/pkg/buddy", "*Driver"),
		},
		Results: []di.Dependency{
			di.ToType(l.module.Import("bud/.cli/command"), "*CLI"),
			&di.Error{},
		},
	})
	if err != nil {
		l.Bail(fmt.Errorf("program unable to wire dependencies > %w", err))
		return
	}
	// Add the imports we find
	for _, im := range state.Provider.Imports {
		l.imports.AddNamed(im.Name, im.Path)
	}
	// Return a list of imports
	state.Imports = l.imports.List()
	return state, nil
}
