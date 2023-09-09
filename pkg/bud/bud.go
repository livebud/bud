package bud

import (
	"context"
	"fmt"

	"github.com/livebud/bud/pkg/cli"
)

func New() *Command {
	return &Command{}
}

func Register(cli *cli.CLI, cmd *Command) {
	fmt.Println("registering!")
}

type Command struct {
}

type Run struct {
}

func (c *Command) Run(ctx context.Context, in *Run) error {
	return nil
}
