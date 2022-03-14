package project

import (
	"context"

	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"
)

type Generator interface {
	Generate(ctx context.Context) error
}

type Command struct {
	FS  *overlay.FileSystem
	Dir string
}

func (c *Command) Compile(ctx context.Context) (*App, error) {
	module, err := gomod.Find(c.Dir)
	if err != nil {
		return nil, err
	}
	if err := c.FS.Sync("bud/.app"); err != nil {
		return nil, err
	}
	if err := gobin.Build(ctx, module.Directory(), "bud/.app/main.go", "bud/app"); err != nil {
		return nil, err
	}
	return &App{module}, nil
}
