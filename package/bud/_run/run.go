package run

import (
	"context"
	"io"
	"io/ioutil"
	"net"
	"os"

	"github.com/livebud/bud/package/bud/expand"
	"github.com/livebud/bud/package/exe"

	"github.com/livebud/bud/package/bud"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/socket"
)

type Command struct {
	Dir      string
	Listener net.Listener
	Flag     *bud.Flag
	Env      bud.Env
	Stdout   io.Writer
	Stderr   io.Writer
}

func (c *Command) compile(ctx context.Context) (*exe.Cmd, error) {
	// Discard stderr and stdout
	if c.Stdout == nil {
		c.Stdout = ioutil.Discard
	}
	if c.Stderr == nil {
		c.Stderr = ioutil.Discard
	}

	// Default flags for running
	if c.Flag == nil {
		c.Flag = &bud.Flag{
			Hot:    ":35729",
			Embed:  false,
			Minify: false,
		}
	}

	// Find go.mod
	module, err := gomod.Find(c.Dir)
	if err != nil {
		return nil, err
	}

	// Default environment
	c.Env = c.Env.Defaults(bud.Env{
		"HOME":       os.Getenv("HOME"),
		"PATH":       os.Getenv("PATH"),
		"GOPATH":     os.Getenv("GOPATH"),
		"GOMODCACHE": module.ModCache(),
		"TMPDIR":     os.TempDir(),
	})

	// Expand macros
	macro := &expand.Command{
		Module: module,
		Flag:   c.Flag,
	}
	if err := macro.Expand(ctx); err != nil {
		return nil, err
	}

	// Run the project CLI
	// $ bud/cli run
	cmd := exe.Command(ctx, "bud/cli", "run")
	cmd.Dir = module.Directory()
	cmd.Env = c.Env.List()
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	// Inject the listener into the command to be available in the subprocess
	if err := socket.Inject(cmd, c.Listener); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (c *Command) Run(ctx context.Context) error {
	cmd, err := c.compile(ctx)
	if err != nil {
		return err
	}
	return cmd.Run()
}

func (c *Command) Start(ctx context.Context) (*exe.Cmd, error) {
	cmd, err := c.compile(ctx)
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, err
}
