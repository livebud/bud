// Package v8client is a client for submitting jobs to the v8server. The purpose
// of this package and the v8server package is because embedding v8 into your Go
// binary takes too much time, so we instead communicate back with the embedded
// V8 in the bud binary during development. When building, we embed V8 directly
// because builds can afford to be slower.
package v8client

import (
	"context"
	"encoding/gob"
	"errors"
	"io"
	"os"
	"os/exec"
	"sync"
)

// Launch the process and return a client
func Launch(ctx context.Context) (c *Client, err error) {
	// Get the BUD_PATH that's been passed in or fail. This should always be set
	// by the compiler
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
	// Close function to shut down the process gracefully
	closer := func() error {
		if cmd.Process == nil {
			return nil
		}
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			return err
		}
		if err := cmd.Wait(); err != nil && err.Error() != "signal: interrupt" {
			return err
		}
		return nil
	}
	return &Client{
		reader: gob.NewDecoder(stdout),
		writer: gob.NewEncoder(stdin),
		closer: closer,
	}, nil
}

// New client for testing
func New(reader io.Reader, writer io.Writer) *Client {
	return &Client{
		reader: gob.NewDecoder(reader),
		writer: gob.NewEncoder(writer),
		closer: func() error { return nil },
	}
}

// Client for evaluating scripts against the V8 server. This client may be used
// concurrently, but you cannot have multiple instances of clients communicating
// with the same server
type Client struct {
	// Synchronize readers, writers and closers
	mu     sync.Mutex
	closer func() error
	reader *gob.Decoder
	writer *gob.Encoder
}

func (c *Client) Script(path, script string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.writer.Encode(Input{Type: "script", Path: path, Code: script}); err != nil {
		return err
	}
	var out Output
	if err := c.reader.Decode(&out); err != nil {
		return err
	}
	if out.Error != "" {
		return errors.New(out.Error)
	}
	return nil
}

func (c *Client) Eval(path, expr string) (value string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.writer.Encode(Input{Type: "eval", Path: path, Code: expr}); err != nil {
		return "", err
	}
	var out Output
	if err := c.reader.Decode(&out); err != nil {
		return "", err
	}
	if out.Error != "" {
		return "", errors.New(out.Error)
	}
	return out.Result, nil
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closer()
}
