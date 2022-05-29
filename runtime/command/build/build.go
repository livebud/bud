package build

import (
	"context"
	"io/fs"

	"github.com/livebud/bud/internal/buildcache"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/runtime/command"
)

func New(module *gomod.Module) *Command {
	return &Command{
		module: module,
		// Default flags
		Flag: &command.Flag{
			Embed:  true,
			Hot:    "",
			Minify: true,
		},
	}
}

// Command to build the project at runtime
type Command struct {
	module *gomod.Module

	// Below are filled in by the CLI
	Flag *command.Flag
	FS   *overlay.FileSystem
}

// Run is triggered by `bud/cli build`
func (c *Command) Run(ctx context.Context) error {
	// Sync the application
	if err := c.FS.Sync("bud/.app"); err != nil {
		return err
	}
	// Ensure that main.go exists
	if _, err := fs.Stat(c.module, "bud/.app/main.go"); err != nil {
		return err
	}
	// Build the application binary
	bcache := buildcache.Default(c.module)
	if err := bcache.Build(ctx, c.module, "bud/.app/main.go", "bud/app"); err != nil {
		return err
	}
	return nil
}
