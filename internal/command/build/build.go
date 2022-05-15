package build

import (
	"context"

	"github.com/livebud/bud/internal/bud"
	"github.com/livebud/bud/internal/command"
)

func New(bud *command.Bud) *Command {
	return &Command{bud}
}

type Command struct {
	bud *command.Bud
}

func (c *Command) Run(ctx context.Context) error {
	// Load the compiler
	compiler, err := bud.Find(c.bud.Dir)
	if err != nil {
		return err
	}
	// Compile the project CLI
	project, err := compiler.Compile(ctx, &c.bud.Flag)
	if err != nil {
		return err
	}
	// Build the project
	app, err := project.Build(ctx)
	if err != nil {
		return err
	}
	_ = app
	return nil
}
