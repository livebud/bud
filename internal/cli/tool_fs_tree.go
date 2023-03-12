package cli

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"

	"github.com/livebud/bud/framework"
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
	abs, err := filepath.Abs(c.Dir)
	if err != nil {
		return err
	}
	fsys := virtual.OS(abs)
	sub, err := fs.Sub(fsys, path.Clean(in.Path))
	if err != nil {
		return err
	}
	// Print out the tree
	tree, err := virtual.Print(sub)
	if err != nil {
		return err
	}
	fmt.Fprintln(c.Stdout, tree)

	return nil
}
