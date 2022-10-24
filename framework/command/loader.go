package command

import (
	"io/fs"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
)

func Load(fsys fs.FS, module *gomod.Module) (*State, error) {
	return (&loader{
		fsys:    fsys,
		module:  module,
		imports: imports.New(),
	}).Load()
}

type loader struct {
	fsys   fs.FS
	module *gomod.Module

	imports *imports.Set
	bail.Struct
}

func (l *loader) Load() (state *State, err error) {
	defer l.Recover2(&err, "app: unable to load state")
	state = new(State)
	l.imports.AddStd("context")
	l.imports.AddNamed("commander", "github.com/livebud/bud/package/commander")
	state.Commands = []*Command{
		{
			FullName: "migrate",
			Name:     "migrate",
			Pascal:   "CommandMigrate",
			Desc:     "migrate the database",
			Type: &Type{
				Path:    l.module.Import("command/migrate"),
				Package: l.imports.Add(l.module.Import("command/migrate")),
				Name:    "*Command",
			},
			Commands: []*Command{
				{
					FullName: "migrate up",
					Name:     "up",
					Pascal:   "Up",
					Desc:     "migrates a database up",
					Input:    l.imports.Add(l.module.Import("command/migrate")) + `.UpInput`,
					Flags: []*Flag{
						{
							Name:    "table",
							Pascal:  "Table",
							Type:    "String",
							Desc:    "table name",
							Default: `"migrations"`,
						},
					},
					Args: []*Arg{
						{
							Name:    "n",
							Pascal:  "N",
							Type:    "Int",
							Default: `1`,
							Desc:    "number of migrations to run up",
						},
					},
					Runnable:   true,
					HasContext: true,
					Package:    "CommandMigrate",
				},
				{
					FullName:   "migrate down",
					Name:       "down",
					Pascal:     "Down",
					Desc:       "migrates a database down",
					Input:      l.imports.Add(l.module.Import("command/migrate")) + `.DownInput`,
					Runnable:   true,
					HasContext: true,
					Package:    "CommandMigrate",
				},
			},
		},
		// {
		// 	Name: "schema",
		// 	Path: "schema",
		// 	Commands: []*Command{
		// 		{
		// 			Name:     "Print",
		// 			Desc:     "print the database schema",
		// 			Runnable: true,
		// 			// Input: &Type{
		// 			// 	Path:    importPath,
		// 			// 	Package: packageName,
		// 			// 	Name:    "*UpInput",
		// 			// },
		// 			Flags: []*Flag{
		// 				// {
		// 				// 	Pascal: "Table",
		// 				// 	Type:   "string",
		// 				// 	Desc:   "table name",
		// 				// },
		// 			},
		// 			Args: []*Arg{
		// 				// {
		// 				// 	Pascal: "URL",
		// 				// 	Type:   "string",
		// 				// 	Desc:   "url to introspect",
		// 				// },
		// 			},
		// 		},
		// 	},
		// },
	}
	state.Imports = l.imports.List()
	return state, nil
}
