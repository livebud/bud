package cli

import (
	"context"

	"github.com/livebud/bud/package/commander"
)

type custom struct {
	Help bool
	Args []string
}

func (c *CLI) runCustom(ctx context.Context, in *custom) error {
	if in.Help || len(in.Args) == 0 {
		return commander.Usage()
	}
	// TODO: generate the app and run a custom command
	return commander.Usage()
}
