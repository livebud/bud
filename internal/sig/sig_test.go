package sig_test

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/livebud/bud/internal/testsub"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/sig"
)

func waitFor(r io.Reader, line string) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		if scanner.Text() == line {
			return scanner.Err()
		}
		fmt.Println(scanner.Text())
	}
	return fmt.Errorf("unable to find line %q", line)
}

func TestInterrupt(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		stdin, err := cmd.StdinPipe()
		is.NoErr(err)
		stderr, err := cmd.StderrPipe()
		is.NoErr(err)
		cmd.Stdout = os.Stdout
		is.NoErr(cmd.Start())
		is.NoErr(waitFor(stderr, "ready"))
		is.NoErr(cmd.Process.Signal(os.Interrupt))
		stdin.Write([]byte("interrupt"))
		stdin.Close()
		is.NoErr(cmd.Wait())
	}
	child := func(t testing.TB) {
		ctx := sig.Trap(context.Background(), os.Interrupt)
		// Should not have received signal
		select {
		case <-ctx.Done():
			is.Fail() // context shouldn't be cancelled yet
		default:
		}
		fmt.Fprintln(os.Stderr, "ready")
		waitFor(os.Stdin, "interrupt")
		// Should have received a signal
		select {
		case <-ctx.Done():
			// Give it 1 second to receive the signal
		case <-time.Tick(time.Second):
			is.Fail() // context should have been cancelled
		}
	}
	testsub.Run(t, parent, child)
}

func TestEither(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		stdin, err := cmd.StdinPipe()
		is.NoErr(err)
		stderr, err := cmd.StderrPipe()
		is.NoErr(err)
		cmd.Stdout = os.Stdout
		is.NoErr(cmd.Start())
		is.NoErr(waitFor(stderr, "ready"))
		is.NoErr(cmd.Process.Signal(syscall.SIGQUIT))
		stdin.Write([]byte("interrupt"))
		stdin.Close()
		is.NoErr(cmd.Wait())
	}
	child := func(t testing.TB) {
		ctx := sig.Trap(context.Background(), os.Interrupt, syscall.SIGQUIT)
		// Should not have received signal
		select {
		case <-ctx.Done():
			is.Fail() // context shouldn't be cancelled yet
		default:
		}
		fmt.Fprintln(os.Stderr, "ready")
		waitFor(os.Stdin, "interrupt")
		// Should have received a SIGQUIT
		select {
		case <-ctx.Done():
		case <-time.Tick(time.Second):
			is.Fail() // context should have been cancelled
		}
	}
	testsub.Run(t, parent, child)
}

func TestMultiple(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		cmd.Stdout = os.Stdout
		stdin, err := cmd.StdinPipe()
		is.NoErr(err)
		stderr, err := cmd.StderrPipe()
		is.NoErr(err)
		is.NoErr(cmd.Start())
		is.NoErr(waitFor(stderr, "ready"))
		is.NoErr(cmd.Process.Signal(os.Interrupt))
		stdin.Write([]byte("interrupt"))
		stdin.Close()
		is.NoErr(cmd.Wait())
	}
	child := func(t testing.TB) {
		ctx1 := sig.Trap(context.Background(), os.Interrupt)
		ctx2 := sig.Trap(context.Background(), os.Interrupt)
		// Should not have received signal
		select {
		case <-ctx1.Done():
			t.Fatalf("context should not be cancelled yet")
		case <-ctx2.Done():
			t.Fatalf("context should not be cancelled yet")
		default:
		}
		fmt.Fprintln(os.Stderr, "ready")
		waitFor(os.Stdin, "interrupt")
		// Should have received a interrupt
		select {
		case <-ctx1.Done():
		case <-time.Tick(time.Second):
			t.Fatalf("context should have been cancelled")
		}
		select {
		case <-ctx2.Done():
		case <-time.Tick(time.Second):
			t.Fatalf("context should have been cancelled")
		}
	}
	testsub.Run(t, parent, child)
}
