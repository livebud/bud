package build

import (
	"context"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/bfs"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/gobuild"
	"github.com/livebud/bud/internal/versions"
)

// New command for bud build
func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud: bud,
		in:  in,
		Flag: &framework.Flag{
			Env:    in.Env,
			Stderr: in.Stderr,
			Stdin:  in.Stdin,
			Stdout: in.Stdout,
		},
	}
}

// Command for running bud build
type Command struct {
	bud  *bud.Command
	in   *bud.Input
	Flag *framework.Flag
}

// Run the build command
func (c *Command) Run(ctx context.Context) error {
	// Find go.mod
	module, err := bud.Module(c.bud.Dir)
	if err != nil {
		return err
	}
	// Ensure we have version alignment between the CLI and the runtime
	if err := bud.EnsureVersionAlignment(ctx, module, versions.Bud); err != nil {
		return err
	}
	// Setup the logger
	log, err := bud.Log(c.in.Stderr, c.bud.Log)
	if err != nil {
		return err
	}
	bfs, err := bfs.Load(c.Flag, log, module)
	if err != nil {
		return err
	}
	defer bfs.Close()
	// Generate the application
	if err := bfs.Sync("bud/internal"); err != nil {
		return err
	}
	builder := gobuild.New(module)
	return builder.Build(ctx, "bud/internal/app/main.go", "bud/app")
}
