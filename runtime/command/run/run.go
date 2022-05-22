package run

import (
	"context"
	"io/fs"
	"net"
	"os"
	"path/filepath"

	"github.com/livebud/bud/internal/buildcache"
	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/watcher"

	"github.com/livebud/bud/package/hot"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/runtime/bud"
	"github.com/livebud/bud/runtime/web"
)

func New(module *gomod.Module) *Command {
	return &Command{
		module: module,
		// Default flags
		Flag: &bud.Flag{
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
	Flag   *bud.Flag
	Listen string
}

func (c *Command) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	// Load the hot server
	hotServer := hot.New()
	eg.Go(func() error { return c.startHot(ctx, hotServer) })
	// Start the web server
	eg.Go(func() error { return c.startApp(ctx, hotServer) })
	return eg.Wait()
}

func (c *Command) compileAndStart(ctx context.Context, ln net.Listener) (*exe.Cmd, error) {
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
	process := exe.Command(ctx, "bud/app")
	process.Stdout = os.Stdout
	process.Stderr = os.Stderr
	process.Env = os.Environ()
	process.Dir = c.module.Directory()
	if err := process.Start(); err != nil {
		return nil, err
	}
	return process, nil
}

func (c *Command) startApp(ctx context.Context, hotServer *hot.Server) error {
	listener, err := socket.Load(c.Listen)
	if err != nil {
		return err
	}
	// Compile and start the project
	process, err := c.compileAndStart(ctx, listener)
	if err != nil {
		// TODO: de-duplicate with the watcher above
		console.Error(err.Error())
		if err := watcher.Watch(ctx, ".", func(path string) error {
			process, err = c.compileAndStart(ctx, listener)
			if err != nil {
				console.Error(err.Error())
				return nil
			}
			// TODO: host should be dynamic
			console.Info("Ready on " + listener.Addr().String())
			return watcher.Stop
		}); err != nil {
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
			process, err = c.compileAndStart(ctx, listener)
			if err != nil {
				console.Error(err.Error())
				return nil
			}
			console.Info("Ready on " + formatAddress(listener))
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
	listener, err := socket.Listen(c.Flag.Hot)
	if err != nil {
		return err
	}
	return web.Serve(ctx, listener, hotServer)
}

func formatAddress(l net.Listener) string {
	return l.Addr().String()
}
