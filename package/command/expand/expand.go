package expand

import (
	"context"

	"gitlab.com/mnm/bud/generator2/cli/command"
	"gitlab.com/mnm/bud/generator2/cli/generator"
	"gitlab.com/mnm/bud/generator2/cli/mainfile"
	"gitlab.com/mnm/bud/generator2/cli/program"
	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

func Load(dir string) (*Command, error) {
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
	ofs.FileGenerator("bud/.cli/command/command.go", command.New(module))
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

func (c *Command) Run(ctx context.Context) error {
	if err := c.overlay.Sync("bud/.cli"); err != nil {
		return err
	}
	// Build the CLI
	if err := gobin.Build(ctx, c.dir, "bud/.cli/main.go", "bud/cli"); err != nil {
		return err
	}
	return nil
}
