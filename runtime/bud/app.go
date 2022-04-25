package bud

import (
	"context"
	"io"
	"net"

	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/socket"
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

func (a *App) command(ctx context.Context, args ...string) *exe.Cmd {
	cmd := exe.Command(ctx, a.Module.Directory("bud", "app"), args...)
	cmd.Dir = a.Module.Directory()
	cmd.Env = a.Env
	cmd.Stderr = a.Stderr
	cmd.Stdout = a.Stdout
	return cmd
}

func (a *App) Executor(ctx context.Context, args ...string) *exe.Cmd {
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

// Start the application
func (a *App) Start(ctx context.Context, listener net.Listener) (*exe.Cmd, error) {
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
	return cmd, nil
}

// Run the app and wait for the result
func (a *App) Run(ctx context.Context, listener net.Listener) error {
	cmd, err := a.Start(ctx, listener)
	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
