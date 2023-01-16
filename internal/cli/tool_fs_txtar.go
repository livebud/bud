package cli

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/virtual"
	"golang.org/x/tools/txtar"
)

type ToolFsTxtar struct {
	Flag *framework.Flag
	Path string
}

func (c *CLI) ToolFsTxtar(ctx context.Context, in *ToolFsTxtar) error {
	// Generate bud files
	generate := &Generate{Flag: in.Flag}
	if err := c.Generate(ctx, generate); err != nil {
		return err
	}

	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	fsys := virtual.OS(wd)

	// Walk the directory, adding all the files to the archive
	ar := new(txtar.Archive)
	dir := path.Clean(in.Path)
	err = fs.WalkDir(fsys, dir, func(path string, de fs.DirEntry, err error) error {
		if err != nil {
			return err
		} else if de.IsDir() {
			return nil
		}
		code, err := fs.ReadFile(fsys, path)
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
	fmt.Fprintln(c.Stdout, string(txtar.Format(ar)))
	return nil
}
