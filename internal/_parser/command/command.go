package command

import (
	"strconv"

	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

// New command parser
func New(module *gomod.Module, parser *parser.Parser) *Parser {
	return &Parser{module, parser}
}

type Parser struct {
	module *gomod.Module
	parser *parser.Parser
}

// TODO: replace with real thing
func (p *Parser) ParseCLI() (*State, error) {
	imset := imports.New()
	imset.AddStd("context")
	imset.AddNamed("generator", p.module.Import("bud/.cli/generator"))
	imset.AddNamed("commander", "gitlab.com/mnm/bud/package/commander")
	imset.AddNamed("build", "gitlab.com/mnm/bud/package/command/build")
	imset.AddNamed("run", "gitlab.com/mnm/bud/package/command/run")

	return &State{
		Imports: imset.List(),
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
							Default: strconv.Quote(p.module.Directory()),
						},
					},
					Deps: []*Dep{
						&Dep{
							Import: &imports.Import{Name: "generator", Path: p.module.Import("bud/.cli/generator")},
							Name:   "Generator",
							Type:   "*generator.Generator",
						},
					},
				},
				&Command{
					Name:     "run",
					Runnable: true,
					Import:   &imports.Import{Name: "run", Path: "gitlab.com/mnm/bud/package/command/run"},
					Flags: []*Flag{
						&Flag{
							Name:    "dir",
							Type:    "string",
							Help:    "project directory",
							Default: strconv.Quote(p.module.Directory()),
						},
					},
					Deps: []*Dep{
						&Dep{
							Import: &imports.Import{Name: "generator", Path: p.module.Import("bud/.cli/generator")},
							Name:   "Generator",
							Type:   "*generator.Generator",
						},
					},
				},
			},
		},
	}, nil
}

// TODO: replace with real thing
func (p *Parser) ParseApp() (*State, error) {
	imset := imports.New()
	imset.AddStd("context")
	imset.AddNamed("generator", p.module.Import("bud/.cli/generator"))
	imset.AddNamed("commander", "gitlab.com/mnm/bud/package/commander")
	imset.AddNamed("build", "gitlab.com/mnm/bud/package/command/build")
	imset.AddNamed("run", "gitlab.com/mnm/bud/package/command/run")

	return &State{
		Imports: imset.List(),
		Command: &Command{
			Name: "app",
			Subs: []*Command{
				&Command{
					Name:     "build",
					Runnable: true,
					Import:   &imports.Import{Name: "build", Path: "gitlab.com/mnm/bud/package/command/build"},
					Flags:    []*Flag{},
					Deps: []*Dep{
						&Dep{
							Import: &imports.Import{Name: "generator", Path: p.module.Import("bud/.cli/generator")},
							Name:   "Generator",
							Type:   "*generator.Generator",
						},
					},
				},
			},
		},
	}, nil
}

// parse state
type parse struct {
	bail.Struct
	imports *imports.Set
	module  *gomod.Module
	parser  *parser.Parser
}

// Parse commands
func (p *parse) Parse() (state *State, err error) {
	defer p.Recover(&err)
	state = new(State)
	// TODO: finish me
	return state, nil
}
