package run

import (
	"context"
	"io"
	"io/fs"
	"net"
	"path/filepath"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/web/webrt"
	"github.com/livebud/bud/internal/cli/bud"
	"github.com/livebud/bud/internal/exe"
	"github.com/livebud/bud/internal/extrafile"
	"github.com/livebud/bud/internal/gobuild"
	"github.com/livebud/bud/internal/prompter"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/budhttp/budsvr"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/overlay"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/package/watcher"
)

// New command for bud run.
func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud:  bud,
		in:   in,
		Flag: new(framework.Flag),
	}
}

// Command for bud run.
type Command struct {
	bud *bud.Command
	in  *bud.Input

	// Flags
	Flag   *framework.Flag
	Listen string // Web listener address
}

// Run the run command. That's a mouthful.
func (c *Command) Run(ctx context.Context) (err error) {
	// Find go.mod
	module, err := bud.Module(c.bud.Dir)
	if err != nil {
		return err
	}
	// Ensure we have version alignment between the CLI and the runtime
	if err := bud.EnsureVersionAlignment(ctx, module, versions.Bud); err != nil {
		return err
	}
	// Setup the logger
	log, err := bud.Log(c.in.Stderr, c.bud.Log)
	if err != nil {
		return err
	}
	// Setup the prompter
	// TODO: move this into New
	var prompter prompter.Prompter
	c.in.Stdout = io.MultiWriter(c.in.Stdout, &prompter.StdOut)
	c.in.Stderr = io.MultiWriter(c.in.Stderr, &prompter.StdErr)
	// Listening on the web listener as soon as possible
	webln := c.in.WebLn
	if webln == nil {
		// Listen and increment if the port is already in use up to 10 times
		webln, err = socket.ListenUp(c.Listen, 10)
		if err != nil {
			return err
		}
		defer webln.Close()
		log.Info("Listening on http://" + webln.Addr().String())
	}
	// Setup the default terminal prompter state
	prompter.Init(webrt.Format(webln))
	// Setup the bud listener
	budln := c.in.BudLn
	if budln == nil {
		budln, err = socket.Listen(":35729")
		if err != nil {
			return err
		}
		defer budln.Close()
		log.Debug("run: bud server is listening", "url", "http://"+budln.Addr().String())
	}
	// Load the generator filesystem
	genfs, close, err := bud.FileSystem(ctx, log, module, c.Flag, c.in)
	if err != nil {
		return err
	}
	defer close()
	// Load V8
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	// Load the file server
	servefs, err := bud.FileServer(log, module, vm, c.Flag)
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
		dir:      module.Directory(),
		builder:  gobuild.New(module),
		prompter: &prompter,
		bus:      bus,
		genfs:    genfs,
		log:      log,
		starter:  starter,
	}
	// Start the servers
	eg, ctx := errgroup.WithContext(ctx)
	// Start the internal bud server
	eg.Go(func() error { return budServer.Run(ctx) })
	// Start the internal app server
	eg.Go(func() error { return appServer.Run(ctx) })
	// Wait until either the hot or web server exits
	err = eg.Wait()
	log.Debug("run: command finished", "err", err)
	return err
}

// budServer runs the bud development server
type budServer struct {
	budln net.Listener
	bus   pubsub.Client
	fsys  fs.FS
	log   log.Interface
}

// Run the bud server
func (s *budServer) Run(ctx context.Context) error {
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	devServer := budsvr.New(s.fsys, s.bus, s.log, vm)
	err = webrt.Serve(ctx, s.budln, devServer)
	s.log.Debug("run: bud server closed", "err", err)
	return err
}

// appServer runs the generated web application
type appServer struct {
	dir      string
	builder  *gobuild.Builder
	prompter *prompter.Prompter
	bus      pubsub.Client
	genfs    *overlay.FileSystem
	log      log.Interface
	starter  *exe.Command
}

// Run the app server
func (a *appServer) Run(ctx context.Context) error {
	// Generate the app
	if err := a.genfs.Sync("bud/internal"); err != nil {
		a.bus.Publish("app:error", []byte(err.Error()))
		a.log.Debug("run: published event", "event", "app:error")
		return err
	}
	// Build the app
	if err := a.builder.Build(ctx, "bud/internal/app/main.go", "bud/app"); err != nil {
		a.bus.Publish("app:error", []byte(err.Error()))
		a.log.Debug("run: published event", "event", "app:error")
		return err
	}
	// Start the built app
	process, err := a.starter.Start(ctx, filepath.Join("bud", "app"))
	if err != nil {
		a.bus.Publish("app:error", []byte(err.Error()))
		a.log.Debug("run: published event", "event", "app:error")
		return err
	}
	// Watch for changes
	return watcher.Watch(ctx, a.dir, catchError(a.prompter, func(events []watcher.Event) error {
		a.log.Debug("run: file changes", "paths", events)
		a.prompter.Reloading(events)
		if canIncrementallyReload(events) {
			a.log.Debug("run: incrementally reloading")
			// Publish the frontend:update event
			a.bus.Publish("frontend:update", nil)
			a.log.Debug("run: published event", "event", "frontend:update")
			// Publish the app:ready event
			a.bus.Publish("app:ready", nil)
			a.log.Debug("run: published event", "event", "app:ready")
			a.prompter.SuccessReload()
			return nil
		}
		now := time.Now()
		a.log.Debug("run: restarting the process")
		if err := process.Close(); err != nil {
			return err
		}
		a.bus.Publish("backend:update", nil)
		a.log.Debug("run: published event", "event", "backend:update")
		// Generate the app
		if err := a.genfs.Sync("bud/internal"); err != nil {
			return err
		}
		// Build the app
		if err := a.builder.Build(ctx, "bud/internal/app/main.go", "bud/app"); err != nil {
			return err
		}
		// Restart the process
		p, err := process.Restart(ctx)
		if err != nil {
			a.bus.Publish("app:error", nil)
			a.log.Debug("run: published event", "event", "app:error")
			return err
		}
		a.prompter.SuccessReload()
		a.log.Debug("restarted the process", "in", time.Since(now))
		process = p
		return nil
	}))
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
