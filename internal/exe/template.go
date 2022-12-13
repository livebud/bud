package exe

import (
	"io"
	"os"
	"os/exec"
)

type Template struct {
	Dir        string
	Stdin      io.Reader
	Stdout     io.Writer
	Stderr     io.Writer
	Env        []string
	ExtraFiles []*os.File
}

// Turn the template into a command
func (t *Template) Command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Stdin = t.Stdin
	cmd.Stdout = t.Stdout
	cmd.Stderr = t.Stderr
	cmd.Env = t.Env
	cmd.Dir = t.Dir
	cmd.ExtraFiles = t.ExtraFiles
	return cmd
}
