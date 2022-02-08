package program

import (
	_ "embed"
	"errors"

	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/gotemplate"
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
	loader := &loader{Generator: g}
	return loader.Load()
}

type loader struct {
	bail.Struct
	*Generator
}

func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	l.Bail(errors.New("not implemented yet"))
	return state, nil
}
