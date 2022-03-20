package bud

import (
	"os"
	"os/exec"
	"strings"
)

type Process struct {
	cmd *exec.Cmd
}

func (p *Process) Close() error {
	sp := p.cmd.Process
	if sp != nil {
		if err := sp.Signal(os.Interrupt); err != nil {
			sp.Kill()
		}
	}
	if err := p.cmd.Wait(); err != nil {
		return err
	}
	return nil
}

func (p *Process) Wait() error {
	return p.cmd.Wait()
}

func isExitStatus(err error) bool {
	return err != nil && strings.Contains(err.Error(), "exit status ")
}

func (p *Process) Restart() error {
	sp := p.cmd.Process
	if sp != nil {
		if err := sp.Signal(os.Interrupt); err != nil {
			sp.Kill()
		}
	}
	if err := p.cmd.Wait(); err != nil {
		if !isExitStatus(err) {
			return err
		}
	}
	// Re-run the command again. cmd.Args[0] is the path, so we skip that.
	cmd := exec.Command(p.cmd.Path, p.cmd.Args[1:]...)
	cmd.Env = p.cmd.Env
	cmd.Stdout = p.cmd.Stdout
	cmd.Stderr = p.cmd.Stderr
	cmd.Stdin = p.cmd.Stdin
	cmd.ExtraFiles = p.cmd.ExtraFiles
	cmd.Dir = p.cmd.Dir
	if err := cmd.Start(); err != nil {
		return err
	}
	// Point to the new command
	p.cmd = cmd
	return nil
}
