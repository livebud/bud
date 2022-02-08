package buddy

import (
	"context"
	"io/fs"

	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gen"

	"gitlab.com/mnm/bud/pkg/parser"
	"gitlab.com/mnm/bud/pkg/pluginfs"

	"gitlab.com/mnm/bud/internal/buddy/build"
	"gitlab.com/mnm/bud/internal/buddy/expand"
	"gitlab.com/mnm/bud/internal/buddy/generate"
	"gitlab.com/mnm/bud/internal/buddy/run"
	"gitlab.com/mnm/bud/pkg/gomod"
)

type Option func(d *Driver)

// Load the driver from a directory
func Load(dir string, options ...Option) (*Driver, error) {
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	pluginFS, err := pluginfs.Load(module)
	if err != nil {
		return nil, err
	}
	genFS := gen.New(pluginFS)
	parser := parser.New(genFS, module)
	injector := di.New(genFS, module, parser)
	return &Driver{
		mod:       module,
		gen:       genFS,
		parser:    parser,
		injector:  injector,
		expander:  expand.New(genFS, injector, module, parser),
		generator: generate.New(module),
		builder:   build.New(module),
		runner:    run.New(module),
	}, nil
}

// Driver is a single, public entrypoint that generators and commands can use to
// extend Bud. The driver itself, should not do much, rather it should delegate
// to various internal implementations.
type Driver struct {
	mod       *gomod.Module
	gen       *gen.FileSystem
	injector  *di.Injector
	expander  *expand.Command
	generator *generate.Command
	builder   *build.Command
	runner    *run.Command
	parser    *parser.Parser
}

// Expand input
type Expand = expand.Input

// Expand commands and user-defined generators and generate a "project CLI"
func (d *Driver) Expand(ctx context.Context, in *Expand) error {
	return d.expander.Expand(ctx, in)
}

// Generate bud files from the project CLI. Depends on Expand. Generate does not
// run go build on the files.
func (d *Driver) Generate(ctx context.Context, options ...generate.Option) error {
	return d.generator.Generate(ctx, options...)
}

// Build an application from the generated files. Depends on Generate.
func (d *Driver) Build(ctx context.Context, options ...build.Option) error {
	return d.builder.Build(ctx, options...)
}

// Run an application from the generated files and watch for changes.
// Depends on Generate.
func (d *Driver) Run(ctx context.Context, options ...run.Option) error {
	return d.runner.Run(ctx, options...)
}

// ImportPath returns an import path within the application module.
func (d *Driver) ImportPath(subpaths ...string) string {
	return d.mod.Import(subpaths...)
}

// Parse a Go package
func (d *Driver) Parse(dir string) (*parser.Package, error) {
	return d.parser.Parse(dir)
}

type Function = di.Function
type Provider = di.Provider

// Wire up a function
func (d *Driver) Wire(fn *Function) (*Provider, error) {
	return d.injector.Wire(fn)
}

// Generator adds a new generator
func (d *Driver) Generator(path string, generator gen.Generator) {
	d.gen.Add(map[string]gen.Generator{path: generator})
}

// Open a file. Implements fs.FS. Open is looped over to generate bud files.
func (d *Driver) Open(name string) (fs.File, error) {
	return d.gen.Open(name)
}
