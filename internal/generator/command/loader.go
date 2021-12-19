package command

import (
	"path/filepath"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/imports"
)

// Load state
func Load(module *mod.Module) (*State, error) {
	loader := &loader{
		imports: imports.New(),
		module:  module,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	imports *imports.Set
	module  *mod.Module
}

func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// Add initial imports
	l.imports.AddStd("context")
	l.imports.AddNamed("commander", "gitlab.com/mnm/bud/commander")
	l.imports.AddNamed("console", "gitlab.com/mnm/bud/log/console")
	l.imports.AddNamed("v8", "gitlab.com/mnm/bud/js/v8")
	l.imports.AddNamed("web", l.module.Import("bud", "web"))
	l.imports.AddNamed("deploy", l.module.Import("command", "deploy"))
	l.imports.AddNamed("new", l.module.Import("command", "new"))
	l.imports.AddNamed("new_view", l.module.Import("command", "new", "view"))
	// Load the commands
	state.Command = l.loadCommand("command", ".")
	// Load the imports
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadCommand(base, dir string) *Command {
	command := new(Command)
	command.Name = filepath.Base(dir)
	if command.Name == "." {
		command.Name = imports.AssumedName(l.module.Import())
	}
	// des, err := fs.ReadDir(l.module, filepath.Join(base, dir))
	// if err != nil {
	// 	l.Bail(err)
	// }
	// fmt.Println(command.Name)
	// for _, de := range des {
	// 	fmt.Println(de.Name())
	// }
	return &Command{
		Name:  command.Name,
		Usage: "start your application",
		Subs: []*Command{
			{
				Name:  "deploy",
				Usage: "deploy your application",
				Deps: []di.Dependency{
					&di.Type{
						Import: "gitlab.com/mnm/bud/js/v8",
						Type:   "*Pool",
					},
				},
				Flags: []*Flag{
					{
						Name:  "access-key",
						Usage: "AWS access key",
						Type:  "string",
					},
					{
						Name:  "secret-key",
						Usage: "AWS secret key",
						Type:  "string",
					},
				},
			},
			{
				Name:  "new",
				Usage: "new scaffolding",
				Subs: []*Command{
					{
						Name:  "view",
						Usage: "new view",
						Args: []*Arg{
							{
								Name:  "name",
								Usage: "name of the view",
								Type:  "string",
							},
						},
						Flags: []*Flag{
							{
								Name:    "with-test",
								Usage:   "include a test",
								Type:    "bool",
								Default: "true",
							},
						},
					},
				},
			},
		},
	}
}
