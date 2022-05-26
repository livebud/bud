package run

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/livebud/bud/internal/buildcache"
	"github.com/livebud/bud/internal/envs"
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
	// Forward existing APP file descriptor to bud/app if it exists.
	files, env, err := extrafile.Prepare("APP", 0, listener)
	if err != nil {
		return nil, err
	}
	// Start the web server
	process := exe.Command(ctx, "bud/app")
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr
	process.Env = envs.From(os.Environ()).Append(env...).List()
	process.Dir = c.module.Directory()
	process.ExtraFiles = append(process.ExtraFiles, files...)
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
	// Compile and start the project
	process, err := c.compileAndStart(ctx, listener)
	if err != nil {
		// Exit without logging if the context has been cancelled. This can
		// occur when the hot reload server failed to start or exits early.
		if errors.Is(err, context.Canceled) {
			return err
		}
		// TODO: de-duplicate with the watcher below
		console.Error(err.Error())
		if err := watcher.Watch(ctx, ".", func(path string) error {
			process, err = c.compileAndStart(ctx, listener)
			if err != nil {
				// Exit without logging if the context has been cancelled. This can
				// occur when the hot reload server failed to start or exits early.
				if errors.Is(err, context.Canceled) {
					return err
				}
				console.Error(err.Error())
				return nil
			}
			console.Info("Ready on " + web.Format(listener))
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
	if err := watcher.Watch(ctx, ".", func(path string) error {
		switch filepath.Ext(path) {
		// Re-compile the app and restart the Go server
		case ".go":
			// Trigger a reload if there's a hot reload server configured
			if hotServer != nil {
				// Exclamation point just means full page reload
				hotServer.Reload("!")
			}
			if err := process.Close(); err != nil {
				console.Error(err.Error())
				return nil
			}
			p, err := c.compileAndStart(ctx, listener)
			if err != nil {
				console.Error(err.Error())
				return nil
			}
			process = p
			console.Info("Ready on " + web.Format(listener))
			return nil
		// Hot reload the page
		default:
			// Trigger a reload if there's a hot reload server configured
			if hotServer != nil {
				hotServer.Reload("*")
			}
			return nil
		}
	}); err != nil {
		return err
	}
	return process.Wait()
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
