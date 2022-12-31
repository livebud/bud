package build

import (
	"context"

	"github.com/livebud/bud/internal/budfs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/cli/bud"
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
	// Setup the listener
	budln, err := bud.BudListener(c.in)
	if err != nil {
		return err
	}
	defer budln.Close()
	// Setup the command shell
	cmd := bud.Shell(c.in, module)
	cmd.Env = append(cmd.Env, "BUD_LISTEN="+budln.Addr().String())
	// Load the budfs
	bfs, err := budfs.Load(cmd, c.Flag, module, log)
	if err != nil {
		return err
	}
	defer bfs.Close(ctx)
	// Start the server
	budServer, err := bud.StartBudServer(ctx, budln, bfs, log)
	if err != nil {
		return err
	}
	defer budServer.Close()
	// Sync the generated files
	if err := bfs.Sync(ctx, module); err != nil {
		return err
	}
	return nil
}
