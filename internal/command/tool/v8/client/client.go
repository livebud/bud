package client

import (
	"context"

	"gitlab.com/mnm/bud/internal/command"
	"gitlab.com/mnm/bud/package/js/v8server"
)

type Command struct {
	Bud *command.Bud
}

func (c *Command) Run(ctx context.Context) error {
	return v8server.Serve(ctx)
}
