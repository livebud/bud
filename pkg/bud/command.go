package bud

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/livebud/bud/pkg/cli"
	"github.com/livebud/bud/pkg/env"
	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/web"
)

func New(env *env.Bud, log *slog.Logger, router *mux.Router, server *web.Server) *Command {
	return &Command{env, log, router, server}
}

func Register(cli cli.Command, cmd *Command) {
	// Run is the root command
	run := &Run{}
	cli.Run(func(ctx context.Context) error {
		return cmd.Run(ctx, run)
	})
	routes := cli.Command("routes", "list the router routes")
	routes.Run(cmd.Routes)
}

type Command struct {
	env    *env.Bud
	log    *slog.Logger
	router *mux.Router
	server *web.Server
}

// Run input
type Run struct {
	Listen string
}

func (c *Command) Run(ctx context.Context, in *Run) error {
	ln, err := web.Listen(c.env.Listen)
	if err != nil {
		return err
	}
	c.log.Info(fmt.Sprintf("Listening on %s", web.Format(ln)))
	return c.server.Serve(ctx, ln)
}

func (c *Command) Routes(ctx context.Context) error {
	for _, route := range c.router.List() {
		fmt.Println(route.String())
	}
	return nil
}
