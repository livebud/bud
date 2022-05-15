package run

import (
	"context"
	"net"
	"path/filepath"

	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/log/console"
	"github.com/livebud/bud/package/watcher"

	"github.com/livebud/bud/package/hot"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/runtime/bud"
)

type Command struct {
	Flag    *bud.Flag
	Project *bud.Project
	Port    string
}

func (c *Command) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)
	// Initialize the hot server
	var hotServer *hot.Server
	if c.Flag.Hot {
		hotServer = hot.New()
		// Start the hot reload server
		eg.Go(func() error { return c.startHot(ctx, hotServer) })
	}
	// Start the web server
	eg.Go(func() error { return c.startApp(ctx, hotServer) })
	return eg.Wait()
}

func (c *Command) compileAndStart(ctx context.Context, ln net.Listener) (*exe.Cmd, error) {
	app, err := c.Project.Compile(ctx, c.Flag)
	if err != nil {
		return nil, err
	}
	process, err := app.Start(ctx, ln)
	if err != nil {
		return nil, err
	}
	return process, nil
}

func (c *Command) startApp(ctx context.Context, hotServer *hot.Server) error {
	listener, err := socket.Load(c.Port)
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
			console.Info("Ready on http://127.0.0.1" + c.Port)
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
			app, err := c.Project.Compile(ctx, c.Flag)
			if err != nil {
				console.Error(err.Error())
				return nil
			}
			process, err = app.Start(ctx, listener)
			if err != nil {
				console.Error(err.Error())
				return nil
			}
			// TODO: host should be dynamic
			console.Info("Ready on http://127.0.0.1" + c.Port)
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
	// TODO: host should be dynamic
	return hotServer.ListenAndServe(ctx, "127.0.0.1:35729")
}
