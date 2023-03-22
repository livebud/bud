package cli

import (
	"context"

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
	// cache, err := dag.Load(log, module.Directory("bud", "bud.db"))
	// if err != nil {
	// 	return err
	// }
	// defer cache.Close()
	gen := genfs.New(dag.Discard, module, log)
	parser := parser.New(gen, module)
	injector := di.New(gen, log, module, parser)
	gen.FileGenerator("bud/cmd/gen/main.go", &mainGenerator{injector, log, module})
	gen.FileGenerator("bud/internal/generator/generator.go", &generatorGenerator{log, module})
	gen.FileGenerator("bud/pkg/transpiler/transpiler.go", &transpilerGenerator{log, module})
	gen.FileGenerator("bud/pkg/viewer/viewer.go", &viewerGenerator{log, module})
	if err := virtual.Sync(log, gen, module, "bud"); err != nil {
		return err
	}

	// Build bud/gen
	cmd := c.command(module.Directory(), "go", "build", "-mod=mod", "-o=bud/gen", "./bud/cmd/gen")
	if err := cmd.Run(); err != nil {
		return err
	}

	// Run bud/gen
	cmd = c.command(module.Directory(), "./bud/gen", in.Flag.Flags()...)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Run bud/app
	// TODO: this should be moved into `bud run`
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
		Params: []*di.Param{
			&di.Param{
				Import: "github.com/livebud/bud/framework",
				Type:   "*Flag",
			},
			&di.Param{
				Import: "github.com/livebud/bud/package/gomod",
				Type:   "*Module",
			},
			&di.Param{
				Import: "github.com/livebud/bud/package/log",
				Type:   "Log",
			},
			&di.Param{
				Import: "github.com/livebud/bud/package/genfs",
				Type:   "FileSystem",
			},
		},
		Aliases: di.Aliases{
			di.ToType("github.com/livebud/bud/package/parser", "*Parser"): di.ToType("github.com/livebud/bud/runtime/gen", "*Parser"),
			di.ToType("github.com/livebud/bud/package/di", "*Injector"):   di.ToType("github.com/livebud/bud/runtime/gen", "*Injector"),
		},
		Results: []di.Dependency{
			di.ToType(g.module.Import("bud/internal/generator"), "*Generator"),
			&di.Error{},
		},
	})
	if err != nil {
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

func NewGenerator(
	genfs generator.FileSystem,
	log log.Log,
	{{- range $generator := $.Generators }}
	{{ $generator.Camel }} *{{ $generator.Import.Name }}.{{ $generator.Type }},
	{{- end }}
) *Generator {
	return generator.NewGenerator(
		genfs,
		log,
		{{- range $generator := $.Generators }}
		{{ $generator.Camel }},
		{{- end }}
	)
}

type Generator = generator.Generator
`

var generatorGen = gotemplate.MustParse("generator.go", generatorTemplate)

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
	// imset.AddNamed("gomod", "github.com/livebud/bud/package/gomod")
	imset.AddNamed("log", "github.com/livebud/bud/package/log")
	appImportPath := g.module.Import("generator/app")
	commandImportPath := g.module.Import("generator/command")
	webImportPath := g.module.Import("generator/web")
	controllerImportPath := g.module.Import("generator/controller")
	viewImportPath := g.module.Import("generator/view")
	generators := []*Generator{
		{
			Import: &imports.Import{
				Name: imset.Add(appImportPath),
				Path: appImportPath,
			},
			Path:  "bud",
			Camel: "app",
			Type:  "Generator",
		},
		{
			Import: &imports.Import{
				Name: imset.Add(commandImportPath),
				Path: commandImportPath,
			},
			Path:  "bud",
			Camel: "command",
			Type:  "Generator",
		},
		{
			Import: &imports.Import{
				Name: imset.Add(webImportPath),
				Path: webImportPath,
			},
			Path:  "bud",
			Camel: "web",
			Type:  "Generator",
		},
		{
			Import: &imports.Import{
				Name: imset.Add(controllerImportPath),
				Path: controllerImportPath,
			},
			Path:  "bud",
			Camel: "controller",
			Type:  "Generator",
		},
		{
			Import: &imports.Import{
				Name: imset.Add(viewImportPath),
				Path: viewImportPath,
			},
			Path:  "bud",
			Camel: "view",
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

type transpilerGenerator struct {
	log    log.Log
	module *gomod.Module
}

const transpilerTemplate = `package transpiler

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{$import.Name}} "{{$import.Path}}"
	{{- end }}
)
{{- end }}

// Load the transpiler
func Load(
	{{- range $transpiler := $.Transpilers }}
	{{ $transpiler.Camel }} *{{ $transpiler.Import.Name }}.Transpiler,
	{{- end }}
) Transpiler {
	tr := transpiler.New()
	{{- range $transpiler := $.Transpilers }}
	{{- range $method := $transpiler.Methods }}
	tr.Add("{{ $method.From }}", "{{ $method.To }}", {{ $transpiler.Camel }}.{{ $method.Pascal }})
	{{- end }}
	{{- end }}
	return &proxy{tr}
}

