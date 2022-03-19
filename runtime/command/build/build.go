package build

import (
	"context"

	"gitlab.com/mnm/bud/runtime/project"
)

type Command struct {
	Project *project.Compiler
	Embed   bool
	Hot     bool
	Minify  bool
}

func (c *Command) Run(ctx context.Context) error {
	_, err := c.Project.Compile(ctx)
	if err != nil {
		return err
	}
	return nil
}
