package client

import (
	"context"

	"github.com/livebud/bud/package/js/v8server"
)

type Command struct {
}

func (c *Command) Run(ctx context.Context) error {
	return v8server.Serve()
}
