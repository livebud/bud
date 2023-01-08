package run

import (
	"context"
	"path/filepath"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/internal/budfs"
	"github.com/livebud/bud/internal/config"
	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/prompter"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/shell"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/watcher"
)

// New command for bud run.
func New(provide config.Provide) *Command {
	return &Command{
		provide: provide,
	}
}

// Command for bud run.
type Command struct {
	provide config.Provide
}

// Run the run command. That's a mouthful.
func (c *Command) Run(ctx context.Context) (err error) {
	module, err := c.provide.Module()
	if err != nil {
		return err
	}
	// Ensure we have version alignment between the CLI and the runtime
	if err := versions.AlignRuntime(ctx, module, versions.Bud); err != nil {
		return err
	}
	// Setup the logger
	log, err := c.provide.Logger()
	if err != nil {
		return err
	}
	webln, err := c.provide.WebListener()
	if err != nil {
		return err
	}
	defer webln.Close()
	log.Info("Listening on http://" + webln.Addr().String())
	// Setup the prompter
	prompter, err := c.provide.Prompter()
	if err != nil {
		return err
	}
	// Setup the bud listener
	budln, err := c.provide.BudListener()
	if err != nil {
		return err
	}
	defer budln.Close()
	log.Debug("run: bud server is listening on %s", "http://"+budln.Addr().String())
	// Load the generator filesystem
	budfs, err := c.provide.BudFileSystem()
	if err != nil {
		return err
	}
	defer budfs.Close(ctx)
	// Create a bus if we don't have one yet
	bus := c.provide.Bus()
	// Initialize the bud server
	budServer, err := c.provide.BudServer()
	if err != nil {
		return err
	}
	// Setup the starter command
	starter := c.provide.Command()
	starter.Env = append(starter.Env, "BUD_LISTEN="+budln.Addr().String())
	// Get the file descriptor for the web listener
	webFile, err := webln.File()
	if err != nil {
		return err
	}
	defer webFile.Close()
	// Inject that file into the starter's extrafiles
	extrafile.Inject(&starter.ExtraFiles, &starter.Env, "WEB", webFile)
	// Initialize the app server
	appServer := &appServer{
		dir:      module.Directory(),
		prompter: prompter,
		bus:      bus,
		budfs:    budfs,
		log:      log,
		module:   module,
		starter:  starter,
	}
	// Start the servers
	eg, ctx := errgroup.WithContext(ctx)
	// Start the internal bud server
	eg.Go(func() error { return budServer.Listen(ctx) })
	// Start the internal app server
	eg.Go(func() error { return appServer.Run(ctx) })
	// Wait until either the hot or web server exits
	err = eg.Wait()
	log.Field("error", err).Debug("run: command finished")
	return err
}

// appServer runs the generated web application
type appServer struct {
	dir      string
	prompter *prompter.Prompter
	bus      pubsub.Client
	budfs    budfs.FileSystem
	log      log.Log
	module   *gomod.Module
	starter  *shell.Command
}

// Run the app server
func (a *appServer) Run(ctx context.Context) error {
	// Generate the app
	if err := a.budfs.Sync(ctx, a.module); err != nil {
		a.bus.Publish("app:error", []byte(err.Error()))
		a.log.Debug("run: published event %q", "app:error")
		return err
	}
	// Start the built app
	process, err := a.starter.Start(ctx, filepath.Join("bud", "app"))
	if err != nil {
		a.bus.Publish("app:error", []byte(err.Error()))
		a.log.Debug("run: published event %q", "app:error")
		return err
	}
	defer process.Close()
	// Watch for changes
	err = watcher.Watch(ctx, a.dir, catchError(a.prompter, func(events []watcher.Event) error {
		// Trigger reloading
		a.prompter.Reloading(events)
		// Inform the bud filesystem of the changes
		changes := make([]string, len(events))
		for i, event := range events {
			a.log.Debug("run: file path changed %q", event.Path)
			changes[i] = event.Path
		}
		if err := a.budfs.Refresh(ctx, changes...); err != nil {
			return err
		}
		// Check if we can incrementally reload
		if canIncrementallyReload(events) {
			a.log.Debug("run: incrementally reloading")
			// Publish the frontend:update event
			a.bus.Publish("frontend:update", nil)
			a.log.Debug("run: published event %q", "frontend:update")
			// Publish the app:ready event
			a.bus.Publish("app:ready", nil)
			a.log.Debug("run: published event %q", "app:ready")
			a.prompter.SuccessReload()
			return nil
		}
		now := time.Now()
		a.log.Debug("run: restarting the process")
		if err := process.Close(); err != nil {
			return err
		}
		a.bus.Publish("backend:update", nil)
		a.log.Debug("run: published event %q", "backend:update")
		// Generate the app
		if err := a.budfs.Sync(ctx, a.module); err != nil {
			return err
		}
		// Restart the process
		p, err := process.Restart(ctx)
		if err != nil {
			a.bus.Publish("app:error", nil)
			a.log.Debug("run: published event %q", "app:error")
			return err
		}
		a.prompter.SuccessReload()
		a.log.Debug("restarted the process in %s", time.Since(now))
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
	return nil
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
