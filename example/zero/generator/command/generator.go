package command

import (
	"context"

	"github.com/livebud/bud/package/commander"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/runtime/generator"
	"golang.org/x/sync/errgroup"
)

func NewGenerator(module *gomod.Module) *Generator {
	return &Generator{module}
}

type Generator struct {
	module *gomod.Module
}

const template = `package command

{{- if $.Imports }}

import (
	{{- range $import := $.Imports }}
	{{$import.Name}} "{{$import.Path}}"
	{{- end }}
)
{{- end }}

func New(
	log log.Log,
	webCmd *web.Command,
) *CLI {
	cli := command.New("app")
	webIn := new(web.Serve)
	cli.Run(func(ctx context.Context) error {
		return command.Go(ctx, log,
			func(ctx context.Context) error { return webCmd.GoServe(ctx, webIn) },
		)
	})

	{ // web

		{ // web:serve
			cmd := cli.Command("web:serve", "serve web requests")
			in := new(web.Serve)
			cmd.Run(func(ctx context.Context) error {
				return webCmd.GoServe(ctx, in)
			})
		}
	}

	return cli
}

type CLI = command.CLI
`

var gen = gotemplate.MustParse("command.gotext", template)

func (g *Generator) Extend(gen generator.FileSystem) {
	// TODO: should bud/ be implied? I don't think we should sync non-bud/
	// directories, it's too risky.
	gen.GenerateFile("bud/internal/command/command.go", g.generateFile)
}

type Routine struct{}

type State struct {
	Imports []*imports.Import
	Name    string
}

func (g *Generator) generateFile(fsys generator.FS, file *generator.File) error {
	imset := imports.New()
	imset.AddStd("context")
	imset.Add("github.com/livebud/bud/package/log")
	// imset.Add("golang.org/x/sync/errgroup")
	imset.Add(g.module.Import("generator/command"))
	// TODO: move to State
	imset.Add(g.module.Import("command/web"))
	code, err := gen.Generate(&State{
		Imports: imset.List(),
		Name:    "app",
	})
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}

////////////////////////////////////////////////
// Runtime from the code generated above
////////////////////////////////////////////////

func New(name string) *CLI {
	return commander.New(name)
}

// Go starts a group of commands in goroutines and waits for them to finish
func Go(ctx context.Context, log log.Log, fns ...func(context.Context) error) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, fn := range fns {
		fn := fn
		eg.Go(func() error {
			return fn(ctx)
		})
	}
	return eg.Wait()
}

type CLI = commander.CLI
