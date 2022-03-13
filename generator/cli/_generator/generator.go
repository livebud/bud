package generator

import (
	_ "embed"

	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/buddy"
	"gitlab.com/mnm/bud/pkg/gen"
)

//go:embed generator.gotext
var template string

var generator = gotemplate.MustParse("generator.gotext", template)

func New(kit buddy.Kit) *Generator {
	return &Generator{kit}
}

type Generator struct {
	kit buddy.Kit
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
	l.imports.AddNamed("buddy", "gitlab.com/mnm/bud/pkg/buddy")
	l.imports.AddNamed("generator", "gitlab.com/mnm/bud/runtime/cli/generator")
	state.Imports = l.imports.List()
	return state, nil
}
