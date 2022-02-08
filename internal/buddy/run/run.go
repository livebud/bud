package run

import (
	"context"

	"gitlab.com/mnm/bud/pkg/gomod"
)

func New(mod *gomod.Module) *Command {
	return &Command{mod}
}

type Command struct {
	mod *gomod.Module
}

type Option func(option *option)

type option struct {
}

func (c *Command) Run(ctx context.Context, options ...Option) error {
	panic("Run not implemented yet")
	// return nil
}
