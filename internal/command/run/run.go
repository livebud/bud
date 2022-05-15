package run

import (
	"context"
	"net"

	"github.com/livebud/bud/internal/bud"
	"github.com/livebud/bud/internal/command"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/socket"
)

func New(bud *command.Bud) *Command {
	return &Command{bud: bud}
}

type Command struct {
	bud  *command.Bud
	Port string
}

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
	// Load the compiler
	compiler, err := bud.Find(c.bud.Dir)
	if err != nil {
		return err
	}
	// Compiler the project CLI
	project, err := compiler.Compile(ctx, &c.bud.Flag)
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
