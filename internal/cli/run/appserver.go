package run

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/livebud/bud/internal/exe"
	"github.com/livebud/bud/internal/gobuild"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/overlay"
)

type appServer struct {
	builder *gobuild.Builder
	bus     pubsub.Client
	genfs   *overlay.FileSystem
	log     log.Interface
	starter *exe.Command
}

func (a *appServer) Run(ctx context.Context) error {
	// Generate the app
	if err := a.genfs.Sync("bud/internal/app"); err != nil {
		return err
	}
	// Build the app
	if err := a.builder.Build(ctx, "bud/internal/app/main.go", "bud/app"); err != nil {
		return err
	}
	// Start the built app
	process, err := a.starter.Start(ctx, filepath.Join("bud", "app"))
	if err != nil {
		return err
	}
	// Subscribe to file change events
	sub := a.bus.Subscribe("watch:backend:update")
	defer sub.Close()
	for {
		select {
		case <-ctx.Done():
			// Wait for the command to exit
			return process.Wait()
		case <-sub.Wait():
		}
		fmt.Println("triggering a restart!")
		now := time.Now()
		// Generate the app
		if err := a.genfs.Sync("bud/internal/app"); err != nil {
			return err
		}
		// Build the app
		if err := a.builder.Build(ctx, "bud/internal/app/main.go", "bud/app"); err != nil {
			return err
		}
		p, err := process.Restart(ctx)
		if err != nil {
			return err
		}
		fmt.Println("restarted in", time.Since(now))
		process = p
	}
}
