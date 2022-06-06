package main

import (
	"context"
	"os"

	cli "github.com/livebud/bud/internal/cli2"
)

//go:generate go run scripts/set-package-json/main.go

// main bud entrypoint. Intentionally simple.
func main() {
	os.Exit(cli.Run(context.Background(), os.Args[1:]...))
}
