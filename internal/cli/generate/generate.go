package generate

import (
	"context"
	"os/exec"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/bfs"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/versions"
)

// New command for bud generate
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

// Command for running bud generate
type Command struct {
	bud  *bud.Command
	in   *bud.Input
	Flag *framework.Flag
	Args []string
}

// Run the generate command
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
	// Load the filesystem
	bfs, err := bfs.Load(c.Flag, log, module)
	if err != nil {
		return err
	}
	defer bfs.Close()
	// Generate the application
	if err := bfs.Sync(); err != nil {
		return err
	}
	// Run go generate if we have any args
	if len(c.Args) > 0 {
		cmd := exec.CommandContext(ctx, "go", append([]string{"generate"}, c.Args...)...)
		cmd.Dir = module.Directory()
		cmd.Env = c.in.Env
		cmd.Stdin = c.in.Stdin
		cmd.Stdout = c.in.Stdout
		cmd.Stderr = c.in.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}
	return nil
}
