package v8client

import (
	"context"
	"encoding/gob"
	"errors"
	"os"
	"os/exec"
)

// Launch the process and return a client
func Load(ctx context.Context) (c *Client, err error) {
	// Try getting the BUD_PATH that's been passed in
	budPath := os.Getenv("BUD_PATH")
	if budPath == "" {
		budPath, err = exec.LookPath("bud")
		if err != nil {
			return nil, err
		}
	}
	cmd := exec.CommandContext(ctx, budPath, "tool", "v8", "client")
	cmd.Env = os.Environ()
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
	return &Client{
		cmd:    cmd,
		stdin:  gob.NewEncoder(stdin),
		stdout: gob.NewDecoder(stdout),
	}, nil
}

type Client struct {
	cmd    *exec.Cmd
	stdin  *gob.Encoder
	stdout *gob.Decoder
}

func (c *Client) Script(path, script string) error {
	if err := c.stdin.Encode(input{Type: "script", Path: path, Code: script}); err != nil {
		return err
	}
	var out output
	if err := c.stdout.Decode(&out); err != nil {
		return err
	}
	if out.Error != "" {
		return errors.New(out.Error)
	}
	return nil
}

func (c *Client) Eval(path, expr string) (value string, err error) {
	if err := c.stdin.Encode(input{Type: "eval", Path: path, Code: expr}); err != nil {
		return "", err
	}
	var out output
	if err := c.stdout.Decode(&out); err != nil {
		return "", err
	}
	if out.Error != "" {
		return "", errors.New(out.Error)
	}
	return out.Result, nil
}

func (c *Client) Close() error {
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
