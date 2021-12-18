package command

import (
	"path/filepath"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/imports"
)

// Load state
func Load(injector *di.Injector, module *mod.Module) (*State, error) {
	loader := &loader{
		imports:  imports.New(),
		injector: injector,
		module:   module,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	imports  *imports.Set
	injector *di.Injector
	module   *mod.Module
}

func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// Add initial imports
	l.imports.AddStd("os", "context", "errors")
	l.imports.AddNamed("commander", "gitlab.com/mnm/bud/commander")
	l.imports.AddNamed("console", "gitlab.com/mnm/bud/log/console")
	l.imports.AddNamed("plugin", "gitlab.com/mnm/bud/plugin")
	l.imports.AddNamed("mod", "gitlab.com/mnm/bud/go/mod")
	l.imports.AddNamed("gen", "gitlab.com/mnm/bud/gen")
	// Load the commands
	state.Command = l.loadCommand("command", ".")
	// Load the provider
	state.Provider = l.loadProvider(state.Command)
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
		Name:  filepath.Base(dir),
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
						Name:    "access-key",
						Usage:   "include a test",
						Type:    "*bool",
						Default: "true",
					},
					{
						Name:    "secret-key",
						Usage:   "include a test",
						Type:    "*bool",
						Default: "true",
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
								Type:    "*bool",
								Default: "true",
							},
						},
					},
				},
			},
		},
	}
}

func (l *loader) loadProvider(cmd *Command) *Provider {
	provider, err := l.injector.Wire(&di.Function{
		Name:   "load",
		Target: l.module.Import("bud", "command"),
		Params: []di.Dependency{
			&di.Type{Import: "gitlab.com/mnm/bud/go/mod", Type: "*Module"},
			&di.Type{Import: "gitlab.com/mnm/bud/gen", Type: "*FileSystem"},
		},
		Results: []di.Dependency{
			&di.Struct{
				Import: l.module.Import("bud", "command"),
				Type:   "*Command",
				Fields: []*di.StructField{
					{
						Name:   "web",
						Import: l.module.Import("bud", "web"),
						Type:   "*Server",
					},
					// {
					// 	Name:   "Deploy",
					// 	Import: l.module.Import("bud", "command"),
					// 	Type:   "*deployCommand",
					// },
				},
			},
			&di.Error{},
		},
	})
	if err != nil {
		l.Bail(err)
	}
	// Add the imports
	for _, im := range provider.Imports {
		l.imports.AddNamed(im.Name, im.Path)
	}
	return provider
}
