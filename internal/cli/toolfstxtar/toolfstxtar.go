package toolfstxtar

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/livebud/bud/internal/config"

	"golang.org/x/tools/txtar"
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
	ar := new(txtar.Archive)
	dir := path.Clean(c.Dir)
	err = fs.WalkDir(budfs, dir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if de.IsDir() {
			return nil
		}
		code, err := fs.ReadFile(budfs, path)
		if err != nil {
			return err
		}
		ar.Files = append(ar.Files, txtar.File{
			Name: path,
			Data: code,
		})
		return nil
	})
	if err != nil {
		return err
	}
	// Print the archive to stdout
	fmt.Fprintln(os.Stdout, string(txtar.Format(ar)))
	return nil
}
