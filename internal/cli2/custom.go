package cli

import (
	"context"
	"fmt"

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
	fmt.Println("running custom with args", in.Args)
	return nil
}
