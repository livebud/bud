package run

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/package/devserver"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/watcher"
	"github.com/livebud/bud/runtime/web"
	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/socket"
)

func New(bud *bud.Command, bus pubsub.Client, webListener, budListener socket.Listener) *Command {
	return &Command{
		bud:         bud,
		bus:         bus,
		webListener: webListener,
		budListener: budListener,
		Flag:        new(framework.Flag),
	}
}

type Command struct {
	bud *bud.Command
	bus pubsub.Client

	// Passed in for testing
	webListener socket.Listener // Can be nil
	budListener socket.Listener // Can be nil

	// Flags
	Flag   *framework.Flag
	Listen string // Web listen address

	// Private
	app *exec.Cmd // Starts as nil
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

	// Start serving bud
	if c.budListener == nil {
		c.budListener, err = socket.Listen(":35729")
		if err != nil {
			return err
		}
		log.Debug("Bud server is listening on http://" + c.budListener.Addr().String())
	}

	// Load the filesystem
	genfs, err := c.bud.FileSystem(log, module, c.Flag)
	if err != nil {
		return err
	}

	servefs, err := c.bud.FileServer(log, module, c.Flag)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Start the internal bud server
	eg.Go(func() error {
		return c.startBud(ctx, servefs, log)
	})

	// Start the app
	eg.Go(func() error {
		return c.startApp(ctx, genfs, log, module)
	})

	// Wait until either the hot or web server exits
	err = eg.Wait()
	log.Debug("run: command finished", "err", err)
	return err
}

func (c *Command) startBud(ctx context.Context, servefs *overlay.Server, log log.Interface) (err error) {
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	devServer := devserver.New(servefs, c.bus, log, vm)
	err = web.Serve(ctx, c.budListener, devServer)
	log.Debug("run: bud server closed", "err", err)
	return err
}

// 1. Trigger reload
// 2. Close existing process
// 3. Generate new codebase
// 4. Start new process
func (c *Command) startApp(ctx context.Context, genfs *overlay.FileSystem, log log.Interface, module *gomod.Module) (err error) {
	if c.webListener == nil {
		c.webListener, err = socket.Listen(c.Listen)
		if err != nil {
			return err
		}
		log.Info("Listening on http://localhost" + c.Listen)
	}
	// Run the start function once upon booting
	if err := c.restart(ctx, genfs, log, module); err != nil {
		log.Error(err.Error())
	}
	// Watch the project
	err = watcher.Watch(ctx, module.Directory(), func(paths []string) error {
		log.Debug("run: files changed", "paths", paths)
		if err := c.restart(ctx, genfs, log, module, paths...); err != nil {
			log.Error(err.Error())
		}
		return nil
	})
	log.Debug("run: watcher closed", "err", err)
	if c.app != nil {
		err := closeProcess(c.app)
		log.Debug("run: app server closed", "err", err)
		return err
	}
	return nil
}

func (c *Command) restart(ctx context.Context, genfs *overlay.FileSystem, log log.Interface, module *gomod.Module, updatePaths ...string) (err error) {
	if c.app != nil {
		if canIncrementallyReload(updatePaths) {
			log.Debug("run: incrementally reloading")
			// Trigger an incremental reload. Star just means any path.
			c.bus.Publish("page:update:*", nil)
			return nil
		}
		// Reload the full server. Exclamation point just means full page reload.
		log.Debug("run: reloading the page")
		c.bus.Publish("page:reload", nil)
		if err := closeProcess(c.app); err != nil {
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
	app, err := c.bud.Start(module, c.webListener, c.budListener, c.Flag)
	if err != nil {
		return err
	}
	go watchProcess(c.bus, app)
	c.app = app
	return nil
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

// watchProcess watches for a process to exit and publishes an event if there
// was an error.
func watchProcess(bus pubsub.Publisher, cmd *exec.Cmd) error {
	if err := cmd.Wait(); err != nil {
		if !isWaitError(err) {
			bus.Publish("cmd:error", []byte(err.Error()))
			return err
		}
	}
	return nil
}

// Close the process down gracefully
func closeProcess(cmd *exec.Cmd) error {
	sp := cmd.Process
	if sp == nil {
		return nil
	}
	if err := sp.Signal(os.Interrupt); err != nil {
		if errors.Is(err, os.ErrProcessDone) {
			return nil
		}
		if err := sp.Kill(); err != nil {
			return err
		}
	}
	if err := cmd.Wait(); err != nil {
		if !isWaitError(err) {
			return err
		}
	}
	return nil
}

func isWaitError(err error) bool {
	return err != nil && strings.Contains(err.Error(), "Wait was already called")
}
