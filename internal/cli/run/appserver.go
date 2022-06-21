package run

import (
	"context"
	"fmt"
	"path/filepath"

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
	goFileSub := a.bus.Subscribe("file:update:go", "file:create:go", "file:delete:go")
	defer goFileSub.Close()
	for {
		select {
		case <-ctx.Done():
			// Wait for the command to exit
			return process.Wait()
		case <-goFileSub.Wait():
			fmt.Println("triggering a restart!")
		}
	}
}
