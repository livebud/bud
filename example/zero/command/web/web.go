package web

import (
	"context"
	"fmt"

	"github.com/livebud/bud/example/zero/bud/pkg/web"
)

func New(server *web.Server) *Command {
	return &Command{server}
}

type Command struct {
	server *web.Server
}

type Serve struct {
	Listen string
}

func (c *Command) GoServe(ctx context.Context, in *Serve) error {
	if in.Listen == "" {
		in.Listen = ":3000"
	}
	fmt.Println("starting the web routine!", in.Listen)
	return c.server.Listen(ctx, in.Listen)
}
