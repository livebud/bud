package bud

import (
	"context"
	"io"
	"net"
	"os/exec"

	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/socket"
)

type App struct {
	Module *gomod.Module
	Env    []string
	Stdout io.Writer
	Stderr io.Writer
}

func (a *App) args(args ...string) []string {
	return args
}

func (a *App) command(ctx context.Context, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, a.Module.Directory("bud", "app"), args...)
	cmd.Dir = a.Module.Directory()
	cmd.Env = a.Env
	cmd.Stderr = a.Stderr
	cmd.Stdout = a.Stdout
	return cmd
}

func (a *App) Executor(ctx context.Context, args ...string) *exec.Cmd {
	return a.command(ctx, a.args(args...)...)
}

// Execute a custom command
func (a *App) Execute(ctx context.Context, args ...string) error {
	cmd := a.Executor(ctx, args...)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) Run(ctx context.Context) error {
	cmd := a.command(ctx)
	return cmd.Run()
}

func (a *App) Start(ctx context.Context, listener net.Listener) (*Process, error) {
	// Pass the socket through
	files, env, err := socket.Files(listener)
	if err != nil {
		return nil, err
	}
	cmd := a.command(ctx)
	cmd.Env = append(a.Env, string(env))
	cmd.ExtraFiles = append(cmd.ExtraFiles, files...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Process{cmd}, nil
}
