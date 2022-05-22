package main

import (
	"context"
	"os"

	"github.com/livebud/bud/internal/cli"
)

//go:generate go run scripts/set-package-json/main.go

// main bud entrypoint. Intentionally simple.
func main() {
	os.Exit(cli.Parse(context.Background(), os.Args[1:]...))
}
