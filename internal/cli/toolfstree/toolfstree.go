package toolfstree

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/livebud/bud/internal/config"
	"github.com/livebud/bud/internal/printfs"
)

func New(provide config.Provide) *Command {
	return &Command{provide: provide}
}

type Command struct {
	provide config.Provide
	Dir     string
}

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
	if err := budfs.Sync(ctx, module); err != nil {
		return err
	}
	tree, err := printfs.Print(budfs, path.Clean(c.Dir))
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, tree)
	return nil
}
