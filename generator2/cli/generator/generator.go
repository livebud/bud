package generator

import (
	"context"
	_ "embed"

	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/overlay"
)

//go:embed generator.gotext
var template string

var generator = gotemplate.MustParse("generator.gotext", template)

func New() *Generator {
	return &Generator{}
}

type Generator struct {
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
	l.imports.AddStd("context")
	l.imports.AddNamed("overlay", "gitlab.com/mnm/bud/package/overlay")
	l.imports.AddNamed("mainfile", "gitlab.com/mnm/bud/generator2/app/mainfile")
	l.imports.AddNamed("program", "gitlab.com/mnm/bud/generator2/app/program")
	l.imports.AddNamed("process", "gitlab.com/mnm/bud/generator2/app/process")
	state.Imports = l.imports.List()
	// TODO: finish state
	state = &State{
		Imports: l.imports.List(),
		Generators: []*DirGenerator{
			&DirGenerator{
				Path:   "bud/.app",
				Import: &imports.Import{Name: "mainfile", Path: "gitlab.com/mnm/bud/generator2/app/mainfile"},
			},
			&DirGenerator{
				Path:   "bud/.app/program",
				Import: &imports.Import{Name: "program", Path: "gitlab.com/mnm/bud/generator2/app/program"},
			},
			&DirGenerator{
				Path:   "bud/.app/process",
				Import: &imports.Import{Name: "process", Path: "gitlab.com/mnm/bud/generator2/app/process"},
			},
		},
	}
	return state, nil
}
