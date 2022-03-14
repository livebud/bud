package buddy

import (
	"context"
	"fmt"
)

// func New(dir string) *Compiler {
// 	return &Compiler{dir}
// }

// type Compiler struct {
// 	dir string
// }

type compileOption interface {
	compile(o *compileConfig)
}

type compileConfig struct {
	Embed  bool
	Minify bool
	Hot    bool
}

// func (c *Compiler) Compile(ctx context.Context, options ...compileOption) (cli *CLI, err error) {
// 	// Start the trace
// 	ctx, span := trace.Start(ctx, "buddy compile")
// 	defer span.End(&err)
// 	// Load the overlay
// 	overlay, err := c.overlay(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Initialize dependencies
// 	parser := parser.New(overlay, c.module)
// 	injector := di.New(overlay, c.module, parser)
// 	// Setup the generators
// 	overlay.FileGenerator("bud/.cli/main.go", mainfile.New(c.module))
// 	overlay.FileGenerator("bud/.cli/program/program.go", program.New(injector, c.module))
// 	overlay.FileGenerator("bud/.cli/command/command.go", command.New(overlay, c.module, parser))
// 	// overlay.FileGenerator("bud/.cli/generator/generator.go", generator.New())
// 	// Sync the generators
// 	if err := c.sync(ctx, overlay); err != nil {
// 		return nil, err
// 	}
// 	// Build the binary
// 	if err := c.build(ctx); err != nil {
// 		return nil, err
// 	}
// 	return &CLI{
// 		module: c.module,
// 		path:   "bud/cli",
// 	}, nil
// }

// // Load the overlay
// func (c *Compiler) overlay(ctx context.Context) (fsys *overlay.FileSystem, err error) {
// 	_, span := trace.Start(ctx, "load the overlay")
// 	defer span.End(&err)
// 	return overlay.Load(c.module)
// }

// // Sync the generators to bud/.cli
// func (c *Compiler) sync(ctx context.Context, overlay *overlay.FileSystem) (err error) {
// 	_, span := trace.Start(ctx, "sync cli", "dir", "bud/.cli")
// 	defer span.End(&err)
// 	err = overlay.Sync("bud/.cli")
// 	return err
// }

// // Build the CLI
// func (c *Compiler) build(ctx context.Context) (err error) {
// 	_, span := trace.Start(ctx, "build cli", "from", "bud/.cli/main.go", "to", "bud/cli")
// 	defer span.End(&err)
// 	return gobin.Build(ctx, c.module.Directory(), "bud/.cli/main.go", "bud/cli")
// }

type runOption interface {
	run(o *runConfig)
}

type runConfig struct {
	Embed  bool
	Minify bool
	Hot    bool
	Port   string
}

func (c *Compiler) Run(ctx context.Context, options ...runOption) (*Process, error) {
	return nil, fmt.Errorf("not implemented")
}

type buildOption interface {
	build(o *buildConfig)
}

type buildConfig struct {
	Embed  bool
	Minify bool
}

func (c *Compiler) Build(ctx context.Context, options ...buildOption) (*App, error) {
	cfg := &buildConfig{
		Embed:  true,
		Minify: true,
	}
	for _, option := range options {
		option.build(cfg)
	}
	return nil, fmt.Errorf("not finished")
}
