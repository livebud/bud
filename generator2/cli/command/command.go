package command

import (
	_ "embed"
	"strconv"

	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command.gotext", template)

func New(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
}

func (g *Generator) GenerateFile(f overlay.F, file *overlay.File) error {
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
	l.imports.AddNamed("generator", l.module.Import("bud/.cli/generator"))
	l.imports.AddNamed("commander", "gitlab.com/mnm/bud/pkg/commander")
	l.imports.AddNamed("build", "gitlab.com/mnm/bud/package/command/build")
	// state.Imports = l.imports.List()

	// TODO: finish state
	state = &State{
		Imports: l.imports.List(),
		Command: &Command{
			Name: "cli",
			Subs: []*Command{
				&Command{
					Name:     "build",
					Runnable: true,
					Import:   &imports.Import{Name: "build", Path: "gitlab.com/mnm/bud/package/command/build"},
					Flags: []*Flag{
						&Flag{
							Name:    "dir",
							Type:    "string",
							Help:    "project directory",
							Default: strconv.Quote(l.module.Directory()),
						},
					},
					Deps: []*Dep{
						&Dep{
							Import: &imports.Import{Name: "generator", Path: l.module.Import("bud/.cli/generator")},
							Name:   "Generator",
							Type:   "*generator.Generator",
						},
					},
				},
				// &Command{
				// 	Name:   "run",
				// 	Import: &imports.Import{Name: "run", Path: ""},
				// },
			},
		},
	}
	return state, nil
}
