package command

import (
	_ "embed"
	"fmt"
	"path/filepath"

	"gitlab.com/mnm/bud/internal/di"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command", template)

type Generator struct {
	Modfile  *mod.File
	Injector *di.Injector
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) (err error) {
	// TODO: consider also building when only commands are present
	if err := gen.Exists(f, "bud/web/web.go"); err != nil {
		return err
	}
	imports := imports.New()
	imports.AddStd("context", "errors", "os")
	// imports.AddStd("fmt")
	imports.AddNamed("web", g.Modfile.ModulePath("bud/web"))
	imports.AddNamed("commander", "gitlab.com/mnm/bud/commander")
	imports.AddNamed("console", "gitlab.com/mnm/bud/log/console")
	imports.AddNamed("plugin", "gitlab.com/mnm/bud/plugin")
	imports.AddNamed("mod", "gitlab.com/mnm/bud/go/mod")
	imports.AddNamed("gen", "gitlab.com/mnm/bud/gen")

	// TODO: pull from DI
	imports.AddNamed("js", "gitlab.com/mnm/bud/js/v8")
	imports.AddNamed("controller1", "gitlab.com/mnm/bud/example/hn/bud/controller")
	imports.AddNamed("tailwind", "gitlab.com/mnm/bud-tailwind/transform/tailwind")
	imports.AddNamed("public", "gitlab.com/mnm/bud/example/hn/bud/public")
	imports.AddNamed("transform", "gitlab.com/mnm/bud/example/hn/bud/transform")
	imports.AddNamed("view", "gitlab.com/mnm/bud/example/hn/bud/view")
	imports.AddNamed("web", "gitlab.com/mnm/bud/example/hn/bud/web")
	imports.AddNamed("controller", "gitlab.com/mnm/bud/example/hn/controller")
	imports.AddNamed("hn", "gitlab.com/mnm/bud/example/hn/internal/hn")
	imports.AddNamed("hot", "gitlab.com/mnm/bud/hot")
	imports.AddNamed("router", "gitlab.com/mnm/bud/router")
	imports.AddNamed("svelte", "gitlab.com/mnm/bud/svelte")

	name := filepath.Base(g.Modfile.Directory())

	// 1. Load all the commands
	state := &State{
		Imports: imports.List(),
		Command: &Command{
			Name:  name,
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
		},
	}

	// // 2. DI a virtual *Command
	state.Provider, err = g.Injector.Wire(&di.Function{
		Name:   "load",
		Target: g.Modfile.ModulePath("bud", "command"),
		Params: []di.Dependency{
			&di.Type{Import: "gitlab.com/mnm/bud/go/mod", Type: "*File"},
			&di.Type{Import: "gitlab.com/mnm/bud/gen", Type: "*FileSystem"},
		},
		Results: []di.Dependency{
			&di.Struct{
				Import: g.Modfile.ModulePath("bud", "command"),
				Type:   "*Command",
				Fields: []*di.StructField{
					{
						Name:   "Web",
						Import: g.Modfile.ModulePath("bud", "web"),
						Type:   "*Server",
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	// 3. Generate
	code, err := generator.Generate(state)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(string(code))
	file.Write(code)
	return nil
}
