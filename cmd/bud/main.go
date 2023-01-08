package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	cli "github.com/livebud/bud/internal/cli2"
	"github.com/livebud/bud/package/log/console"
)

//go:generate go run ../../scripts/set-package-json/main.go

// main is bud's entrypoint
func main() {
	if err := run(); err != nil {
		console.Error(err.Error())
		fmt.Println("exited with error")
		os.Exit(1)
	}
	fmt.Println("exited successfully")
	os.Exit(0)
}

// Run the CLI with the default configuration and return any resulting errors.
func run() error {
	// Initialize the CLI
	cli := cli.New()
	// Run the cli
	if err := cli.Parse(context.Background(), os.Args[1:]...); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}
	return nil
}
