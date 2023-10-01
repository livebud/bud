package serve

import (
	"context"
	"fmt"
	"net"

	"github.com/livebud/bud/cli"
	"github.com/livebud/bud/internal/socket"
	"github.com/livebud/bud/log"
	"github.com/livebud/bud/web"
)

func New(log log.Log, server web.Service) *Command {
	return &Command{":8080", log, server}
}

type Command struct {
	Listen string
	log    log.Log
	server web.Service
}

func (c *Command) Mount(cli cli.Command) {
	serve := cli.Command("serve", "serve the app")
	serve.Flag("listen", "address to listen on").String(&c.Listen).Default(c.Listen)
	serve.Run(c.Serve)
}

func (c *Command) Serve(ctx context.Context) error {
	ln, err := socket.ListenUp(c.Listen, 5)
	if err != nil {
		return err
	}
	defer ln.Close()
	log.Infof("Listening on %s", format(ln.Addr()))
	return c.server.Serve(ctx, ln)
}

func format(addr net.Addr) string {
	address := addr.String()
	if addr.Network() == "unix" {
		return address
	}
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		// Give up trying to format.
		// TODO: figure out if this can occur.
		return address
	}
	// https://serverfault.com/a/444557
	if host == "::" {
		host = "0.0.0.0"
	}
	return fmt.Sprintf("http://%s:%s", host, port)
}
