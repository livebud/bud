package app

import (
	"context"
	"os"

	"github.com/livebud/bud/example/zero/generator/command"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/runtime/generator"
)

func New(injector *di.Injector, module *gomod.Module) *Generator {
	return &Generator{injector, module}
}

type Generator struct {
	injector *di.Injector
	module   *gomod.Module
}

const appTemplate = `package main

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{$import.Name}} "{{$import.Path}}"
	{{- end }}
)
{{- end }}

func main() {
	app.Main({{ $.Provider.Name }})
}

{{ $.Provider.Function }}
`

var appGen = gotemplate.MustParse("app.gotext", appTemplate)

func (g *Generator) Extend(gen generator.FileSystem) {
	// TODO: should bud/ be implied? I don't think we should sync non-bud/
	// directories, it's too risky.
	gen.GenerateFile("bud/cmd/app/main.go", g.generateFile)
}

func (g *Generator) generateFile(fsys generator.FS, file *generator.File) error {
	type State struct {
		Imports  []*imports.Import
		Provider *di.Provider
	}
	imset := imports.New()
	imset.Add(g.module.Import("generator/app"))
	provider, err := g.injector.Wire(&di.Function{
		Name:    "loadCLI",
		Imports: imset,
		Params: []*di.Param{
			&di.Param{
				Import: "github.com/livebud/bud/package/log",
				Type:   "Log",
			},
		},
		Results: []di.Dependency{
			di.ToType(g.module.Import("bud/internal/command"), "*CLI"),
			&di.Error{},
		},
	})
	if err != nil {
		return err
	}
	code, err := appGen.Generate(State{
		Imports:  imset.List(),
		Provider: provider,
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}

// Runtime from the code generated above

type loadCLI func(log log.Log) (*command.CLI, error)

func Main(loadCLI loadCLI) {
	log := log.New(console.New(os.Stderr))
	if err := run(log, loadCLI); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func run(log log.Log, loadCLI loadCLI) error {
	cli, err := loadCLI(log)
	if err != nil {
		return err
	}
	ctx := context.Background()
	return cli.Parse(ctx, os.Args[1:])
}
