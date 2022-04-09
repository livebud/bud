package build

import (
	"context"

	"gitlab.com/mnm/bud/internal/bud"
	"gitlab.com/mnm/bud/internal/command"
)

type Command struct {
	Bud *command.Bud
}

func (c *Command) Run(ctx context.Context) error {
	// Load the compiler
	compiler, err := bud.Find(c.Bud.Dir)
	if err != nil {
		return err
	}
	// Compile the project CLI
	project, err := compiler.Compile(ctx, &c.Bud.Flag)
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
