package exe

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/livebud/bud/internal/once"
)

type Command struct {
	Dir        string
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	Env        []string
	ExtraFiles []*os.File
}

func (c *Command) cmd(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = c.Stdin
	cmd.Stdout = c.Stdout
	cmd.Stderr = c.Stderr
	cmd.Env = c.Env
	cmd.Dir = c.Dir
	cmd.ExtraFiles = c.ExtraFiles
	return cmd
}

func (c *Command) Run(ctx context.Context, name string, args ...string) error {
	cmd := c.cmd(name, args...)
	return Run(ctx, cmd)
}

func (c *Command) Start(ctx context.Context, name string, args ...string) (*Process, error) {
	cmd := c.cmd(name, args...)
	return Start(ctx, cmd)
}

func Run(ctx context.Context, cmd *exec.Cmd) error {
	p, err := Start(ctx, cmd)
	if err != nil {
		return err
	}
	return p.Wait()
}

func Start(ctx context.Context, cmd *exec.Cmd) (*Process, error) {
	p := wrap(ctx, cmd)
	if err := p.Start(); err != nil {
		return nil, err
	}
	return p, nil
}

// Wrap a command in a process
func wrap(ctx context.Context, cmd *exec.Cmd) *Process {
	return &Process{
		ctx:    ctx,
		cmd:    cmd,
		exitCh: make(chan error),
	}
}

type Process struct {
	cmd       *exec.Cmd
	ctx       context.Context
	exitCh    chan error
	closeOnce once.Error
}

func (p *Process) Start() error {
	if err := p.cmd.Start(); err != nil {
		return err
	}
	go p.wait()
	return nil
}

func (p *Process) Run() error {
	if err := p.Start(); err != nil {
		return err
	}
	return p.Wait()
}

func (p *Process) wait() {
	p.exitCh <- p.cmd.Wait()
}

// Close command
func (p *Process) close() error {
	sp := p.cmd.Process
	if sp == nil {
		return nil
	}
	expectError := isInterrupt
	if err := sp.Signal(os.Interrupt); err != nil {
		if isProcessDone(err) {
			return nil
		}
		expectError = isKilled
		if err := sp.Kill(); err != nil {
			return err
		}
	}
	err := <-p.exitCh
	close(p.exitCh)
	if err != nil {
		if !expectError(err) {
			return err
		}
	}
	return nil
}

func (p *Process) Close() error {
	return p.closeOnce.Do(p.close)
}

func (p *Process) Wait() error {
	select {
	case <-p.ctx.Done():
		if err := p.Close(); err != nil {
			return err
		}
		return nil
	case err := <-p.exitCh:
		return err
	}
}

func (p *Process) Restart(ctx context.Context) (*Process, error) {
	// Close the process first
	if err := p.Close(); err != nil {
		return nil, err
	}
	cmd := p.cmd
	// Re-run the command again. cmd.Args[0] is the path, so we skip that.
	next := exec.Command(cmd.Path, cmd.Args[1:]...)
	next.Env = cmd.Env
	next.Stdout = cmd.Stdout
	next.Stderr = cmd.Stderr
	next.Stdin = cmd.Stdin
	next.ExtraFiles = cmd.ExtraFiles
	next.Dir = cmd.Dir
	return Start(ctx, next)
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
