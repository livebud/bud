package run

import (
	"context"
	"os"
	"os/exec"

	"gitlab.com/mnm/bud/internal/gobin"
	"gitlab.com/mnm/bud/package/trace"
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

func (c *Command) Run(ctx context.Context) (err error) {
	ctx, span := trace.Start(ctx, "run app", "dir", c.Dir)
	defer span.End(&err)
	if err := c.generate(ctx); err != nil {
		return err
	}
	if err := c.build(ctx); err != nil {
		return err
	}
	if err := c.start(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Command) generate(ctx context.Context) (err error) {
	ctx, span := trace.Start(ctx, "generate app")
	defer span.End(&err)
	if err := c.Generator.Generate(ctx); err != nil {
		return err
	}
	return nil
}

func (c *Command) build(ctx context.Context) (err error) {
	ctx, span := trace.Start(ctx, "build app")
	defer span.End(&err)
	if err := gobin.Build(ctx, c.Dir, "bud/.app/main.go", "bud/app"); err != nil {
		return err
	}
	return nil
}

func (c *Command) start(ctx context.Context) (err error) {
	ctx, span := trace.Start(ctx, "start app")
	defer span.End(&err)
	cmd := exec.CommandContext(ctx, "bud/app")
	cmd.Dir = c.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
