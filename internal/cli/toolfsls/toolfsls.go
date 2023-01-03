package toolfsls

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"sort"

	"github.com/livebud/bud/internal/config"
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
	des, err := fs.ReadDir(budfs, path.Clean(c.Dir))
	if err != nil {
		return err
	}
	// Directories come first
	sort.Slice(des, func(i, j int) bool {
		if des[i].IsDir() && !des[j].IsDir() {
			return true
		} else if !des[i].IsDir() && des[j].IsDir() {
			return false
		}
		return des[i].Name() < des[j].Name()
	})
	// Print out list
	for _, de := range des {
		name := de.Name()
		if de.IsDir() {
			name += "/"
		}
		fmt.Fprintln(os.Stdout, name)
	}
	return nil
}
