package build

import (
	"context"

	"gitlab.com/mnm/bud/internal/bud"
	"gitlab.com/mnm/bud/internal/command"
	"gitlab.com/mnm/bud/package/trace"
)

type Command struct {
	Bud *command.Bud
}

func (c *Command) Run(ctx context.Context) error {
	ctx, shutdown, err := c.Bud.Tracer(ctx)
	if err != nil {
		return err
	}
	defer shutdown(&err)
	_, span := trace.Start(ctx, "bud build")
	defer span.End(&err)
	// Load the compiler
	compiler, err := bud.Find(c.Bud.Dir)
	if err != nil {
		return err
	}
	// Compile the project CLI
	project, err := compiler.Compile(ctx, c.Bud.Flag)
	if err != nil {
		return err
	}
	// Build the project
	if err := project.Build(ctx); err != nil {
		return err
	}
	return nil
}