type Transpiler = transpiler.Interface

type proxy struct {
	Transpiler
}

func (p *proxy) Transpile(fromExt, toExt string, code []byte) ([]byte, error) {
	transpiled, err := p.Transpiler.Transpile(fromExt, toExt, code)
	if err != nil {
		if !errors.Is(err, transpiler.ErrNoPath) {
			return nil, err
		}
		return code, nil
	}
	return transpiled, nil
}
`

var transpilerGen = gotemplate.MustParse("transpiler.go", transpilerTemplate)

func (g *transpilerGenerator) GenerateFile(fsys genfs.FS, file *genfs.File) error {
	g.log.Info("generating file", file.Path())
	type Method struct {
		Pascal string // Method name in pascal
		From   string // From extension
		To     string // To extension
	}
	type Transpiler struct {
		Import  *imports.Import
		Camel   string
		Methods []*Method
	}
	type State struct {
		Imports     []*imports.Import
		Transpilers []*Transpiler
	}
	imset := imports.New()
	imset.AddStd("errors")
	imset.Add("github.com/livebud/bud/runtime/transpiler")
	tailwindImport := g.module.Import("transpiler/tailwind")
	goldmarkImport := g.module.Import("transpiler/goldmark")
	code, err := transpilerGen.Generate(State{
		Transpilers: []*Transpiler{
			{
				Import: &imports.Import{
					Name: imset.Add(tailwindImport),
					Path: tailwindImport,
				},
				Camel: "tailwind",
				Methods: []*Method{
					{
						Pascal: "GohtmlToGohtml",
						From:   ".gohtml",
						To:     ".gohtml",
					},
				},
			},
			{
				Import: &imports.Import{
					Name: imset.Add(goldmarkImport),
					Path: goldmarkImport,
				},
				Camel: "goldmark",
				Methods: []*Method{
					{
						Pascal: "MdToGohtml",
						From:   ".md",
						To:     ".gohtml",
					},
				},
			},
		},
		Imports: imset.List(),
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}

type viewerGenerator struct {
	log    log.Log
	module *gomod.Module
}

const viewerTemplate = `package viewer

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{$import.Name}} "{{$import.Path}}"
	{{- end }}
)
{{- end }}

// Load the viewer
func New(
	{{- range $viewer := $.Viewers }}
	{{ $viewer.Camel }} *{{ $viewer.Import.Name }}.Viewer,
	{{- end }}
) Map {
	return Map{
		{{- range $viewer := $.Viewers }}
		"{{ $viewer.Ext }}": {{ $viewer.Camel }},
		{{- end }}
	}
}

type Map map[string]view.Viewer

// var _ view.Viewer = Map{}

func (viewers Map) Register(r *router.Router, pages []*view.Page) {
	for _, viewer := range viewers {
		viewer.Register(r, pages)
	}
}

// func (viewers Map) Render(ctx context.Context, page *view.Page, propMap view.PropMap) ([]byte, error) {
// 	viewer, ok := viewers[page.Ext]
// 	if !ok {
// 		return nil, fmt.Errorf("viewer: no viewer for extension %q to render %s", page.Ext, page.Path)
// 	}
// 	return viewer.Render(ctx, page, propMap)
// }

// func (viewers Map) RenderError(ctx context.Context, page *view.Page, propMap view.PropMap, err error) []byte {
// 	viewer, ok := viewers[page.Ext]
// 	if !ok {
// 		return []byte(fmt.Sprintf("viewer: no viewer for extension %q to render error %s", page.Ext, err))
// 	}
// 	return viewer.RenderError(ctx, page, propMap, err)
// }

func (viewers Map) Bundle(ctx context.Context, fsys view.Writable, pages []*view.Page) error {
	return fmt.Errorf("viewer: bundle not implemented")
}
`

var viewerGen = gotemplate.MustParse("viewer.go", viewerTemplate)

func (g *viewerGenerator) GenerateFile(fsys genfs.FS, file *genfs.File) error {
	g.log.Info("generating viewer", file.Path())
	type Viewer struct {
		Import *imports.Import
		Ext    string
		Camel  string
		Pascal string
	}
	type State struct {
		Imports []*imports.Import
		Viewers []*Viewer
	}
	imset := imports.New()
	imset.AddStd("context", "fmt")
	imset.Add("github.com/livebud/bud/runtime/view")
	imset.Add("github.com/livebud/bud/package/router")
	gohtmlPath := g.module.Import("viewer/gohtml")
	code, err := viewerGen.Generate(State{
		Viewers: []*Viewer{
			{
				Import: &imports.Import{
					Name: imset.Add(gohtmlPath),
					Path: gohtmlPath,
				},
				Ext:    ".gohtml",
				Camel:  "gohtml",
				Pascal: "Gohtml",
			},
		},
		Imports: imset.List(),
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
