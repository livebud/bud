package cmds

import (
	"context"
	"net/http"

	"github.com/livebud/bud/pkg/cli"
	"github.com/livebud/bud/pkg/logs"
	"github.com/livebud/bud/pkg/mux"
)

func Serve(log logs.Log, router *mux.Router) *serveCmd {
	return &serveCmd{log, router}
}

type serveCmd struct {
	log    logs.Log
	router *mux.Router
}

func (c *serveCmd) Usage(cmd cli.Command) {
}

func (c *serveCmd) Run(ctx context.Context) error {
	c.log.Infof("Listening on http://localhost%s", ":8080")
	return http.ListenAndServe(":8080", c.router)
}
