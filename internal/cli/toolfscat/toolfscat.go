package toolfscat

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/livebud/bud/internal/config"
)

func New(provide config.Provide) *Command {
	return &Command{provide: provide}
}

type Command struct {
	provide config.Provide
	Path    string
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
	code, err := fs.ReadFile(budfs, path.Clean(c.Path))
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(code))
	return nil
}
