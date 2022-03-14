package project

import (
	"context"
	"net"
	"os"
	"os/exec"

	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/socket"
)

type App struct {
	module *gomod.Module
}

func (a *App) Command(ctx context.Context, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, a.module.Directory("bud", "app"), args...)
}

func (a *App) Start(ctx context.Context, ln net.Listener) (*Process, error) {
	files, env, err := socket.Files(ln)
	if err != nil {
		return nil, err
	}
	cmd := a.Command(ctx)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Dir = a.module.Directory()
	cmd.ExtraFiles = append(cmd.ExtraFiles, files...)
	cmd.Env = append(os.Environ(), string(env))
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Process{cmd}, nil
}
