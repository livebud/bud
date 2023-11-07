package app

import (
	"context"
	"os"

	"github.com/livebud/bud/cli"
	"github.com/livebud/bud/log"
	"github.com/livebud/bud/mux"
)

// Run a program, returning an exit code.
func Run(cli cli.Parser) int {
	ctx := context.Background()
	if err := Parse(ctx, cli, os.Args[1:]...); err != nil {
		log.Error(err.Error())
		return 1
	}
	return 0
}

// Parse runs a program, returning an error if there is one.
func Parse(ctx context.Context, cli cli.Parser, args ...string) error {
	return cli.Parse(ctx, args...)
}

type Config struct {
	Router func() *mux.Router
}
