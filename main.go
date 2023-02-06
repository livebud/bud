package main

import (
	"context"
	"errors"
	"os"

	cli "github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/errs"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/package/log/console"
)

//go:generate go run scripts/set-package-json/main.go

// main is bud's entrypoint
func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		console.Error(errs.Format(err))
		os.Exit(1)
	}
	os.Exit(0)
}

// Run the CLI with the default configuration and return any resulting errors.
func run(ctx context.Context) error {
	closer := new(once.Closer)
	defer closer.Close()
	// Initialize the CLI
	cli := cli.New(closer)
	// Run the cli
	if err := cli.Parse(ctx, os.Args[1:]...); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}
	return nil
}
