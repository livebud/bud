package cli

import (
	"context"

	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/package/commander"
)

type Custom struct {
	Closer *once.Closer
	Help   bool
	Args   []string
}

func (c *CLI) Custom(ctx context.Context, in *Custom) error {
	if in.Help {
		return commander.Usage()
	}
	return commander.Usage()
}
