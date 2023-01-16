package cli

import (
	"context"
	"os"

	"github.com/livebud/bud/framework"
)

type ToolCacheClean struct {
	Flag      *framework.Flag
	ListenDev string
}

func (c *CLI) ToolCacheClean(ctx context.Context, in *ToolCacheClean) error {
	module, err := c.findModule()
	if err != nil {
		return err
	}
	// Remove the old cache
	// TODO: this can be removed in a future release
	if err := os.RemoveAll(module.Directory("bud", ".cache")); err != nil {
		return err
	}
	// Remove the new SQLite cache
	if err := os.RemoveAll(module.Directory("bud", "bud.db")); err != nil {
		return err
	}
	return nil
}
