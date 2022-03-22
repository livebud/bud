package run

import (
	"context"

	"gitlab.com/mnm/bud/pkg/socket"
	"gitlab.com/mnm/bud/runtime/project"
)

type Command struct {
	Project *project.Compiler
	Flag    project.Flag
	Port    string
}

func (c *Command) Run(ctx context.Context) error {
	app, err := c.Project.Compile(ctx, &c.Flag)
	if err != nil {
		return err
	}
	listener, err := socket.Load(c.Port)
	if err != nil {
		return err
	}
	process, err := app.Start(ctx, listener)
	if err != nil {
		return err
	}
	return process.Wait()
}
