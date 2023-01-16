package cli

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/printfs"
	"github.com/livebud/bud/package/virtual"
)

type ToolFsTree struct {
	Flag *framework.Flag
	Path string
}

func (c *CLI) ToolFsTree(ctx context.Context, in *ToolFsTree) error {
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

	// Print out the tree
	tree, err := printfs.Print(fsys, path.Clean(in.Path))
	if err != nil {
		return err
	}
	fmt.Fprintln(c.Stdout, tree)

	return nil
}
