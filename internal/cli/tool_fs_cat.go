package cli

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/livebud/bud/framework"
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

	// Read the file out
	code, err := os.ReadFile(path.Clean(in.Path))
	if err != nil {
		return err
	}

	// Print it out
	fmt.Fprintln(c.Stdout, string(code))
	return nil
}
