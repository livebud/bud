package expand

import (
	"context"

	"gitlab.com/mnm/bud/generator2/cli/command"
	"gitlab.com/mnm/bud/generator2/cli/generator"
	"gitlab.com/mnm/bud/generator2/cli/mainfile"
	"gitlab.com/mnm/bud/generator2/cli/program"
	"gitlab.com/mnm/bud/internal/gobin"
	commandparser "gitlab.com/mnm/bud/internal/parser/command"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/package/trace"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

func Load(ctx context.Context, dir string) (cmd *Command, err error) {
	_, span := trace.Start(ctx, "load expander", "dir", dir)
	defer span.End(&err)
	module, err := gomod.Find(dir)
	if err != nil {
		return nil, err
	}
	ofs, err := overlay.Load(module)
	if err != nil {
		return nil, err
	}
	parser := parser.New(ofs, module)
	injector := di.New(ofs, module, parser)
	ofs.FileGenerator("bud/.cli/main.go", mainfile.New(module))
	ofs.FileGenerator("bud/.cli/program/program.go", program.New(injector, module))
	ofs.FileGenerator("bud/.cli/command/command.go", command.New(commandparser.New(module, parser)))
	ofs.FileGenerator("bud/.cli/generator/generator.go", generator.New())
	return &Command{
		dir:     dir,
		overlay: ofs,
		module:  module,
	}, nil
}

type Command struct {
	dir     string
	overlay *overlay.FileSystem
	module  *gomod.Module
}

func (c *Command) Run(ctx context.Context) (err error) {
	ctx, span := trace.Start(ctx, "run expander", "dir", c.dir)
	defer span.End(&err)
	if err := c.sync(ctx); err != nil {
		return err
	}
	if err := c.build(ctx); err != nil {
		return err
	}
	return nil
}

// Sync the generators to bud/.cli
func (c *Command) sync(ctx context.Context) (err error) {
	_, span := trace.Start(ctx, "sync cli", "dir", "bud/.cli")
	defer span.End(&err)
	return c.overlay.Sync("bud/.cli")
}

// Build the CLI
func (c *Command) build(ctx context.Context) (err error) {
	_, span := trace.Start(ctx, "build cli", "from", "bud/.cli/main.go", "to", "bud/cli")
	defer span.End(&err)
	return gobin.Build(ctx, c.dir, "bud/.cli/main.go", "bud/cli")
}
