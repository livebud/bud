package cli

import (
	"context"
	"fmt"

	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"

	"github.com/livebud/bud/package/gotemplate"

	"github.com/livebud/bud/package/virtual"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/package/genfs"
)

type Generate2 struct {
	Flag      *framework.Flag
	ListenDev string
	Packages  []string
}

func (c *CLI) Generate2(ctx context.Context, in *Generate2) (err error) {
	// Load the logger if not already provided
	log, err := c.loadLog()
	if err != nil {
		return err
	}
	log = log.Field("method", "Generate2").Field("package", "cli")

	// Find the module if not already provided
	module, err := c.findModule()
	if err != nil {
		return err
	}

	gen := genfs.New(dag.Discard, module, log)
	parser := parser.New(gen, module)
	injector := di.New(gen, log, module, parser)
	gen.FileGenerator("bud/cmd/gen/main.go", &mainGenerator{injector, log, module})
	gen.FileGenerator("bud/internal/gen/generator/generator.go", &generatorGenerator{log, module})
	if err := virtual.Sync(log, gen, module, "bud"); err != nil {
		return err
	}

	// Build bud/gen
	cmd := c.command(module.Directory(), "go", "build", "-mod=mod", "-o=bud/gen", "./bud/cmd/gen")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Run bud/gen
	cmd = c.command(module.Directory(), "./bud/gen")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Build bud/app
	cmd = c.command(module.Directory(), "go", "build", "-mod=mod", "-o=bud/app", "./bud/cmd/app")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Run bud/app
	cmd = c.command(module.Directory(), "./bud/app")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

type mainGenerator struct {
	injector *di.Injector
	log      log.Log
	module   *gomod.Module
}

const mainTemplate = `package main

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{$import.Name}} "{{$import.Path}}"
	{{- end }}
)
{{- end }}

func main() {
	gen.Main({{ $.Provider.Name }})
}

{{ $.Provider.Function }}
`

var mainGen = gotemplate.MustParse("bud/cmd/gen/main.go", mainTemplate)

func (g *mainGenerator) GenerateFile(fsys genfs.FS, file *genfs.File) error {
	g.log.Info("generating file", file.Path())
	type State struct {
		Imports  []*imports.Import
		Provider *di.Provider
	}
	imset := imports.New()
	imset.AddNamed("gen", "github.com/livebud/bud/runtime/gen")
	provider, err := g.injector.Wire(&di.Function{
		Name:    "loadGenerator",
		Imports: imset,
		Aliases: di.Aliases{
			di.ToType("github.com/livebud/bud/package/log", "Log"):       di.ToType("github.com/livebud/bud/runtime/gen", "Log"),
			di.ToType("github.com/livebud/bud/package/gomod", "*Module"): di.ToType("github.com/livebud/bud/runtime/gen", "*Module"),
		},
		Results: []di.Dependency{
			di.ToType(g.module.Import("bud/internal/gen/generator"), "*Generator"),
			&di.Error{},
		},
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	code, err := mainGen.Generate(State{
		Imports:  imset.List(),
		Provider: provider,
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}

type generatorGenerator struct {
	log    log.Log
	module *gomod.Module
}

const generatorTemplate = `package generator

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{$import.Name}} "{{$import.Path}}"
	{{- end }}
)
{{- end }}

func New(
	log log.Log,
	module *gomod.Module,
	{{- range $generator := $.Generators }}
	{{ $generator.Camel }} *{{ $generator.Import.Name }}.{{ $generator.Type }},
	{{- end }}
) *Generator {
	return generator.New(log, module, &generator.Schema{
		GenerateFiles: map[string]generator.GenerateFile{
			{{- range $generator := $.Generators }}
			"{{ $generator.Path }}": {{ $generator.Camel }}.GenerateFile,
			{{- end }}
		},
	})
}

type Generator = generator.Generator
`

var generatorGen = gotemplate.MustParse("bud/internal/gen/generator/generator.go", generatorTemplate)

func (g *generatorGenerator) GenerateFile(fsys genfs.FS, file *genfs.File) error {
	g.log.Info("generating file", file.Path())
	type Generator struct {
		Import *imports.Import
		Path   string // Path that triggers the generator (e.g. "bud/cmd/app/main.go")
		Camel  string
		Type   string
	}
	type State struct {
		Imports    []*imports.Import
		Generators []*Generator
	}
	imset := imports.New()
	// imset.AddStd("fmt")
	imset.AddNamed("generator", "github.com/livebud/bud/runtime/generator")
	imset.AddNamed("gomod", "github.com/livebud/bud/package/gomod")
	imset.AddNamed("log", "github.com/livebud/bud/package/log")
	appImportPath := g.module.Import("generator/cmd/app")
	generators := []*Generator{
		{
			Import: &imports.Import{
				Name: imset.Add(appImportPath),
				Path: appImportPath,
			},
			Path:  "bud/cmd/app/main.go",
			Camel: "app",
			Type:  "Generator",
		},
	}

	code, err := generatorGen.Generate(State{
		Imports:    imset.List(),
		Generators: generators,
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
