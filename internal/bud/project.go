package bud

import (
	"context"
	"io"
	"net"

	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/runtime/bud"
)

type Project struct {
	Module *gomod.Module
	Env    Env
	Stdout io.Writer
	Stderr io.Writer
}

func (p *Project) command(ctx context.Context, args ...string) *exe.Cmd {
	cmd := exe.Command(ctx, p.Module.Directory("bud", "cli"), args...)
	cmd.Dir = p.Module.Directory()
	cmd.Env = p.Env.List()
	cmd.Stderr = p.Stderr
	cmd.Stdout = p.Stdout
	return cmd
}

func (p *Project) Executor(ctx context.Context, args ...string) *exe.Cmd {
	return p.command(ctx, args...)
}

// Execute a custom command
func (p *Project) Execute(ctx context.Context, args ...string) error {
	cmd := p.Executor(ctx, args...)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (p *Project) Builder(ctx context.Context) *exe.Cmd {
	return p.command(ctx, "build")
}

func (p *Project) Build(ctx context.Context) (*bud.App, error) {
	cmd := p.Builder(ctx)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return &bud.App{
		Module: p.Module,
		Env:    p.Env.List(),
		Stderr: p.Stderr,
		Stdout: p.Stdout,
	}, nil
}

func (p *Project) Runner(ctx context.Context, listener net.Listener) (*exe.Cmd, error) {
	// Pass the socket through
	files, env, err := socket.Files(listener)
	if err != nil {
		return nil, err
	}
	cmd := p.command(ctx, "run")
	cmd.Env = append(p.Env.List(), string(env))
	cmd.ExtraFiles = append(cmd.ExtraFiles, files...)
	return cmd, nil
}

func (p *Project) Run(ctx context.Context, listener net.Listener) (*exe.Cmd, error) {
	cmd, err := p.Runner(ctx, listener)
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}
