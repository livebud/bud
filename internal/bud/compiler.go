package bud

import (
	"context"
	"os"

	"gitlab.com/mnm/bud/generator/cli/command"
	"gitlab.com/mnm/bud/generator/cli/generator"
	"gitlab.com/mnm/bud/generator/cli/mainfile"
	"gitlab.com/mnm/bud/generator/cli/program"
	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/package/trace"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

func Find(dir string) (*Compiler, error) {
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	return &Compiler{module}, nil
}

func New(module *gomod.Module) *Compiler {
	return &Compiler{module}
}

type Compiler struct {
	module *gomod.Module
}

// Load the module
func (c *Compiler) findModule(ctx context.Context, dir string) (module *gomod.Module, err error) {
	_, span := trace.Start(ctx, "find the module")
	defer span.End(&err)
	return module.Find(dir)
}

// Load the overlay
func (c *Compiler) loadOverlay(ctx context.Context, module *gomod.Module) (fsys *overlay.FileSystem, err error) {
	_, span := trace.Start(ctx, "load the overlay")
	defer span.End(&err)
	return overlay.Load(module)
}

// Sync the generators to bud/.cli
func (c *Compiler) sync(ctx context.Context, overlay *overlay.FileSystem) (err error) {
	_, span := trace.Start(ctx, "sync cli", "dir", "bud/.cli")
	defer span.End(&err)
	err = overlay.Sync("bud/.cli")
	return err
}

// Build the CLI
func (c *Compiler) goBuild(ctx context.Context, module *gomod.Module) (err error) {
	_, span := trace.Start(ctx, "build cli", "from", "bud/.cli/main.go", "to", "bud/cli")
	defer span.End(&err)
	return gobin.Build(ctx, module.Directory(), "bud/.cli/main.go", "bud/cli")
}

func (c *Compiler) Compile(ctx context.Context, flag Flag) (p *Project, err error) {
	// Start the trace
	ctx, span := trace.Start(ctx, "compile project cli")
	defer span.End(&err)
	// Load the overlay
	overlay, err := c.loadOverlay(ctx, c.module)
	if err != nil {
		return nil, err
	}
	// Initialize dependencies
	parser := parser.New(overlay, c.module)
	injector := di.New(overlay, c.module, parser)
	// Setup the generators
	overlay.FileGenerator("bud/.cli/main.go", mainfile.New(c.module))
	overlay.FileGenerator("bud/.cli/program/program.go", program.New(injector, c.module))
	overlay.FileGenerator("bud/.cli/command/command.go", command.New(overlay, c.module, parser))
	overlay.FileGenerator("bud/.cli/generator/generator.go", generator.New(overlay, c.module, parser))
	// Sync the generators
	if err := c.sync(ctx, overlay); err != nil {
		return nil, err
	}
	// Build the binary
	if err := c.goBuild(ctx, c.module); err != nil {
		return nil, err
	}
	return &Project{
		Module: c.module,
		Flag:   flag,
		Env: map[string]string{
			"HOME":   os.Getenv("HOME"),
			"PATH":   os.Getenv("PATH"),
			"GOPATH": os.Getenv("GOPATH"),
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, nil
}
