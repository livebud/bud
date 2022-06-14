package run

import (
	"context"
	"path/filepath"

	"github.com/livebud/bud/package/devserver"
	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/hot"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/watcher"
	"github.com/livebud/bud/runtime/web"
	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/command"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/socket"
)

func New(bud *command.Bud, webListener, budListener socket.Listener) *Command {
	return &Command{
		bud:         bud,
		webListener: webListener,
		budListener: budListener,
		Flag:        new(framework.Flag),
		hotServer:   hot.New(),
	}
}

type Command struct {
	bud *command.Bud

	// Passed in for testing
	webListener socket.Listener // Can be nil
	budListener socket.Listener // Can be nil

	// Flags
	Flag   *framework.Flag
	Listen string // Web listen address

	// Private
	hotServer *hot.Server
	app       *exe.Cmd // Starts as nil
}

func (c *Command) Run(ctx context.Context) (err error) {
	module, err := c.bud.Module()
	if err != nil {
		return err
	}
	log, err := c.bud.Logger()
	if err != nil {
		return err
	}

	// Load the filesystem
	genfs, err := c.bud.FileSystem(module, c.Flag)
	if err != nil {
		return err
	}

	// Start serving bud
	if c.budListener == nil {
		c.budListener, err = socket.Listen(":0")
		if err != nil {
			return err
		}
		log.Debug("Bud server is listening on http://" + c.budListener.Addr().String())
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Start the bud server
	eg.Go(func() error {
		return c.startBud(ctx, genfs)
	})

	// Start the app
	eg.Go(func() error {
		return c.startApp(ctx, genfs, module, log)
	})

	// Wait until either the hot or web server exits
	return eg.Wait()
}

func (c *Command) startBud(ctx context.Context, genfs *overlay.FileSystem) (err error) {
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	ps := pubsub.New()
	devServer := devserver.New(genfs, ps, vm)
	return web.Serve(ctx, c.budListener, devServer)
}

// 1. Trigger reload
// 2. Close existing process
// 3. Generate new codebase
// 4. Start new process
func (c *Command) startApp(ctx context.Context, genfs *overlay.FileSystem, module *gomod.Module, log log.Interface) (err error) {
	webListener := c.webListener
	if webListener == nil {
		webListener, err = socket.Listen(c.Listen)
		if err != nil {
			return err
		}
		log.Info("Listening on http://localhost" + c.Listen)
	}

	var app *exe.Cmd
	starter := logWrap(log, func(paths []string) error {
		if app != nil {
			if canIncrementallyReload(paths) {
				// Trigger an incremental reload. Star just means any path.
				c.hotServer.Reload("*")
				return nil
			}
			// Reload the full server. Exclamation point just means full page reload.
			c.hotServer.Reload("!")
			if err := app.Close(); err != nil {
				return err
			}
		}
		// Generate the app
		if err := genfs.Sync("bud/internal/app"); err != nil {
			return err
		}
		// Build the app
		if err := c.bud.Build(ctx, module, "bud/internal/app/main.go", "bud/app"); err != nil {
			return err
		}
		// Start the app
		app, err = c.bud.Start(module, webListener, c.budListener)
		if err != nil {
			return err
		}
		return nil
	})
	// Run the start function once upon booting
	if err := starter([]string{}); err != nil {
		return err
	}
	// Watch the project
	return watcher.Watch(ctx, module.Directory(), starter)
}

// Wrap the watch function to allow errors to be returned but logged instead of
// stopping the watcher
func logWrap(log log.Interface, fn func(paths []string) error) func(paths []string) error {
	return func(paths []string) error {
		if err := fn(paths); err != nil {
			log.Error(err.Error())
			return nil
		}
		return nil
	}
}

// canIncrementallyReload returns true if we can incrementally reload a page
func canIncrementallyReload(paths []string) bool {
	for _, path := range paths {
		if filepath.Ext(path) == ".go" {
			return false
		}
	}
	return true
}
