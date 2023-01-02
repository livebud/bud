package build

import (
	"context"

	"github.com/livebud/bud/internal/config"
)

// New command for bud build
func New(provide config.Provide) *Command {
	return &Command{provide}
}

// Command for running bud build
type Command struct {
	provide config.Provide
}

// Run the build command
func (c *Command) Run(ctx context.Context) error {
	module, err := c.provide.Module()
	if err != nil {
		return err
	}
	budServer, err := c.provide.BudServer()
	if err != nil {
		return err
	}
	defer budServer.Close()
	budfs, err := c.provide.BudFileSystem()
	if err != nil {
		return err
	}
	defer budfs.Close(ctx)
	// // Find go.mod
	// module, err := bud.Module(c.bud.Dir)
	// if err != nil {
	// 	return err
	// }
	// // Ensure we have version alignment between the CLI and the runtime
	// if err := bud.EnsureVersionAlignment(ctx, module, versions.Bud); err != nil {
	// 	return err
	// }
	// // Setup the logger
	// log, err := bud.Log(c.in.Stderr, c.bud.Log)
	// if err != nil {
	// 	return err
	// }
	// // Setup the listener
	// budln, err := bud.BudListener(c.in)
	// if err != nil {
	// 	return err
	// }
	// defer budln.Close()
	// // Setup the command shell
	// cmd := bud.Shell(c.in, module)
	// cmd.Env = append(cmd.Env, "BUD_LISTEN="+budln.Addr().String())
	// // Load the budfs
	// bfs, err := budfs.Load(budln, cmd, c.Flag, module, log)
	// if err != nil {
	// 	return err
	// }
	// defer bfs.Close(ctx)
	// // Start the server
	// budServer, err := bud.StartBudServer(ctx, budln, bfs, log)
	// if err != nil {
	// 	return err
	// }
	// defer budServer.Close()
	// Sync the generated files
	if err := budfs.Sync(ctx, module); err != nil {
		return err
	}
	return nil
}
