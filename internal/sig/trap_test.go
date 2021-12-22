package sig_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/sig"
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
	ctx, cancel := sig.Trap(context.Background(), os.Interrupt)
	defer cancel()
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

func TestCancel(t *testing.T) {
	is := is.New(t)
	ctx, cancel := sig.Trap(context.Background(), os.Interrupt)
	// Should not have received signal
	select {
	case <-ctx.Done():
		is.Fail() // context shouldn't be cancelled yet
	default:
	}
	cancel()
	// Should have received a signal
	select {
	case <-ctx.Done():
	default:
		is.Fail() // context should have been cancelled
	}
}
