package main

import (
	"context"
	"errors"
	"os"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/config"
	"github.com/livebud/bud/package/log/console"
)

//go:generate go run scripts/set-package-json/main.go

// main is bud's entrypoint
func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

// Run the CLI with the default configuration and return any resulting errors.
func run() error {
	// Initialize the CLI
	cli := cli.New(config.New())
	// Run the cli
	if err := cli.Run(context.Background(), os.Args[1:]...); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}
	return nil
}
