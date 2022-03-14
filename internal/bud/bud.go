package bud

import (
	"context"
	"fmt"
	"os"
	"time"

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

type Command struct {
	Flag
	Trace bool
	Dir   string
	Args  []string
}

func (c *Command) Tracer(ctx context.Context) (context.Context, func(*error), error) {
	tracer, ctx, err := trace.Serve(ctx)
	if err != nil {
		return nil, nil, err
	}
	shutdown := func(outerError *error) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		trace, err := tracer.Print(ctx)
		if err != nil {
			*outerError = err
			return
		}
		fmt.Fprintf(os.Stderr, "\n%s", trace)
		if err := tracer.Shutdown(ctx); err != nil {
			*outerError = err
			return
		}
	}
	return ctx, shutdown, nil
}

// Load the module
func (c *Command) module(ctx context.Context) (module *gomod.Module, err error) {
	_, span := trace.Start(ctx, "find the module")
	defer span.End(&err)
	return module.Find(c.Dir)
}

// Load the overlay
func (c *Command) overlay(ctx context.Context, module *gomod.Module) (fsys *overlay.FileSystem, err error) {
	_, span := trace.Start(ctx, "load the overlay")
	defer span.End(&err)
	return overlay.Load(module)
}

// Sync the generators to bud/.cli
func (c *Command) sync(ctx context.Context, overlay *overlay.FileSystem) (err error) {
	_, span := trace.Start(ctx, "sync cli", "dir", "bud/.cli")
	defer span.End(&err)
	err = overlay.Sync("bud/.cli")
	return err
}

// Build the CLI
func (c *Command) gobuild(ctx context.Context, module *gomod.Module) (err error) {
	_, span := trace.Start(ctx, "build cli", "from", "bud/.cli/main.go", "to", "bud/cli")
	defer span.End(&err)
	return gobin.Build(ctx, module.Directory(), "bud/.cli/main.go", "bud/cli")
}

// Compile the project CLI
func (c *Command) Compile(ctx context.Context, module *gomod.Module) (cli *CLI, err error) {
	// Start the trace
	ctx, span := trace.Start(ctx, "compile project cli")
	defer span.End(&err)
	// Load the overlay
	overlay, err := c.overlay(ctx, module)
	if err != nil {
		return nil, err
	}
	// Initialize dependencies
	parser := parser.New(overlay, module)
	injector := di.New(overlay, module, parser)
	// Setup the generators
	overlay.FileGenerator("bud/.cli/main.go", mainfile.New(module))
	overlay.FileGenerator("bud/.cli/program/program.go", program.New(injector, module))
	overlay.FileGenerator("bud/.cli/command/command.go", command.New(overlay, module, parser))
	overlay.FileGenerator("bud/.cli/generator/generator.go", generator.New(overlay, module, parser))
	// Sync the generators
	if err := c.sync(ctx, overlay); err != nil {
		return nil, err
	}
	// Build the binary
	if err := c.gobuild(ctx, module); err != nil {
		return nil, err
	}
	return &CLI{c.Flag, module}, nil
}

// func (c *Command) Build(ctx context.Context, dir string) (string, error) {
// 	// // generator, err := generator.Load(dir)
// 	// // if err != nil {
// 	// // 	return "", err
// 	// // }
// 	// // if err := generator.Generate(ctx); err != nil {
// 	// // 	return "", err
// 	// // }
// 	// mainPath := filepath.Join(dir, "bud", "main.go")
// 	// // Check to see if we generated a main.go
// 	// if _, err := os.Stat(mainPath); err != nil {
// 	// 	return "", err
// 	// }
// 	// cacheDir, err := os.UserCacheDir()
// 	// if err != nil {
// 	// 	return "", err
// 	// }
// 	// // Building over an existing binary is faster for some reason, so we'll use
// 	// // the cache directory for a consistent place to output builds
// 	// binPath := filepath.Join(cacheDir, filepath.ToSlash(generator.Module().Import()), "bud", "main")
// 	// if err := gobin.Build(ctx, dir, mainPath, binPath); err != nil {
// 	// 	return "", err
// 	// }
// 	// return binPath, nil
// 	panic("Not implemented")
// }

// Run a custom command
func (c *Command) Run(ctx context.Context) (err error) {
	// Start the trace
	ctx, span := trace.Start(ctx, "running bud")
	defer span.End(&err)
	// Find the module
	module, err := c.module(ctx)
	if err != nil {
		return err
	}
	// Compile the project CLI
	cli, err := c.Compile(ctx, module)
	if err != nil {
		return err
	}
	// Run the custom command
	return cli.Custom(ctx, c.Args...)
	// if err := c.Compile(ctx, module); err != nil {
	// 	return err
	// }
	// // Find the project directory
	// dir, err := gomod.Absolute(c.Chdir)
	// if err != nil {
	// 	return err
	// }
	// // Generate the code
	// binPath, err := c.Build(ctx, dir)
	// if err != nil {
	// 	if !errors.Is(err, fs.ErrNotExist) {
	// 		return err
	// 	}
	// 	return fmt.Errorf("unknown command %q", c.Args)
	// }
	// Run the built binary
	// cmd := exec.Command("bud/cli", c.Args...)
	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout
	// err = cmd.Run()
	// if err != nil {
	// 	return err
	// }
	// return nil
}
