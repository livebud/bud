package tester

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/vfs"
)

func New(t testing.TB) *Tester {
	t.Helper()
	is := is.New(t)
	tmpdir := filepath.Join("_tmp", t.Name())
	err := os.MkdirAll(tmpdir, 0755)
	is.NoErr(err)
	t.Cleanup(cleanup(t, "_tmp", tmpdir))
	return &Tester{
		I:    is,
		name: t.Name(),
		env:  []string{"PWD=" + tmpdir},
		dir:  tmpdir,
	}
}

// Cleanup individual files and root if no files left
func cleanup(t testing.TB, root, dir string) func() {
	t.Helper()
	is := is.New(t)
	return func() {
		if t.Failed() {
			return
		}
		is.NoErr(os.RemoveAll(dir))
		fis, err := os.ReadDir(root)
		if err != nil {
			return
		}
		if len(fis) > 0 {
			return
		}
		is.NoErr(os.RemoveAll(root))
	}
}

type Tester struct {
	*is.I
	name string
	env  []string
	dir  string
}

func (t *Tester) Env(key, value string) {
	t.Helper()
	t.env = append(t.env, key+"="+value)
}

func (t *Tester) WriteFiles(files map[string]string) {
	t.Helper()
	t.NoErr(vfs.WriteTo(t.dir, vfs.Map(files)))
}

func (t *Tester) Exists(path string) bool {
	path = filepath.Join(t.dir, path)
	if _, err := os.Stat(path); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false
		}
		t.NoErr(err)
	}
	return true
}

func (t *Tester) WaitFile(path string, deadline time.Duration) {

}

func (t *Tester) WaitPort(port int, deadline time.Duration) {

}

func (t *Tester) Run(command string, args ...string) *CommandResult {
	t.Helper()
	cmd := exec.Command(command, args...)
	cmd.Env = t.env
	cmd.Dir = t.dir
	stdout := new(bytes.Buffer)
	cmd.Stdout = stdout
	stderr := new(bytes.Buffer)
	cmd.Stderr = stderr
	err := cmd.Run()
	return &CommandResult{
		t:      t,
		err:    err,
		Stdout: &Stdio{stdout, t},
		Stderr: &Stdio{stderr, t},
	}
}

type Stdio struct {
	*bytes.Buffer
	t *Tester
}

func (s *Stdio) Equal(expect string) {
	s.t.Equal(s.String(), expect)
}

func (s *Stdio) Match(pattern string) {
	s.t.Match(s.String(), pattern)
}

type CommandResult struct {
	t      *Tester
	err    error
	Stdout *Stdio
	Stderr *Stdio
}

func (c *CommandResult) NoErr() *CommandResult {
	if c.err != nil {
		fmt.Fprint(os.Stderr, c.Stderr.String())
		c.t.NoErr(c.err)
	}
	return c
}

func (c *CommandResult) Error() string {
	if c.err == nil {
		return ""
	}
	return c.err.Error()
}

func (t *Tester) Start(command string, args ...string) (*Command, error) {
	t.Helper()
	cmd := exec.Command(command, args...)
	cmd.Env = t.env
	cmd.Dir = t.dir
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Start()
	t.NoErr(err)
	return &Command{cmd, stdout, stderr}, err
}

type Command struct {
	cmd    *exec.Cmd
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func (c *Command) Stdout() string {
	return c.stdout.String()
}

func (c *Command) Stderr() string {
	return c.stderr.String()
}

func (c *Command) Close() error {
	p := c.cmd.Process
	if p != nil {
		if err := p.Signal(os.Interrupt); err != nil {
			p.Kill()
		}
	}
	return c.cmd.Wait()
}

func (t *Tester) Match(actual, pattern string) {
	t.Helper()
}
