package buddy

import (
	"context"
	"io/fs"

	"gitlab.com/mnm/bud/internal/dsync"
	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gen"

	"gitlab.com/mnm/bud/pkg/parser"
	"gitlab.com/mnm/bud/pkg/pluginfs"

	"gitlab.com/mnm/bud/pkg/gomod"
)

type Generator = gen.Generator

type Kit interface {
	// DirFS(subpaths ...string) vfs.ReadWritable
	ImportPath(subpaths ...string) string
	Parse(dir string) (*parser.Package, error)
	Wire(fn *Function) (*Provider, error)
	Generators(generators map[string]Generator) error
	Open(name string) (fs.File, error)
	Sync(from, to string) error
	Go() Go
}

type Go interface {
	Build(ctx context.Context, mainPath, outPath string) error
}

type Option func(*Kit)

// Load the driver from a directory
func Load(dir string, options ...Option) (Kit, error) {
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
	return &kit{
		mod:      module,
		gen:      genFS,
		parser:   parser,
		injector: injector,
		golang:   &golang{module},
	}, nil
}

type kit struct {
	mod      *gomod.Module
	gen      *gen.FileSystem
	injector *di.Injector
	parser   *parser.Parser
	golang   *golang
}

// ImportPath returns an import path within the application module.
func (k *kit) ImportPath(subpaths ...string) string {
	return k.mod.Import(subpaths...)
}

// Parse a Go package
func (k *kit) Parse(dir string) (*parser.Package, error) {
	return k.parser.Parse(dir)
}

// TODO: this is just going to get worse over time.
type Dependency = di.Dependency
type Function = di.Function
type Provider = di.Provider
type ErrorType = di.Error

var ToType = di.ToType

// Wire up a function
func (k *kit) Wire(fn *Function) (*Provider, error) {
	return k.injector.Wire(fn)
}

// Generator adds a new generator
func (k *kit) Generators(generators map[string]Generator) error {
	k.gen.Add(generators)
	return nil
}

// Open a file. Implements fs.FS. Open is looped over to generate bud files.
func (k *kit) Open(name string) (fs.File, error) {
	return k.gen.Open(name)
}

// Sync the generators with the filesystem.
func (k *kit) Sync(from, to string) error {
	return dsync.Dir(k.gen, from, k.mod.DirFS(to), ".")
}

// Get the Go commands.
func (k *kit) Go() Go {
	return k.golang
}

type golang struct {
	mod *gomod.Module
}

// Run `go build`
func (g *golang) Build(ctx context.Context, mainPath, outPath string) error {
	return gobin.Build(ctx, g.mod.Directory(), mainPath, outPath)
}
