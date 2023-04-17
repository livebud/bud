package web

import (
	"context"

	"github.com/livebud/bud/example/zero/bud/pkg/web"
	"github.com/livebud/bud/example/zero/env"
	"github.com/livebud/bud/package/log"
)

func New(env *env.Env, log log.Log, server *web.Server) *Command {
	return &Command{env, log, server}
}

type Command struct {
	env    *env.Env
	log    log.Log
	server *web.Server
}

type Serve struct {
	Listen string
}

func (c *Command) GoServe(ctx context.Context, in *Serve) error {
	if in.Listen == "" {
		in.Listen = ":3000"
	}
	c.log.Infof("connecting to database: %s", c.env.Database.URL)
	c.log.Infof("starting the web server on: http://localhost%s", in.Listen)
	return c.server.Listen(ctx, in.Listen)
}
