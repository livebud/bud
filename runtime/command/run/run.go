package run

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/livebud/bud/internal/buildcache"
	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/package/watcher"

	"github.com/livebud/bud/package/hot"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/runtime/command"
	"github.com/livebud/bud/runtime/command/run/prompter"
	"github.com/livebud/bud/runtime/web"
)

func New(module *gomod.Module) *Command {
	return &Command{
		module: module,
		// Default flags
		Flag: &command.Flag{
			Embed:  false,
			Hot:    ":35729",
			Minify: false,
		},
	}
}

// Command to run the project at runtime
type Command struct {
	module *gomod.Module
	// Below are filled in by the CLI
	FS     *overlay.FileSystem
	Flag   *command.Flag
	Listen string
}

func (c *Command) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	// Start the live reload server
	hotServer := hot.New()
	eg.Go(func() error { return c.startHot(ctx, hotServer) })
	// Start the app server
	eg.Go(func() error { return c.startApp(ctx, hotServer) })
	// Wait until either the hot or web server exists
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// The prompt need stdout and stderr's content.
var (
	prompt prompter.Prompter
	stdout = io.MultiWriter(os.Stdout, &prompt.StdOut)
	stderr = io.MultiWriter(os.Stderr, &prompt.StdErr)
)

func (c *Command) compileAndStart(ctx context.Context, listener socket.Listener) (*exe.Cmd, error) {
	// Sync the application
	if err := c.FS.Sync("bud/.app"); err != nil {
		return nil, err
	}
	// Ensure that main.go exists
	if _, err := fs.Stat(c.module, "bud/.app/main.go"); err != nil {
		return nil, err
	}
	// Build the application binary
	bcache := buildcache.Default(c.module)
	if err := bcache.Build(ctx, c.module, "bud/.app/main.go", "bud/app"); err != nil {
		return nil, err
	}

	// Start the web server
	process := exe.Command(ctx, "bud/app")
	process.Stdout = stdout
	process.Stderr = stderr
	process.Env = os.Environ()
	process.Dir = c.module.Directory()

	// Forward V8 read-write pipes to bud/app
	extrafile.Forward(&process.ExtraFiles, &process.Env, "V8")

	// Forward APP listener to bud/app
	fileListener, err := listener.File()
	if err != nil {
		return nil, err
	}
	extrafile.Inject(&process.ExtraFiles, &process.Env, "APP", fileListener)

	if err := process.Start(); err != nil {
		return nil, err
	}
	return process, nil
}

func (c *Command) startApp(ctx context.Context, hotServer *hot.Server) error {
	listener, err := web.Listen("APP", c.Listen)
	if err != nil {
		return err
	}

	// Create a terminal prompter
	prompt.Init(web.Format(listener))

	// Compile and start the project
	process, err := c.compileAndStart(ctx, listener)
	if err != nil {
		// Exit without logging if the context has been cancelled. This can
		// occur when the hot reload server failed to start or exits early.
		if errors.Is(err, context.Canceled) {
			return err
		}
		console.Error(err.Error())
		// TODO: de-duplicate with the watcher below
		if err := watcher.Watch(ctx, ".", func(paths []string) error {
			prompt.Reloading(paths)
			process, err = c.compileAndStart(ctx, listener)
			if err != nil {
				// Exit without logging if the context has been cancelled. This can
				// occur when the hot reload server failed to start or exits early.
				if errors.Is(err, context.Canceled) {
					return err
				}
				prompt.FailReload(err.Error())
				return nil
			}
			prompt.SuccessReload()
			return watcher.Stop
		}); err != nil {
			return err
		}
		// The watcher has been cancelled before we ever got an active process, so
		// we'll return the original error.
		if process == nil {
			return err
		}
	}
	defer process.Close()
	// Start watching
	if err := watcher.Watch(ctx, ".", func(paths []string) error {
		prompt.Reloading(paths)
		// Check if the changed paths support an incremental reload
		if canIncrementallyReload(paths) {
			// Trigger a reload if there's a hot reload server configured
			if hotServer != nil {
				hotServer.Reload("*")
				prompt.SuccessReload()
			} else {
				prompt.MadeNoReload()
			}
			return nil
		}
		// Otherwise trigger a full reload if there's a hot reload server configured
		if hotServer != nil {
			// Exclamation point just means full page reload
			hotServer.Reload("!")
		}
		if err := process.Close(); err != nil {
			prompt.FailReload(err.Error())
			return nil
		}
		p, err := c.compileAndStart(ctx, listener)
		if err != nil {
			prompt.FailReload(err.Error())
			return nil
		}
		process = p
		prompt.SuccessReload()
		return nil

	}); err != nil {
		return err
	}
	return process.Wait()
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

func (c *Command) startHot(ctx context.Context, hotServer *hot.Server) error {
	if c.Flag.Hot == "" {
		return nil
	}
	listener, err := web.Listen("HOT", c.Flag.Hot)
	if err != nil {
		return err
	}
	return web.Serve(ctx, listener, hotServer)
}
