package command

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/livebud/bud/cli"
	"github.com/livebud/bud/env"
	"github.com/livebud/bud/mux"
	"github.com/livebud/bud/pkg/web"
)

func New(env *env.Bud, log *slog.Logger, router *mux.Router, server *web.Server) *Bud {
	return &Bud{env, log, router, server}
}

func Register(cli cli.Command, cmd *Bud) {
	// Run is the root command
	run := &Run{}
	cli.Run(func(ctx context.Context) error {
		return cmd.Run(ctx, run)
	})
	routes := cli.Command("routes", "list the router routes")
	routes.Run(cmd.Routes)
}

type Bud struct {
	env    *env.Bud
	log    *slog.Logger
	router *mux.Router
	server *web.Server
}

// Run input
type Run struct {
	Listen string
}

func (c *Bud) Run(ctx context.Context, in *Run) error {
	ln, err := web.Listen(c.env.Listen)
	if err != nil {
		return err
	}
	c.log.Info(fmt.Sprintf("Listening on %s", web.Format(ln)))
	return c.server.Serve(ctx, ln)
}

func (c *Bud) Routes(ctx context.Context) error {
	for _, route := range c.router.List() {
		fmt.Println(route.String())
	}
	return nil
}
