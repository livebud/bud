package run

import (
	"context"
	"fmt"
	"net"
	"path/filepath"
	"runtime"
	"time"

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

func deletePreviousLine() {
	// Move cursor and delete line based on ANSI Escape Sequences: https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
	switch runtime.GOOS {
	case "windows":
		// TODO: Test this on Windows, don't know will it work or not
		fmt.Print("\033[1A\033[0K")
	case "linux":
		// Move cursor up and delete line
		fmt.Print("\033[1A\033[0K")
	default:
		console.Error("Can't detect your operating system, please submit an issue at https://github.com/livebud/bud/issues")
	}
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

	totalRun := 0
	timeElapsed := 0
	// Start watching
	if err := watcher.Watch(ctx, ".", func(path string) error {
		// Count how much watcher.Watch got called
		totalRun++
		// Start timer
		startTime := time.Now()
		// Handle "Reloading..." message
		deletePreviousLine()
		console.Info("Reloading...")
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
				console.Error("")
				return nil
			}
			app, err := c.Project.Compile(ctx, c.Flag)
			if err != nil {
				console.Error(err.Error())
				console.Error("")
				return nil
			}
			process, err = app.Start(ctx, listener)
			if err != nil {
				console.Error(err.Error())
				console.Error("")
				return nil
			}
			timeElapsed += int(time.Since(startTime).Milliseconds())
			if totalRun%4 == 0 {
				deletePreviousLine()
				msg := fmt.Sprintf("Ready in %dms (x%d)", timeElapsed, totalRun/4)
				console.Info(msg)
				timeElapsed = 0
			}
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
