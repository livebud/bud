package build

import (
	"context"
	"fmt"

	"gitlab.com/mnm/bud/runtime/project"
)

type Command struct {
	Project *project.Command
	Embed   bool
	Hot     bool
	Minify  bool
}

func (c *Command) Run(ctx context.Context) error {
	fmt.Println("building app")
	return nil
}
