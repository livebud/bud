package goplugin

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
)

// Conn interface
type Conn interface {
	io.ReadWriteCloser
}

type host struct {
	rc  io.ReadCloser
	wc  io.WriteCloser
	cmd *exec.Cmd
}

var _ Conn = (*host)(nil)

func (h *host) Read(p []byte) (int, error) {
	return h.rc.Read(p)
}

func (h *host) Write(p []byte) (int, error) {
	return h.wc.Write(p)
}

func (h *host) Close() error {
	h.rc.Close()
	h.wc.Close()
	return h.cmd.Wait()
}

// Start starts up a command and returns the connection
func Start(cmd string, args ...string) (Conn, error) {
	r1, w1, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	r2, w2, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	command := exec.Command(cmd, args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	// Use another file descriptor so both stdout & stderr
	// continue to work inside the plugin.
	command.Env = append([]string{"FD=3"}, os.Environ()...)
	command.ExtraFiles = []*os.File{r1, w2}

	if err := command.Start(); err != nil {
		return nil, err
	}

	return &host{r2, w1, command}, nil
}

type plugin struct {
	rc io.ReadCloser
	wc io.WriteCloser
}

func (pl *plugin) Read(p []byte) (int, error) {
	return pl.rc.Read(p)
}

func (pl *plugin) Write(p []byte) (int, error) {
	return pl.wc.Write(p)
}

func (pl *plugin) Close() error {
	pl.rc.Close()
	pl.wc.Close()
	return nil
}

// Serve a connection
func Serve(name string) (Conn, error) {
	fdEnv := os.Getenv("FD")
	if fdEnv == "" {
		return nil, fmt.Errorf("%s shouldn't be called directly", name)
	}
	fd, err := strconv.Atoi(fdEnv)
	if err != nil {
		return nil, fmt.Errorf("%s file descriptor env is invalid", name)
	}

	reader := os.NewFile(uintptr(fd), "pipe")
	writer := os.NewFile(uintptr(fd+1), "pipe")

	return &plugin{reader, writer}, nil
}
