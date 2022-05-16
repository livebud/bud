package client

import (
	"context"

	"github.com/livebud/bud/package/js/v8server"
)

func New() *Command {
	return &Command{}
}

type Command struct {
}

func (c *Command) Run(ctx context.Context) error {
	return v8server.Serve()
}
