package run

import (
	"context"
	"net"

	"gitlab.com/mnm/bud/internal/bud"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/log/console"
	"gitlab.com/mnm/bud/pkg/socket"
)

type Command struct {
	Bud  *bud.Command
	Port string
}

// func (c *Command) Start(ctx context.Context) (*Process, error) {
// 	return nil, nil
// }

func (c *Command) Run(ctx context.Context) error {
	// Start listening on the port
	listener, err := socket.Load(c.Port)
	if err != nil {
		return err
	}
	defer listener.Close()
	host, port, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		return err
	}
	// https://serverfault.com/a/444557
	if host == "::" {
		host = "0.0.0.0"
	}
	console.Info("Listening on http://" + host + ":" + port)
	// Find go.mod
	module, err := gomod.Find(c.Bud.Dir)
	if err != nil {
		return err
	}
	// Compile the project CLI
	cli, err := c.Bud.Compile(ctx, module)
	if err != nil {
		return err
	}
	if err := cli.Run(ctx, listener); err != nil {
		return err
	}
	return nil
}
