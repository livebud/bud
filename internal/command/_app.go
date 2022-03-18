package command

import (
	"context"
	"os/exec"

	"gitlab.com/mnm/bud/pkg/gomod"
)

type App struct {
	module *gomod.Module
}

func (a *App) Command(ctx context.Context, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, a.module.Directory("bud", "app"), args...)
}

func (a *App) Start(ctx context.Context) (*Process, error) {
	return nil, nil
}
