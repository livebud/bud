package toolcache

import (
	"context"
	"os"

	"github.com/livebud/bud/internal/config"
)

func New(provide config.Provide) *Command {
	return &Command{provide: provide}
}

type Command struct {
	provide config.Provide
}

func (c *Command) Run(ctx context.Context) error {
	module, err := c.provide.Module()
	if err != nil {
		return err
	}
	// Remove the old cache
	if err := os.RemoveAll(module.Directory("bud", ".cache")); err != nil {
		return err
	}
	// Remove the new SQLite cache
	if err := os.RemoveAll(module.Directory("bud", "bud.db")); err != nil {
		return err
	}
	return nil
}
