package program

import (
	"context"
	"errors"
	"os"

	"github.com/livebud/bud/internal/sig"
	"github.com/livebud/bud/package/log/console"
)

func Run(fn func(ctx context.Context, args ...string) error) int {
	ctx := sig.Trap(context.Background(), os.Interrupt)
	if err := fn(ctx, os.Args[1:]...); err != nil && !errors.Is(err, context.Canceled) {
		console.Error(err.Error())
		return 1
	}
	return 0
}
