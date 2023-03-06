package cli

import (
	"context"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"sort"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/virtual"
)

type ToolFsLs struct {
	Flag *framework.Flag
	Path string
}

func (c *CLI) ToolFsLs(ctx context.Context, in *ToolFsLs) error {
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
	// Read the directory out
	des, err := fs.ReadDir(fsys, path.Clean(in.Path))
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
		fmt.Fprintln(c.Stdout, name)
	}
	return nil
}
