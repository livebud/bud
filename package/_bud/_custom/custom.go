package custom

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"github.com/livebud/bud/package/bud/expand"
	"github.com/livebud/bud/package/exe"

	"github.com/livebud/bud/package/bud"
	"github.com/livebud/bud/package/gomod"
)

type Command struct {
	Dir    string
	Args   []string
	Env    bud.Env
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Command) Run(ctx context.Context) error {
	// Default inputs
	if c.Stdout == nil {
		c.Stdout = ioutil.Discard
	}
	if c.Stderr == nil {
		c.Stderr = ioutil.Discard
	}

	// Find the go.mod
	module, err := gomod.Find(c.Dir)
	if err != nil {
		return err
	}

	// Expand macros
	macro := &expand.Command{
		Module: module,
	}
	if err := macro.Expand(ctx); err != nil {
		return err
	}

	// Default environment
	c.Env = c.Env.Defaults(bud.Env{
		"HOME":       os.Getenv("HOME"),
		"PATH":       os.Getenv("PATH"),
		"GOPATH":     os.Getenv("GOPATH"),
		"GOMODCACHE": module.ModCache(),
		"TMPDIR":     os.TempDir(),
	})

	// Run the project CLI
	// $ bud/cli [args...]
	cmd := exe.Command(ctx, "bud/cli", c.Args...)
	cmd.Dir = module.Directory()
	cmd.Env = c.Env.List()
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	return cmd.Run()
}
