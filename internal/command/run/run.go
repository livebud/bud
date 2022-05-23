package run

import (
	"context"
	"fmt"
	"net"
	"time"

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

func displayDashboard(host, port string, timeElapsed int) {
	/*
		The dashboard should looks like this:
			| Listening on: http://127.0.0.1:3000
			| Ready in 270ms
	*/
	console.Info(fmt.Sprintf("Listening on: http://%s:%s", host, port))
	console.Info(fmt.Sprintf("Ready in %dms", timeElapsed))
}

func (c *Command) Run(ctx context.Context) error {
	// Start timer
	startTime := time.Now()
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
	// Load the compiler
	compiler, err := bud.Find(c.bud.Dir)
	if err != nil {
		return err
	}
	// Compile the project CLI
	project, err := compiler.Compile(ctx, &c.bud.Flag)
	if err != nil {
		return err
	}
	// Run the project
	process, err := project.Run(ctx, listener)
	if err != nil {
		return err
	}
	// Measure elapsed time
	timeElapsed := int(time.Since(startTime).Milliseconds())
	displayDashboard(host, port, timeElapsed)
	return process.Wait()
}
