package run

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/exe"
	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/gobuild"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/socket"
)

func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud:  bud,
		in:   in,
		Flag: new(framework.Flag),
	}
}

type Command struct {
	bud *bud.Command
	in  *bud.Input

	// Flags
	Flag   *framework.Flag
	Listen string // Web listener address
}

func (c *Command) Run(ctx context.Context) (err error) {
	// Find go.mod
	module, err := bud.Module(c.bud.Dir)
	if err != nil {
		return err
	}
	// Setup the logger
	log, err := bud.Log(c.in.Stderr, c.bud.Log)
	if err != nil {
		return err
	}
	// Listening on the web listener as soon as possible
	webln := c.in.WebLn
	if webln == nil {
		webln, err = socket.Listen(c.Listen)
		if err != nil {
			return err
		}
		log.Info("Listening on http://" + webln.Addr().String())
	}
	// Setup the bud listener
	budln := c.in.BudLn
	if budln == nil {
		budln, err = socket.Listen(":35729")
		if err != nil {
			return err
		}
		log.Debug("run: bud server is listening", "url", "http://"+budln.Addr().String())
	}
	// Load the generator filesystem
	genfs, err := bud.FileSystem(log, module, c.Flag)
	if err != nil {
		return err
	}
	// Load the file server
	servefs, err := bud.FileServer(log, module, c.Flag)
	if err != nil {
		return err
	}
	// Create a bus if we don't have one yet
	bus := c.in.Bus
	if bus == nil {
		bus = pubsub.New()
	}
	// Initialize the bud server
	budServer := &budServer{
		budln: budln,
		bus:   bus,
		fsys:  servefs,
		log:   log,
	}
	// Setup the starter command
	starter := &exe.Command{
		Stdin:  c.in.Stdin,
		Stdout: c.in.Stdout,
		Stderr: c.in.Stderr,
		Dir:    module.Directory(),
		Env: append(c.in.Env,
			"BUD_LISTEN="+budln.Addr().String(),
		),
	}
	// Get the file descriptor for the web listener
	webFile, err := webln.File()
	if err != nil {
		return err
	}
	// Inject that file into the starter's extrafiles
	extrafile.Inject(&starter.ExtraFiles, &starter.Env, "WEB", webFile)
	// Initialize the app server
	appServer := &appServer{
		builder: gobuild.New(module),
		bus:     bus,
		genfs:   genfs,
		log:     log,
		starter: starter,
	}
	// Start watching the filesystem
	watchfs := &watchfs{
		bus: bus,
		dir: module.Directory(),
		log: log,
	}
	// Start the servers
	eg, ctx := errgroup.WithContext(ctx)
	// Start the internal bud server
	eg.Go(func() error { return budServer.Run(ctx) })
	// Start the internal app server
	eg.Go(func() error { return appServer.Run(ctx) })
	// Start the watcher
	eg.Go(func() error { return watchfs.Run(ctx) })
	// Wait until either the hot or web server exits
	err = eg.Wait()
	log.Debug("run: command finished", "err", err)
	return err
}

// // 1. Trigger reload
// // 2. Close existing process
// // 3. Generate new codebase
// // 4. Start new process
// func startApp(ctx context.Context, genfs *overlay.FileSystem, log log.Interface, module *gomod.Module, webln socket.Listener) (err error) {
// 	// Generate the app
// 	if err := genfs.Sync("bud/internal/app"); err != nil {
// 		return err
// 	}
// 	cmd := exec.Command(filepath.Join("bud", "app"))
// 	// cmd.Stdin = stdin
// 	// cmd.Stdout = stdout
// 	// cmd.Stderr = stderr
// 	// cmd.Env = env
// 	// // Run always runs the bud listener. This allows the app to connect to the bud
// 	// // server.
// 	// cmd.Env = append(cmd.Env, "BUD_LISTEN="+budln.Addr().String())
// 	cmd.Dir = module.Directory()
// 	// Inject the web listener into the app
// 	webFile, err := webln.File()
// 	if err != nil {
// 		return err
// 	}
// 	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "WEB", webFile)
// 	// Start the command
// 	if err := cmd.Start(); err != nil {
// 		return err
// 	}

// 	// Run the start function once upon booting
// 	if err := c.restart(ctx, genfs, log, module); err != nil {
// 		log.Error(err.Error())
// 	}
// 	// Watch the project
// 	err = watcher.Watch(ctx, module.Directory(), func(paths []string) error {
// 		log.Debug("run: files changed", "paths", paths)
// 		if err := c.restart(ctx, genfs, log, module, paths...); err != nil {
// 			log.Error(err.Error())
// 		}
// 		return nil
// 	})
// 	log.Debug("run: watcher closed", "err", err)
// 	if c.app != nil {
// 		err := closeProcess(c.app)
// 		log.Debug("run: app server closed", "err", err)
// 		return err
// 	}
// 	return nil
// }

// func restart(ctx context.Context, genfs *overlay.FileSystem, log log.Interface, module *gomod.Module, updatePaths ...string) (err error) {
// 	if c.app != nil {
// 		if canIncrementallyReload(updatePaths) {
// 			// Trigger an incremental reload.
// 			log.Debug("run: publishing event", "topic", "page:update:*", "paths", updatePaths)
// 			c.Bus.Publish("page:update:*", nil)
// 			return nil
// 		}
// 		// Reload the full server. Exclamation point just means full page reload.
// 		log.Debug("run: publishing event", "topic", "page:reload", "paths", updatePaths)
// 		c.Bus.Publish("page:reload", nil)
// 		if err := closeProcess(c.app); err != nil {
// 			return err
// 		}
// 	}
// 	// Generate the app
// 	if err := genfs.Sync("bud/internal/app"); err != nil {
// 		return err
// 	}
// 	// Build the app
// 	// if err := c.bud.Build(ctx, module, "bud/internal/app/main.go", "bud/app"); err != nil {
// 	// 	return err
// 	// }
// 	// Start the app
// 	cmd := exec.Command(filepath.Join("bud", "app"))
// 	// cmd.Stdin = c.bud.Stdin
// 	// cmd.Stdout = c.bud.Stdout
// 	// cmd.Stderr = c.bud.Stderr
// 	// cmd.Env = c.bud.Env
// 	// Run always runs the bud listener. This allows the app to connect to the bud
// 	// server.
// 	cmd.Env = append(cmd.Env, "BUD_LISTEN="+c.BudLn.Addr().String())
// 	cmd.Dir = module.Directory()
// 	// Inject the web listener into the app
// 	webFile, err := c.WebLn.File()
// 	if err != nil {
// 		return err
// 	}
// 	extrafile.Inject(&cmd.ExtraFiles, &cmd.Env, "WEB", webFile)
// 	// Start the command
// 	if err := cmd.Start(); err != nil {
// 		return err
// 	}
// 	go watchProcess(c.Bus, cmd)
// 	c.app = cmd
// 	return nil
// }

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
