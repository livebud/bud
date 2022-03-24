package client

import (
	"context"

	"gitlab.com/mnm/bud/internal/command"
	v8client "gitlab.com/mnm/bud/package/js/v8client"
)

type Command struct {
	Bud *command.Bud
}

func (c *Command) Run(ctx context.Context) error {
	return v8client.Serve(ctx)
}
