package bud

import (
	"context"
	"net"
	"os"
	"os/exec"

	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/socket"
)

type CLI struct {
	flag   Flag
	module *gomod.Module
}

func (c *CLI) Command(ctx context.Context, args ...string) *exec.Cmd {
	return exec.CommandContext(ctx, c.module.Directory("bud", "cli"), args...)
}

// Custom executes a custom command
func (c *CLI) Custom(ctx context.Context, args ...string) error {
	cmd := c.Command(ctx, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (c *CLI) Build(ctx context.Context) error {
	args := append([]string{"build"}, c.flag.Args()...)
	cmd := c.Command(ctx, args...)
	cmd.Dir = c.module.Directory()
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (c *CLI) Run(ctx context.Context, listener net.Listener) error {
	// Pass the socket through
	files, env, err := socket.Files(listener)
	if err != nil {
		return err
	}
	// Create the command
	args := append([]string{"run"}, c.flag.Args()...)
	cmd := c.Command(ctx, args...)
	cmd.Dir = c.module.Directory()
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
