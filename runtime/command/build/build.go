package build

import (
	"context"

	"gitlab.com/mnm/bud/runtime/project"
)

type Command struct {
	Project *project.Compiler
	Flag    project.Flag
}

func (c *Command) Run(ctx context.Context) error {
	_, err := c.Project.Compile(ctx, &c.Flag)
	if err != nil {
		return err
	}
	return nil
}
