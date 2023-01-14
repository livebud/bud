package cli

import (
	"context"
	"path/filepath"
	"time"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/prompter"
	"github.com/livebud/bud/internal/shell"
	"github.com/livebud/bud/package/watcher"
)

type Run struct {
	Closer    *once.Closer
	Flag      *framework.Flag
	ListenAFS string
	ListenDev string
	ListenWeb string
	Watch     bool
}

func (c *CLI) Run(ctx context.Context, in *Run) error {
	// When Run is called directly, it owns the Closer
	if in.Closer == nil {
		in.Closer = new(once.Closer)
		defer in.Closer.Close()
	}
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
	in.Closer.Add(webLn.Close)
	log.Info("Listening on http://" + webLn.Addr().String())

	webFile, err := c.listenFileWeb(webLn)
	if err != nil {
		return err
	}
	in.Closer.Add(webFile.Close)

	afsLn, err := c.listenAFS(in.ListenAFS)
	if err != nil {
		return err
	}
	in.Closer.Add(afsLn.Close)

	devLn, err := c.listenDev(in.ListenDev)
	if err != nil {
		return err
	}
	in.Closer.Add(devLn.Close)
	bus := c.newBus()
	db, err := c.openDB(log, module)
	if err != nil {
		return err
	}
	in.Closer.Add(db.Close)
	generate := &Generate{
		Closer:    in.Closer,
		Flag:      in.Flag,
		ListenAFS: in.ListenAFS,
		ListenDev: in.ListenDev,
		bus:       bus,
		db:        db,
		log:       log,
		module:    module,
		afsLn:     afsLn,
		devLn:     devLn,
	}
	if err := c.Generate(ctx, generate); err != nil {
		return err
	}

	// Start the application process
	appCommand, err := c.loadCommandApp(module, afsLn, devLn, webFile)
	if err != nil {
		return err
	}
	appProcess, err := shell.Start(ctx, appCommand)
	if err != nil {
		return err
	}
	in.Closer.Add(appProcess.Close)

	if !in.Watch {
		return appProcess.Wait()
	}

	prompter := c.newPrompter(webLn)

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
