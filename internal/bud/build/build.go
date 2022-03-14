package build

import (
	"context"

	"gitlab.com/mnm/bud/internal/bud"
	"gitlab.com/mnm/bud/package/trace"
	"gitlab.com/mnm/bud/pkg/gomod"
)

type Command struct {
	Bud *bud.Command
}

func (c *Command) Run(ctx context.Context) error {
	ctx, shutdown, err := c.Bud.Tracer(ctx)
	if err != nil {
		return err
	}
	defer shutdown(&err)
	_, span := trace.Start(ctx, "bud build")
	defer span.End(&err)
	// Find the project directory
	module, err := gomod.Find(c.Bud.Dir)
	if err != nil {
		return err
	}
	cli, err := c.Bud.Compile(ctx, module)
	if err != nil {
		return err
	}
	if err := cli.Build(ctx); err != nil {
		return err
	}
	return nil
}
