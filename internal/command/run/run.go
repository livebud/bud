package run

import (
	"context"
	"net"
	"time"

	"github.com/livebud/bud/internal/bud"
	"github.com/livebud/bud/internal/command"
	"github.com/livebud/bud/internal/version"
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

func frontEnd() string {
	// TODO: Change to the corresponding front end in next releases
	return "Svelte " + version.Svelte
}

func displayDashboard(host, port string, timeElapsed time.Duration) {
	/*
		The dashboard should looks something like this:
			|   bud dev server is running:
			|
			| > Listening on: http://127.0.0.1:3000
			| > Front end: Svelte 3.47.0
			|
			|   Ready in 270.131758ms
			|
	*/
	address := "http://" + host + ":" + port
	console.Info("  bud dev server is running:")
	console.Info("")
	console.Info("> Listening on: " + address)
	console.Info("> Front end: " + frontEnd())
	console.Info("")
	console.Info("  Ready in " + timeElapsed.String())
	console.Info("")
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
	// Measure elapsed time
	timeElapsed := time.Since(startTime)
	displayDashboard(host, port, timeElapsed)
	return process.Wait()
}
