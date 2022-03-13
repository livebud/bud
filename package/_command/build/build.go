// Package build is typically run inside the generated CLI, but could also be
// run programmically.
package build

import (
	"context"

	"gitlab.com/mnm/bud/internal/gobin"
)

type Generator interface {
	Generate(ctx context.Context) error
}

func Load(generator Generator) *Command {
	return &Command{Generator: generator}
}

type Command struct {
	// TODO switch to private once we change the command API
	// that will separate dependencies from flags/args
	Generator Generator
	Dir       string // Project directory
}

func (c *Command) Run(ctx context.Context) error {
	if err := c.Generator.Generate(ctx); err != nil {
		return err
	}
	if err := gobin.Build(ctx, c.Dir, "bud/.app/main.go", "bud/app"); err != nil {
		return err
	}
	return nil
}
