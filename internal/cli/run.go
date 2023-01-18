package cli

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/livebud/bud/internal/dag"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/prompter"
	"github.com/livebud/bud/package/watcher"
)

type Run struct {
	Flag      *framework.Flag
	ListenAFS string
	ListenDev string
	ListenWeb string
	Watch     bool
}

func (c *CLI) Run(ctx context.Context, in *Run) error {
	// Find the module
	module, err := c.findModule()
	if err != nil {
		return err
	}
	log, err := c.loadLog()
	if err != nil {
		return err
	}
	webLn, err := c.listenWeb(in.ListenWeb)
	if err != nil {
		return err
	}
	log.Info("Listening on http://" + webLn.Addr().String())

	webFile, err := c.listenWebFile(webLn)
	if err != nil {
		return err
	}

	afsLn, err := c.listenAFS(in.ListenAFS)
	if err != nil {
		return err
	}

	devLn, err := c.listenDev(in.ListenDev)
	if err != nil {
		return err
	}
	bus := c.bus()
	db, err := c.openDB(log, module)
	if err != nil {
		return err
	}
	generate := &Generate{
		Flag:      in.Flag,
		ListenAFS: in.ListenAFS,
		ListenDev: in.ListenDev,
	}
	if err := c.Generate(ctx, generate); err != nil {
		return err
	}

	// Start the application process
	appProcess, err := c.startApp(ctx, module, afsLn, devLn, webFile)
	if err != nil {
		return err
	}

	if !in.Watch {
		return appProcess.Wait()
	}

	prompter := c.prompter(webLn)

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
		if err := db.Delete(changes...); err != nil {
			if errors.Is(err, dag.ErrDatabaseMoved) {
				log.Error("run: bud/bud.db no longer exists")
				return watcher.Stop
			}
			return err
		}
		// Check if we can incrementally reload
		if canIncrementallyReload(events) {
			log.Debug("run: incrementally reloading")
			// Publish the frontend:update event
			bus.Publish("frontend:update", nil)
			log.Debug("run: published event %q", "frontend:update")
			// Publish the app:ready event
			bus.Publish("app:ready", nil)
			log.Debug("run: published event %q", "app:ready")
			prompter.SuccessReload()
			return nil
		}
		now := time.Now()
		log.Debug("run: restarting the process")
		if err := appProcess.Close(); err != nil {
			return err
		}
		bus.Publish("backend:update", nil)
		log.Debug("run: published event %q", "backend:update")
		// Generate the app
		if err := c.Generate(ctx, generate); err != nil {
			if errors.Is(err, dag.ErrDatabaseMoved) {
				log.Error("run: bud/bud.db no longer exists")
				return watcher.Stop
			}
			return err
		}
		// Restart the process
		p, err := appProcess.Restart(ctx)
		if err != nil {
			bus.Publish("app:error", nil)
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

// catchError wraps the watch function in a handler that logs the error instead of
// returning the error (and canceling the watcher)
func catchError(prompter *prompter.Prompter, fn func(events []watcher.Event) error) func(events []watcher.Event) error {
	return func(events []watcher.Event) error {
		if err := fn(events); err != nil {
			if errors.Is(err, watcher.Stop) {
				return err
			}
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
