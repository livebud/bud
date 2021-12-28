package v8client

import (
	"context"
	"encoding/gob"
	"os"
	"os/exec"
)

func launchBudToolV8(ctx context.Context, command string, args ...string) (*Command, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stderr = os.Stderr
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Command{
		cmd:    cmd,
		stdin:  gob.NewEncoder(stdin),
		stdout: gob.NewDecoder(stdout),
	}, nil
}

type Command struct {
	cmd    *exec.Cmd
	stdin  *gob.Encoder
	stdout *gob.Decoder
}

func (c *Command) Eval(path, expr string) (value string, err error) {
	if err := c.stdin.Encode(expr); err != nil {
		return "", err
	}
	var raw string
	if err := c.stdout.Decode(&raw); err != nil {
		return "", err
	}
	return string(raw), nil
}

func (c *Command) Close() error {
	if c.cmd.Process == nil {
		return nil
	}
	if err := c.cmd.Process.Signal(os.Interrupt); err != nil {
		return err
	}
	if err := c.cmd.Wait(); err != nil && err.Error() != "signal: interrupt" {
		return err
	}
	return nil
}
