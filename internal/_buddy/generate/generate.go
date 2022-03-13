package generate

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

func (c *Command) Generate(ctx context.Context, options ...Option) error {
	panic("Generate not implemented yet")
	// return nil
}
