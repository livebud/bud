package run

import (
	"context"
	"path/filepath"

	"github.com/livebud/bud/package/exe"
	"github.com/livebud/bud/package/hot"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/watcher"
	"github.com/livebud/bud/runtime/web"
	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/command"
	"github.com/livebud/bud/package/socket"
)

func New(bud *command.Bud, web, hot socket.Listener) *Command {
	return &Command{
		bud:  bud,
		web:  web,
		hot:  hot,
		Flag: new(framework.Flag),
	}
}

type Command struct {
	bud *command.Bud

	// Passed in for testing
	web socket.Listener // Can be nil
	hot socket.Listener // Can be nil

	// Flags
	Flag   *framework.Flag
	Hot    bool   // Enable hot reload
	Listen string // Web listen address

	// Private
	app *exe.Cmd // Starts as nil
}

func (c *Command) Run(ctx context.Context) (err error) {
	log, err := c.bud.Logger()
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	hotServer := hot.New()

	// Setup and start the hot reload server
	if c.Hot {
		if c.hot == nil {
			// Listen on any free TCP port
			c.hot, err = socket.Listen(":0")
			if err != nil {
				return err
			}
		}
		// Set the framework flag to the address of the hot server
		c.Flag.Hot = c.hot.Addr().String()
		// Start serving requests from the listener with the hot server
		eg.Go(func() error {
			return web.Serve(ctx, c.hot, hotServer)
		})
	}

	// Start the app
	eg.Go(func() error {
		return c.startApp(ctx, log, hotServer)
	})

	// Wait until either the hot or web server exits
	return eg.Wait()

	// { // 1. Listen
	// 	// Bind the web listener to the listen address
	// 	if c.web == nil && !disabled(c.Listen) {
	// 		c.web, err = socket.Listen(c.Listen)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		defer c.web.Close()
	// 	}
	// 	// Bind the hot listener to the hot address
	// 	if c.hot == nil && !disabled(c.Flag.Hot) {
	// 		c.hot, err = socket.Listen(c.Flag.Hot)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		defer c.hot.Close()
	// 	}

	// 	// TODO: Setup the Bud pipe
	// }

	// 1. Trigger reload
	// 2. Close existing process
	// 3. Generate new codebase
	// 4. Start new process

	// Start watching
	// return c.bud.Watch(ctx, module, log, func(isBooting, canHotReload bool) error {
	// 	if isBooting {

	// 		fmt.Println("BOOTING")
	// 		return nil
	// 	}

	// 	if canHotReload {
	// 		hotServer.Reload("*")
	// 		return nil
	// 	}

	// 	// Otherwise trigger a full reload if there's a hot reload server configured
	// 	// Exclamation point just means full page reload
	// 	hotServer.Reload("!")
	// 	// if err := process.Close(); err != nil {
	// 	// console.Error(err.Error())
	// 	// return nil
	// 	// }
	// 	// p, err := c.compileAndStart(ctx, listener)
	// 	// if err != nil {
	// 	// console.Error(err.Error())
	// 	// return nil
	// 	// }
	// 	// process = p
	// 	// console.Info("Ready on " + web.Format(listener))
	// 	return nil
	// })

	// // Generating the app
	// if err := c.bud.Generate(module, c.Flag, "bud/internal/app"); err != nil {
	// 	// TODO: watch and retry
	// 	return err
	// }

	// // { // 2. Compile
	// // 	budCompiler := c.bud.Compiler(log, module)
	// // 	err := budCompiler.Compile(ctx, &compiler.Flag{
	// // 		Minify: c.Minify,
	// // 		Embed:  c.Embed,
	// // 		Hot:    c.Hot,
	// // 	})
	// // 	if err != nil {
	// // 		return err
	// // 	}
	// // }

	// { // 3. Start the application
	// 	cmd := exec.Command("bud/app")
	// 	cmd.Dir = module.Directory()
	// 	cmd.Env = c.bud.Env
	// 	cmd.Stdin = c.bud.Stdin
	// 	cmd.Stdout = c.bud.Stdout
	// 	cmd.Stderr = c.bud.Stderr
	// 	if err := cmd.Start(); err != nil {
	// 		// TODO: watch for changes and retry
	// 		return err
	// 	}
	// }

	// { // 4. Watch for changes, recompile and start.

	// }

	// // 1. Listen
	// // 2. Compile
	// //   a. Generate cli
	// //   	 i. Generate bud/internal/cli
	// //     ii. Build bud/cli
	// //     iii. Run bud/cli
	// //   b. Generate app
	// //     i. Generate bud/internal/app
	// //     ii. Build bud/app
	// // 3. Start
	// // 4. Watch
	// // (...Compile, Start)
	// return nil
}

// func (c *Command) listen() (web, hot socket.Listener, err error) {
// 	return c.WebListener, c.HotListener
// }

// // Bind the web listener to the listen address
// func (c *Command) listenWeb() (socket.Listener, error) {
// 	if c.WebListener != nil {
// 		return c.WebListener, nil
// 	} else if disabled(c.Listen) {
// 		return nil, nil
// 	}
// 	listener, err := socket.Listen(c.Listen)
// 	if err != nil {
// 		return nil, err
// 	}

// }

// func (c *Command) listenHot() (socket.Listener, error) {

// }

func (c *Command) startHot(ctx context.Context) (err error) {
	listener := c.hot
	if listener == nil {
		listener, err = socket.Listen(":0")
		if err != nil {
			return err
		}
	}
	// Serve requests on from listener with the hot server
	return web.Serve(ctx, listener, hot.New())
}

// 1. Trigger reload
// 2. Close existing process
// 3. Generate new codebase
// 4. Start new process
func (c *Command) startApp(ctx context.Context, log log.Interface, hotServer *hot.Server) (err error) {
	webListener := c.web
	if webListener == nil {
		log.Info("Listening on http://localhost" + c.Listen)
		webListener, err = socket.Listen(c.Listen)
		if err != nil {
			return err
		}
	}
	module, err := c.bud.Module()
	if err != nil {
		return err
	}
	var app *exe.Cmd
	starter := logWrap(log, func(paths []string) error {
		if app != nil {
			if canIncrementallyReload(paths) {
				// Trigger an incremental reload. Star just means any path.
				hotServer.Reload("*")
				return nil
			}
			// Reload the full server. Exclamation point just means full page reload.
			hotServer.Reload("!")
			if err := app.Close(); err != nil {
				return err
			}
		}
		// Generate the app
		if err := c.bud.Generate(module, c.Flag, "bud/internal/app"); err != nil {
			return err
		}
		// Build the app
		if err := c.bud.Build(ctx, module, "bud/internal/app/main.go", "bud/app"); err != nil {
			return err
		}
		// Start the app
		app, err = c.bud.Start(module, webListener)
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
