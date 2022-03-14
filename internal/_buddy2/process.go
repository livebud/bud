package buddy

import (
	"os/exec"
)

type Process struct {
	cmd *exec.Cmd
}

func (p *Process) Wait() error {
	return nil
}

func (p *Process) Close() error {
	return nil
}
