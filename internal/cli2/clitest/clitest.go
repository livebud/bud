package clitest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"github.com/livebud/bud/internal/pubsub"

	"github.com/livebud/bud/package/socket"

	"golang.org/x/sync/errgroup"

	cli "github.com/livebud/bud/internal/cli2"
	"github.com/livebud/bud/internal/envs"
)

func New(dir string) *CLI {
	return &CLI{
		dir: dir,
		bus: pubsub.New(),
		Env: envs.Map{
			"NO_COLOR": "1",
			"HOME":     os.Getenv("HOME"),
			"PATH":     os.Getenv("PATH"),
			"GOPATH":   os.Getenv("GOPATH"),
			"TMPDIR":   os.TempDir(),
		},
		Stdin: nil,
	}
}

type CLI struct {
	dir   string
	bus   pubsub.Client
	Env   envs.Map
	Stdin io.Reader
}

func (c *CLI) toCLI(stdout, stderr io.Writer) *cli.CLI {
	return &cli.CLI{
		Dir:    c.dir,
		Bus:    c.bus,
		Env:    c.Env.List(),
		Stdin:  c.Stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
}

func (c *CLI) Run(ctx context.Context, args ...string) (*Result, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cli := c.toCLI(stdout, stderr)
	err := cli.Run(ctx, args...)
	return &Result{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}, err
}

type Result struct {
	Stdout string
	Stderr string
}

func (c *CLI) Start(ctx context.Context, args ...string) (*App, error) {
	stdout, stdoutWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	stderr, stderrWriter, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	// TODO: listen unix and create client
	web, err := socket.Listen(":0")
	if err != nil {
		return nil, err
	}
	// TODO: listen unix and create client
	hot, err := socket.Listen(":0")
	if err != nil {
		return nil, err
	}
	cli := c.toCLI(stdoutWriter, stderrWriter)
	cli.Web = web
	cli.Hot = hot
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return cli.Run(ctx, args...)
	})
	return &App{
		eg:     eg,
		bus:    c.bus,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

type App struct {
	eg        *errgroup.Group
	bus       pubsub.Client
	stdout    io.Reader
	stderr    io.Reader
	webClient *http.Client
	hotClient *http.Client
}

// Subscribe to an event
func (a *App) Subscribe(topics ...string) pubsub.Subscription {
	return a.bus.Subscribe(topics...)
}

// Publish an event
func (a *App) Publish(topic string, payload []byte) {
	a.bus.Publish(topic, payload)
}
