package cli

import (
	"context"

	"gitlab.com/mnm/bud/framework2/cli/command"
	"gitlab.com/mnm/bud/framework2/cli/mainfile"
	"gitlab.com/mnm/bud/framework2/cli/program"
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
	ctx, span := trace.Start(ctx, "cli compile")
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
	overlay.FileGenerator("bud/.cli/main.go", mainfile.New(c.module))
	overlay.FileGenerator("bud/.cli/program/program.go", program.New(injector, c.module))
	overlay.FileGenerator("bud/.cli/command/command.go", command.New(overlay, c.module, parser))
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
	err = overlay.Sync("bud/.cli")
	return err
}

// Build the CLI
func (c *Compiler) build(ctx context.Context) (err error) {
	_, span := trace.Start(ctx, "build cli", "from", "bud/.cli/main.go", "to", "bud/cli")
	defer span.End(&err)
	return gobin.Build(ctx, c.module.Directory(), "bud/.cli/main.go", "bud/cli")
}
