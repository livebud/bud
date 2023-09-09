package subtest_test

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/livebud/bud/pkg/cli/internal/subtest"
	"github.com/matryer/is"
)

func TestRun(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		cmd.Stdin = strings.NewReader("hello")
		is.NoErr(cmd.Start())
		is.NoErr(cmd.Wait())
	}
	child := func(t testing.TB) {
		is := is.New(t)
		buf, err := io.ReadAll(os.Stdin)
		is.NoErr(err)
		is.Equal(string(buf), "hello")
	}
	subtest.Run(t, parent, child)
}

func TestRunError(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		cmd.Stdin = strings.NewReader("hello")
		is.NoErr(cmd.Start())
		err := cmd.Wait()
		is.True(err != nil)
		is.Equal(err.Error(), "exit status 1")
	}
	child := func(t testing.TB) {
		is := is.New(t)
		buf, err := io.ReadAll(os.Stdin)
		is.NoErr(err)
		is.Equal(string(buf), "helloz")
	}
	subtest.Run(t, parent, child)
}
