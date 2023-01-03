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
	// Sync the generated files
	if err := budfs.Sync(ctx, module); err != nil {
		return err
	}
	return nil
}
