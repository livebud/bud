package exe

import (
	"context"
	"errors"
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
	cmd := (*exec.Cmd)(c)
	sp := cmd.Process
	if sp != nil {
		if err := sp.Signal(os.Interrupt); err != nil {
			if isProcessDone(err) {
				return nil
			}
			sp.Kill()
		}
	}
	if err := c.Wait(); err != nil {
		return err
	}
	return nil
}

func (c *Cmd) Wait() error {
	if err := (*exec.Cmd)(c).Wait(); err != nil {
		if canIgnore(err) {
			return nil
		}
		return err
	}
	return nil
}

// Errors we can safely ignore when closing the process
// TODO: if we find ourselves squelching real errors, we will want to revisit
// this as it might be overly aggressive.
func canIgnore(err error) bool {
	return isExitStatus(err) ||
		isInterrupt(err) ||
		isKilled(err) ||
		isWaitError(err)
}

func isExitStatus(err error) bool {
	return err != nil && strings.Contains(err.Error(), "exit status ")
}

func isProcessDone(err error) bool {
	return errors.Is(err, os.ErrProcessDone)
}

func isInterrupt(err error) bool {
	return err != nil && err.Error() == `signal: interrupt`
}

func isKilled(err error) bool {
	return err != nil && err.Error() == `signal: killed`
}

func isWaitError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Wait was already called")
}

func (c *Cmd) Run() error {
	if err := (*exec.Cmd)(c).Run(); err != nil {
		if canIgnore(err) {
			return nil
		}
		return err
	}
	return nil
}

func (c *Cmd) Start() error {
	return (*exec.Cmd)(c).Start()
}

func (c *Cmd) Restart(ctx context.Context) error {
	// Close the process first
	if err := c.Close(); err != nil {
		return err
	}
	cmd := (*exec.Cmd)(c)
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
