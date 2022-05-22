package build

import (
	"context"
	"io"
	"io/ioutil"
	"os"

	"github.com/livebud/bud/package/bud"
	"github.com/livebud/bud/package/bud/expand"
	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/gomod"
)

// Command for `bud build`
type Command struct {
	Dir    string
	Flag   *bud.Flag
	Env    bud.Env
	Stdout io.Writer
	Stderr io.Writer
}

func (c *Command) Build(ctx context.Context) error {
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
			Hot:    "",   // Live Reload is disabled by default
			Embed:  true, // Embed assets by default
			Minify: true, // Minify assets by default
		}
	}

	// Find go.mod
	module, err := gomod.Find(c.Dir)
	if err != nil {
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

	// Expand macros
	macro := &expand.Command{
		Module: module,
		Flag:   c.Flag,
	}
	if err := macro.Expand(ctx); err != nil {
		return err
	}

	// Run the project CLI
	// $ bud/cli build
	cmd := exe.Command(ctx, "bud/cli", "build")
	cmd.Dir = module.Directory()
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
