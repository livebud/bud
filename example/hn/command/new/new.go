package new

import (
	"context"
	"fmt"

	"gitlab.com/mnm/bud/package/js"
)

type Command struct {
	VM     js.VM
	DryRun bool `flag:"dry-run" help:"run but don't write" default:"false"`
}

func (c *Command) Run(ctx context.Context) error {
	fmt.Println("creating new", c.DryRun)
	return nil
}
