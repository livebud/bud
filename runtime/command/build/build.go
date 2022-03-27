package build

import (
	"context"

	"gitlab.com/mnm/bud/runtime/bud"
)

type Command struct {
	Project *bud.Project
	Flag    bud.Flag
}

func (c *Command) Run(ctx context.Context) error {
	_, err := c.Project.Compile(ctx, &c.Flag)
	if err != nil {
		return err
	}
	return nil
}
