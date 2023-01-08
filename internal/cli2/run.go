package cli

import (
	"context"
	"path/filepath"
	"time"

	"github.com/livebud/bud/internal/prompter"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/watcher"

	v8 "github.com/livebud/bud/package/js/v8"

	"github.com/livebud/bud"
	"golang.org/x/sync/errgroup"
)

func (c *CLI) Run(ctx context.Context, in *bud.Run) error {
	cmd := c.Command.Clone()

	// Find go.mod
	module, err := c.module(cmd)
	if err != nil {
		return err
	}

	// TODO: add config alignment check

	// Setup the logger
	log, err := c.logger()
	if err != nil {
		return err
	}

	// Listen on the web address
	webLn, err := c.listenWeb(in.WebAddress)
	if err != nil {
		return err
	}
	defer webLn.Close()
	log.Info("Listening on http://" + webLn.Addr().String())

	prompter := c.prompter(webLn)

	// Load the database
	db, err := c.openDatabase(log, module)
	if err != nil {
		return err
	}
	defer db.Close()

	// Listen on the dev address
	devLn, err := c.listenDev(in.DevAddress)
	if err != nil {
		return err
	}
	defer devLn.Close()
	cmd.Env = append(cmd.Env, "BUD_DEV_URL="+devLn.Addr().String())
	log.Debug("run: dev server is listening on http://" + devLn.Addr().String())

	vm, err := v8.Load()
	if err != nil {
		return err
	}
	defer vm.Close()

	// Start the file server listener
	fileLn, err := c.listenFile(":0")
	if err != nil {
		return err
	}
	defer fileLn.Close()
	log.Debug("run: file server is listening on http://" + devLn.Addr().String())

	remoteClient, err := remotefs.Dial(ctx, fileLn.Addr().String())
	if err != nil {
		return err
	}
	defer remoteClient.Close()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		err := c.serveDev(ctx, &in.Config, remoteClient, devLn, log, vm)
		return err
	})

	// Run the generator
	// TODO: inline sync along with remotefs
	budfs := c.budfs(&in.Config, cmd, db, fileLn, log, module)
	if err := budfs.Sync(ctx); err != nil {
		return err
	}
	defer budfs.Close()

	// Get the file descriptor for the web listener
	webFile, err := webLn.File()
	if err != nil {
		return err
	}
	defer webFile.Close()
	// Inject that file under the WEB prefix
	cmd.Inject("WEB", webFile)

	if !in.Watch {
		return cmd.Run(ctx, filepath.Join("bud", "app"))
	}

	process, err := cmd.Start(ctx, filepath.Join("bud", "app"))
	if err != nil {
		return err
	}
	defer process.Close()

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
		if err := budfs.Refresh(ctx, changes...); err != nil {
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
			// prompter.SuccessReload()
			return nil
		}
		now := time.Now()
		log.Debug("run: restarting the process")
		if err := process.Close(); err != nil {
			return err
		}
		c.Bus.Publish("backend:update", nil)
		log.Debug("run: published event %q", "backend:update")
		// Generate the app
		if err := budfs.Sync(ctx); err != nil {
			return err
		}
		// Restart the process
		p, err := process.Restart(ctx)
		if err != nil {
			c.Bus.Publish("app:error", nil)
			log.Debug("run: published event %q", "app:error")
			return err
		}
		prompter.SuccessReload()
		log.Debug("restarted the process in %s", time.Since(now))
		process = p
		return nil
	}))
	if err != nil {
		return err
	}
	// Close the final process. This process is most likely different than the
	// deferred process.
	if err := process.Close(); err != nil {
		return err
	}

	return eg.Wait()
}

// logWrap wraps the watch function in a handler that logs the error instead of
// returning the error (and canceling the watcher)
func catchError(prompter *prompter.Prompter, fn func(events []watcher.Event) error) func(events []watcher.Event) error {
	return func(events []watcher.Event) error {
		if err := fn(events); err != nil {
			prompter.FailReload(err.Error())
		}
		return nil
	}
}

// canIncrementallyReload returns true if we can incrementally reload a page
func canIncrementallyReload(events []watcher.Event) bool {
	for _, event := range events {
		if event.Op != watcher.OpUpdate || filepath.Ext(event.Path) == ".go" {
			return false
		}
	}
	return true
}
