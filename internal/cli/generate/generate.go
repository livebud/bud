package generate

import (
	"context"
	"os/exec"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/cli/bud"
)

// New command for bud run.
func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud: bud,
		Flag: &framework.Flag{
			Env:    in.Env,
			Stderr: in.Stderr,
			Stdin:  in.Stdin,
			Stdout: in.Stdout,
		},
	}
}

// Command for bud run.
type Command struct {
	bud   *bud.Command
	Flag  *framework.Flag
	Paths []string
}

func (c *Command) Run(ctx context.Context) error {
	// Load the logger
	log, err := bud.Log(c.Flag.Stdout, c.bud.Log)
	if err != nil {
		return err
	}
	// Load the module
	module, err := bud.Module(c.bud.Dir)
	if err != nil {
		return err
	}
	// Setup the filesystem
	bfs, err := bud.FileSystem(ctx, log, module, c.Flag)
	if err != nil {
		return err
	}
	defer bfs.Close()
	// Generate the application
	if err := bfs.Sync(module, "bud/internal"); err != nil {
		return err
	}
	// Fast-track for empty `go generate`
	if len(c.Paths) == 0 {
		return nil
	}
	// Run `go generate [paths...]`
	args := append([]string{"generate"}, c.Paths...)
	cmd := exec.CommandContext(ctx, "go", args...)
	cmd.Dir = module.Directory()
	cmd.Stdout = c.Flag.Stdout
	cmd.Stderr = c.Flag.Stderr
	cmd.Stdin = c.Flag.Stdin
	cmd.Env = c.Flag.Env
	return cmd.Run()
}
