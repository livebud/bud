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

type ToolFsCat struct {
	Flag *framework.Flag
	Path string
}

func (c *CLI) ToolFsCat(ctx context.Context, in *ToolFsCat) error {
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

	// Read the file out
	code, err := fs.ReadFile(fsys, path.Clean(in.Path))
	if err != nil {
		return err
	}

	// Print it out
	fmt.Fprintln(c.Stdout, string(code))
	return nil
}
