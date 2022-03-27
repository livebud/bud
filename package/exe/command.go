package exe

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

func Command(ctx context.Context, name string, args ...string) (c *Cmd) {
	return (*Cmd)(exec.CommandContext(ctx, name, args...))
}

type Cmd exec.Cmd

func (c *Cmd) cmd() *exec.Cmd {
	return (*exec.Cmd)(c)
}

func (c *Cmd) Close() error {
	cmd := c.cmd()
	sp := cmd.Process
	if sp != nil {
		if err := sp.Signal(os.Interrupt); err != nil {
			sp.Kill()
		}
	}
	if err := cmd.Wait(); err != nil {
		if !isExitStatus(err) && !isWaitError(err) {
			return err
		}
	}
	return nil
}

func (c *Cmd) Wait() error {
	return c.cmd().Wait()
}

func isExitStatus(err error) bool {
	return err != nil && strings.Contains(err.Error(), "exit status ")
}

func isWaitError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Wait was already called")
}

func (c *Cmd) Run() error {
	return c.cmd().Run()
}

func (c *Cmd) Start() error {
	return c.cmd().Start()
}

func (c *Cmd) Restart(ctx context.Context) error {
	// Close the process first
	if err := c.Close(); err != nil {
		return err
	}
	cmd := c.cmd()
	// Re-run the command again. cmd.Args[0] is the path, so we skip that.
	next := Command(ctx, cmd.Path, cmd.Args[1:]...)
	next.Env = cmd.Env
	next.Stdout = cmd.Stdout
	next.Stderr = cmd.Stderr
	next.Stdin = cmd.Stdin
	next.ExtraFiles = cmd.ExtraFiles
	next.Dir = cmd.Dir
	if err := next.cmd().Start(); err != nil {
		return err
	}
	// Point to the new command
	*c = *next
	return nil
}
