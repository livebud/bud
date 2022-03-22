package run

import (
	"context"
	"net"

	"gitlab.com/mnm/bud/internal/bud"
	"gitlab.com/mnm/bud/internal/command"
	"gitlab.com/mnm/bud/pkg/log/console"
	"gitlab.com/mnm/bud/pkg/socket"
)

type Command struct {
	Bud  *command.Bud
	Port string
}

func (c *Command) Run(ctx context.Context) error {
	ctx, shutdown, err := c.Bud.Tracer(ctx)
	if err != nil {
		return err
	}
	defer shutdown(&err)
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
	// Load the compiler
	compiler, err := bud.Find(c.Bud.Dir)
	if err != nil {
		return err
	}
	// Compiler the project CLI
	project, err := compiler.Compile(ctx, c.Bud.Flag)
	if err != nil {
		return err
	}
	// Run the project
	process, err := project.Run(ctx, listener)
	if err != nil {
		return err
	}
	return process.Wait()
}
