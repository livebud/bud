package generate

import (
	"context"
	"strings"

	"github.com/livebud/bud/internal/config"
)

// New command for bud generate
func New(provide config.Provide) *Command {
	return &Command{provide: provide}
}

// Command for running bud generate
type Command struct {
	provide config.Provide
	Args    []string
}

// Run the generate command
func (c *Command) Run(ctx context.Context) error {
	module, err := c.provide.Module()
	if err != nil {
		return err
	}
	budsvr, err := c.provide.BudServer()
	if err != nil {
		return err
	}
	defer budsvr.Close()
	budfs, err := c.provide.BudFileSystem()
	if err != nil {
		return err
	}
	defer budfs.Close(ctx)
	// Sync the directories
	if err := budfs.Sync(ctx, module, selectBudDirs(c.Args)...); err != nil {
		return err
	}
	return nil
}

func selectBudDirs(patterns []string) (paths []string) {
	for _, pattern := range patterns {
		// Only sync from within the bud directory
		if !strings.HasPrefix(pattern, "bud/") {
			continue
		}
		// Trim the wildcard suffix since SyncDirs is recursive already
		// TODO: support non-recursive syncs
		pattern = strings.TrimSuffix(pattern, "/...")
		// Add the file or directory
		paths = append(paths, pattern)
	}
	return paths
}
