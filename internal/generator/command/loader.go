package command

import (
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/bail"
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
	// state.Structs = l.loadStructs(flatten(state.Command))
	// state.Functions = l.loadFunctions(flatten(state.Command))
	// Load the imports
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadCommand(base, dir string) *Command {
	return &Command{
		Slug:  imports.AssumedName(l.module.Import()),
		Usage: "start your application",
		Subs: []*Command{
			{
				Name:  "deploy",
				Usage: "deploy your application",
				Import: &imports.Import{
					Path: l.module.Import("command", "deploy"),
					Name: l.imports.Add(l.module.Import("command", "deploy")),
				},
				Fields: []*Field{
					&Field{
						Import: "gitlab.com/mnm/bud/js/v8",
						Name:   "VM",
						Type:   "*v8.Pool",
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
				Parents: []string{},
			},
			{
				Name:    "new",
				Usage:   "new scaffolding",
				Parents: []string{},
				Import: &imports.Import{
					Path: l.module.Import("command", "new"),
					Name: l.imports.Add(l.module.Import("command", "new")),
				},
				Subs: []*Command{
					{
						Name:  "view",
						Usage: "new view",
						Import: &imports.Import{
							Path: l.module.Import("command", "new", "view"),
							Name: l.imports.Add(l.module.Import("command", "new", "view")),
						},
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
						Parents: []string{"new"},
					},
				},
			},
		},
	}
}

// func (l *loader) loadFunctions(commands []*Command) (fns []*Function) {
// 	return []*Function{
// 		&Function{
// 			Name: "deploy",
// 			Params: []*Param{
// 				&Param{
// 					Name: "vm",
// 					Type: "*v8.Pool",
// 				},
// 			},
// 		},
// 		&Function{
// 			Name: "new",
// 			Params: []*Param{
// 				&Param{
// 					Name: "vm",
// 					Type: "*v8.Pool",
// 				},
// 			},
// 		},
// 		&Function{
// 			Name: "new view",
// 		},
// 	}
// }

// func (l *loader) loadStructs(commands []*Command) (stcts []*Struct) {
// 	return []*Struct{
// 		&Struct{
// 			Name: "Command",
// 			Fields: []*Field{
// 				&Field{
// 					Name: "web",
// 					Type: "*web.Server",
// 				},
// 				&Field{
// 					Name: "Deploy",
// 					Type: "*DeployCommand",
// 				},
// 				&Field{
// 					Name: "New",
// 					Type: "*NewCommand",
// 				},
// 			},
// 		},
// 		&Struct{
// 			Name: "DeployCommand",
// 			Fields: []*Field{
// 				&Field{
// 					Name: "Command",
// 					Type: "*deploy.Command",
// 				},
// 			},
// 		},
// 		&Struct{
// 			Name: "NewCommand",
// 			Fields: []*Field{
// 				&Field{
// 					Name: "Command",
// 					Type: "*new.Command",
// 				},
// 				&Field{
// 					Name: "View",
// 					Type: "*NewViewCommand",
// 				},
// 			},
// 		},
// 		&Struct{
// 			Name: "NewViewCommand",
// 			Fields: []*Field{
// 				&Field{
// 					Name: "Command",
// 					Type: "*new_view.Command",
// 				},
// 			},
// 		},
// 	}
// }
