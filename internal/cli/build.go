package cli

import (
	"context"

	"github.com/livebud/bud/framework"
)

type Build struct {
	Flag *framework.Flag
}

func (c *CLI) Build(ctx context.Context, in *Build) error {
	return c.Generate(ctx, &Generate{Flag: in.Flag})
}
