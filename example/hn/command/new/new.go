package new

import (
	"context"
	"fmt"

	v8 "gitlab.com/mnm/bud/js/v8"
)

type Command struct {
	V8     *v8.Pool
	DryRun bool `flag:"dry-run" help:"run but don't write" default:"false"`
}

func (c *Command) View(ctx context.Context) error {
	fmt.Println("creating new", c.DryRun)
	return nil
}
