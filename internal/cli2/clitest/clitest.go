package clitest

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"os"

	"github.com/livebud/bud/package/socket"

	"golang.org/x/sync/errgroup"

	cli "github.com/livebud/bud/internal/cli2"
	"github.com/livebud/bud/internal/envs"
)

func Run(ctx context.Context, dir string, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	return New(dir).Run(ctx, args...)
}

func Start(ctx context.Context, dir string, args ...string) (*App, error) {
	return New(dir).Start(ctx, args...)
}

func New(dir string) *CLI {
	return &CLI{
		dir: dir,
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
	Env   envs.Map
	Stdin io.Reader
}

func (c *CLI) toCli(stdout, stderr io.Writer) *cli.CLI {
	return &cli.CLI{
		Dir:    c.dir,
		Env:    c.Env.List(),
		Stdin:  c.Stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
}

func (c *CLI) Run(ctx context.Context, args ...string) (*bytes.Buffer, *bytes.Buffer, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cli := c.toCli(stdout, stderr)
	return stdout, stderr, cli.Run(ctx, args...)
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
	cli := c.toCli(stdoutWriter, stderrWriter)
	cli.Web = web
	cli.Hot = hot
	eg := new(errgroup.Group)
	eg.Go(func() error {
		return cli.Run(ctx, args...)
	})
	return &App{
		eg:     eg,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

type App struct {
	eg *errgroup.Group

	stdout    io.Reader
	stderr    io.Reader
	webClient *http.Client
	hotClient *http.Client
}
