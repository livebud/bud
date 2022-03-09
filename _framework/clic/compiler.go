package clic

import (
	"context"

	"gitlab.com/mnm/bud/framework/maing"
	"gitlab.com/mnm/bud/framework/programg"
	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/package/trace"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

// New CLI compiler
func New(module *gomod.Module) *Compiler {
	return &Compiler{module}
}

// Compiler for the project CLI
type Compiler struct {
	module *gomod.Module
}

func (c *Compiler) Compile(ctx context.Context) (err error) {
	// Start the trace
	ctx, span := trace.Start(ctx, "clic compile")
	defer span.End(&err)
	// Load the overlay
	overlay, err := c.overlay(ctx)
	if err != nil {
		return err
	}
	// Initialize dependencies
	parser := parser.New(overlay, c.module)
	injector := di.New(overlay, c.module, parser)
	// Setup the generators
	overlay.GenerateFile("bud/.cli/main.go", maing.New(c.module.Import("bud/.cli/program")))
	overlay.GenerateFile("bud/.cli/program/program.go", programg.New(injector, c.module, &di.Function{
		Name:   "loadCLI",
		Target: c.module.Import("bud/.cli/program"),
		Params: []di.Dependency{
			di.ToType("gitlab.com/mnm/bud/pkg/di", "*Injector"),
			di.ToType("gitlab.com/mnm/bud/pkg/gomod", "*Module"),
			di.ToType("gitlab.com/mnm/bud/package/overlay", "*FileSystem"),
			di.ToType("gitlab.com/mnm/bud/pkg/parser", "*Parser"),
		},
		Results: []di.Dependency{
			di.ToType(c.module.Import("bud/.cli/command"), "*CLI"),
			&di.Error{},
		},
	}))
	// overlay.FileGenerator("bud/.cli/command/command.go", command.New(commandparser.New(module, parser)))
	// overlay.FileGenerator("bud/.cli/generator/generator.go", generator.New())
	// Sync the generators
	if err := c.sync(ctx, overlay); err != nil {
		return err
	}
	// Build the binary
	if err := c.build(ctx); err != nil {
		return err
	}
	return nil
}

// Load the overlay
func (c *Compiler) overlay(ctx context.Context) (fsys *overlay.FileSystem, err error) {
	_, span := trace.Start(ctx, "load the overlay")
	defer span.End(&err)
	return overlay.Load(c.module)
}

// Sync the generators to bud/.cli
func (c *Compiler) sync(ctx context.Context, overlay *overlay.FileSystem) (err error) {
	_, span := trace.Start(ctx, "sync cli", "dir", "bud/.cli")
	defer span.End(&err)
	return overlay.Sync("bud/.cli")
}

// Build the CLI
func (c *Compiler) build(ctx context.Context) (err error) {
	_, span := trace.Start(ctx, "build cli", "from", "bud/.cli/main.go", "to", "bud/cli")
	defer span.End(&err)
	return gobin.Build(ctx, c.module.Directory(), "bud/.cli/main.go", "bud/cli")
}
