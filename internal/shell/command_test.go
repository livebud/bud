package shell_test

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/shell"
	"github.com/livebud/bud/internal/sig"
	"github.com/livebud/bud/internal/testsub"
)

func TestCancel(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		ctx, cancel := context.WithCancel(context.Background())
		p, err := shell.Start(ctx, cmd)
		is.NoErr(err)
		time.Sleep(100 * time.Millisecond)
		cancel()
		is.NoErr(p.Wait())
	}
	child := func(t testing.TB) {
		ctx := sig.Trap(context.Background(), os.Interrupt)
		<-ctx.Done()
	}
	testsub.Run(t, parent, child)
}

func TestDone(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		p, err := shell.Start(context.Background(), cmd)
		is.NoErr(err)
		is.NoErr(p.Wait())
	}
	child := func(t testing.TB) {
		time.Sleep(100 * time.Millisecond)
	}
	testsub.Run(t, parent, child)
}

func TestDoneError(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		p, err := shell.Start(context.Background(), cmd)
		is.NoErr(err)
		err = p.Wait()
		is.True(err != nil)
		is.Equal(err.Error(), "exit status 1")
	}
	child := func(t testing.TB) {
		time.Sleep(100 * time.Millisecond)
		os.Exit(1)
	}
	testsub.Run(t, parent, child)
}

func TestRun(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		ctx := context.Background()
		is.NoErr(shell.Run(ctx, cmd))
	}
	child := func(t testing.TB) {
		time.Sleep(100 * time.Millisecond)
	}
	testsub.Run(t, parent, child)
}

func TestRunError(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		ctx := context.Background()
		is.NoErr(shell.Run(ctx, cmd))
	}
	child := func(t testing.TB) {
		time.Sleep(100 * time.Millisecond)
	}
	testsub.Run(t, parent, child)
}

func TestCloseWaitWait(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		ctx := context.Background()
		p, err := shell.Start(ctx, cmd)
		is.NoErr(err)
		is.NoErr(p.Close())
		is.NoErr(p.Wait())
		is.NoErr(p.Wait())
	}
	child := func(t testing.TB) {
		ctx := sig.Trap(context.Background(), os.Interrupt)
		<-ctx.Done()
	}
	testsub.Run(t, parent, child)
}
