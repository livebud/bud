package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/livebud/bud/internal/current"

	"github.com/livebud/bud/internal/sig"
	"github.com/livebud/bud/package/watcher"
)

func run(ctx context.Context) error {
	dirname, err := current.Directory()
	if err != nil {
		return err
	}
	ctx = sig.Trap(ctx, os.Interrupt)
	return watcher.Watch(ctx, dirname, func(events []watcher.Event) error {
		fmt.Println("-> triggered", events)
		return nil
	})
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}
