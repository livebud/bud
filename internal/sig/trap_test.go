package sig_test

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/sig"
)

// Based on: https://github.com/golang/go/issues/19326
func raise(sig os.Signal) error {
	p, err := os.FindProcess(os.Getpid())
	if err != nil {
		return err
	}
	return p.Signal(sig)
}

func TestInterrupt(t *testing.T) {
	is := is.New(t)
	ctx := sig.Trap(context.Background(), os.Interrupt)
	// Should not have received signal
	select {
	case <-ctx.Done():
		is.Fail() // context shouldn't be cancelled yet
	default:
	}
	is.NoErr(raise(os.Interrupt))
	// Should have received a signal
	select {
	case <-ctx.Done():
	// Give it 1 second to receive the signal
	case <-time.Tick(time.Second):
		is.Fail() // context should have been cancelled
	}
}

func TestEither(t *testing.T) {
	is := is.New(t)
	ctx := sig.Trap(context.Background(), os.Interrupt, syscall.SIGQUIT)
	// Should not have received signal
	select {
	case <-ctx.Done():
		is.Fail() // context shouldn't be cancelled yet
	default:
	}
	is.NoErr(raise(syscall.SIGQUIT))
	// Should have received a SIGQUIT
	select {
	case <-ctx.Done():
	case <-time.Tick(time.Second):
		is.Fail() // context should have been cancelled
	}
}
