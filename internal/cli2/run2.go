package cli

import (
	"context"
	"errors"
	"io/fs"
	"time"

	"github.com/livebud/bud"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/package/watcher"
)

func (c *CLI) Run2(ctx context.Context, in *bud.Run2) error {
	provider := newProvider(ctx, c)
	defer provider.Close()
	// budfs := newBudFS(provider)
	// defer budfs.Close()

	log, err := provider.Logger()
	if err != nil {
		return err
	}

	webLn, err := provider.WebListener()
	if err != nil {
		return err
	}
	defer webLn.Close()
	log.Info("Listening on http://" + webLn.Addr().String())

	// Start the dev server
	devServer, err := provider.DevServer()
	if err != nil {
		return err
	}
	go devServer.Listen(ctx)
	defer devServer.Close()

	// Generate the afs
	genfs, err := provider.GenFS()
	if err != nil {
		return err
	}

	module, err := provider.Module()
	if err != nil {
		return err
	}

	db, err := provider.Database()
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Reset(); err != nil {
		return err
	}

	if err := dsync.To(genfs, module, "bud", dsync.WithSkip(afsSkips...), dsync.WithLog(log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	cmd := provider.Command()

	// Build the afs binary
	if err := cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+afsBinPath, afsMainPath); err != nil {
		return err
	}

	// Start the AFS process
	afsCommand, err := provider.AFSCommand()
	if err != nil {
		return err
	}
	afsProcess, err := afsCommand.Start(ctx, module.Directory("bud", "afs"))
	if err != nil {
		return err
	}
	defer afsProcess.Close()

	// Load the remote client
	afsClient, err := provider.AFSClient()
	if err != nil {
		return err
	}
	defer afsClient.Close()

	// Sync the app files again with the remote filesystem
	if err := dsync.To(afsClient, module, "bud", dsync.WithSkip(appSkips...), dsync.WithLog(log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// Build the application binary
	if err := cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+appBinPath, appMainPath); err != nil {
		return err
	}

	// Start the application process
	appCommand, err := provider.AppCommand()
	if err != nil {
		return err
	}
	appProcess, err := appCommand.Start(ctx, module.Directory("bud", "app"))
	if err != nil {
		return err
	}
	defer appProcess.Close()

	if !in.Watch {
		return appProcess.Wait()
	}

	prompter, err := provider.Prompter()
	if err != nil {
		return err
	}

	// Watch for changes
	err = watcher.Watch(ctx, module.Directory(), catchError(prompter, func(events []watcher.Event) error {
		// Trigger reloading
		prompter.Reloading(events)
		// Inform the bud filesystem of the changes
		changes := make([]string, len(events))
		for i, event := range events {
			log.Debug("run: file path changed %q", event.Path)
			changes[i] = event.Path
		}
		// Refresh the cache
		// TODO: optimize later
		if err := db.Reset(); err != nil {
			return err
		}
		// Check if we can incrementally reload
		if canIncrementallyReload(events) {
			log.Debug("run: incrementally reloading")
			// Publish the frontend:update event
			c.Bus.Publish("frontend:update", nil)
			log.Debug("run: published event %q", "frontend:update")
			// Publish the app:ready event
			c.Bus.Publish("app:ready", nil)
			log.Debug("run: published event %q", "app:ready")
			prompter.SuccessReload()
			return nil
		}
		now := time.Now()
		log.Debug("run: restarting the process")
		if err := appProcess.Close(); err != nil {
			return err
		}
		c.Bus.Publish("backend:update", nil)
		log.Debug("run: published event %q", "backend:update")
		// Re-sync the app files again with the remote filesystem
		if err := dsync.To(afsClient, module, "bud", dsync.WithSkip(appSkips...), dsync.WithLog(log)); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}
		// Build the application binary
		if err := cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+appBinPath, appMainPath); err != nil {
			return err
		}
		// Restart the process
		p, err := appProcess.Restart(ctx)
		if err != nil {
			c.Bus.Publish("app:error", nil)
			log.Debug("run: published event %q", "app:error")
			return err
		}
		prompter.SuccessReload()
		log.Debug("restarted the process in %s", time.Since(now))
		appProcess = p
		return nil
	}))
	if err != nil {
		return err
	}

	// Close the final process. This process is most likely different than the
	// deferred process.
	if err := appProcess.Close(); err != nil {
		return err
	}

	return nil
}
