package bud

import (
	"context"
	"fmt"

	"log/slog"

	"github.com/livebud/bud/pkg/cli"
	"github.com/livebud/bud/pkg/env"
	"github.com/livebud/bud/pkg/web"
)

func New(env *env.Bud, log *slog.Logger, server *web.Server) *Command {
	return &Command{env, log, server}
}

func Register(cli cli.Command, cmd *Command) {
	// Run is the root command
	run := &Run{}
	cli.Run(func(ctx context.Context) error {
		return cmd.Run(ctx, run)
	})
}

type Command struct {
	env    *env.Bud
	log    *slog.Logger
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
