package program

import (
	"context"
	"os"

	"github.com/livebud/bud/cli"
	"github.com/livebud/bud/di"
	"github.com/livebud/bud/log"
)

// Load a CLI parser from the dependency injector.
func Load(in di.Injector) (cli.Parser, error) {
	if di.Has[cli.Parser](in) {
		return di.Load[cli.Parser](in)
	}
	return di.Load[*cli.CLI](in)
}

// Parse runs a program, returning an error if there is one.
func Parse(ctx context.Context, in di.Injector, args ...string) error {
	cli, err := Load(in)
	if err != nil {
		return err
	}
	return cli.Parse(ctx, args...)
}

// Run a program, returning an exit code.
func Run(in di.Injector) int {
	ctx := context.Background()
	if err := Parse(ctx, in, os.Args[1:]...); err != nil {
		log.Error(err.Error())
		return 1
	}
	return 0
}
