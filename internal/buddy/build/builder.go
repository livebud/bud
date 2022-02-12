package build

import (
	"context"

	"gitlab.com/mnm/bud/pkg/buddy"
)

func New(kit buddy.Kit) *Command {
	return &Command{kit}
}

type Command struct {
	kit buddy.Kit
}

type Option func(option *option)

type option struct {
}

func (c *Command) Build(ctx context.Context, options ...Option) error {
	panic("Build not implemented yet")
	// return nil
}
