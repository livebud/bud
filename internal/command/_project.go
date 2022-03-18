package command

import (
	"context"
	"net"
	"os"
	"os/exec"

	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/socket"
)

type project struct {
	module *gomod.Module
	Env    map[string]string
}

func (p *project) Command(ctx context.Context, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, p.module.Directory("bud", "cli"), args...)
}

// Execute a custom command
func (p *project) Execute(ctx context.Context, args ...string) error {
	cmd := p.Command(ctx, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (p *project) Build(ctx context.Context) error {
	args := append([]string{"build"})
	cmd := p.Command(ctx, args...)
	cmd.Dir = p.module.Directory()
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (p *project) Run(ctx context.Context, listener net.Listener) error {
	// Pass the socket through
	files, env, err := socket.Files(listener)
	if err != nil {
		return err
	}
	// Create the command
	args := append([]string{"run"}, p.config.Flags()...)
	cmd := p.Command(ctx, args...)
	cmd.Dir = p.module.Directory()
	cmd.Env = append(os.Environ(), string(env))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.ExtraFiles = append(cmd.ExtraFiles, files...)
	// Run the command
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
