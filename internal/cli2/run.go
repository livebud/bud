package cli

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
	"time"

	"github.com/livebud/bud/internal/config"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/versions"

	"github.com/livebud/bud/framework/afs"
	"github.com/livebud/bud/framework/generator"
	"github.com/livebud/bud/internal/prompter"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/virtual"
	"github.com/livebud/bud/package/watcher"

	v8 "github.com/livebud/bud/package/js/v8"

	"github.com/livebud/bud"
	"golang.org/x/sync/errgroup"
)

func (c *CLI) Run(ctx context.Context, in *bud.Run) error {
	cfg := &in.Config
	cmd := c.Command.Clone()

	// Find go.mod
	module, err := c.module()
	if err != nil {
		return err
	}
	cmd.Env = append(cmd.Env, "GOMODCACHE="+module.ModCache())

	// TODO: add config alignment check
	if err := config.EnsureVersionAlignment(ctx, module, versions.Bud); err != nil {
		return err
	}

	// Setup the logger
	log, err := c.logger()
	if err != nil {
		return err
	}

	// Listen on the web address
	webLn, err := c.listenWeb(in.WebAddress)
	if err != nil {
		return err
	}
	defer webLn.Close()
	log.Info("Listening on http://" + webLn.Addr().String())

	// Load the database
	db, err := c.openDatabase(log, module)
	if err != nil {
		return err
	}
	defer db.Close()

	// Listen on the dev address
	devLn, err := c.listenDev(in.DevAddress)
	if err != nil {
		return err
	}
	defer devLn.Close()
	cmd.Env = append(cmd.Env, "BUD_DEV_URL="+devLn.Addr().String())
	log.Debug("run: dev server is listening on http://" + devLn.Addr().String())

	vm, err := v8.Load()
	if err != nil {
		return err
	}
	defer vm.Close()

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		err := c.serveDev(ctx, cfg, devLn, log)
		return err
	})

	// Setup genfs
	fsys := virtual.Exclude(module, func(path string) bool {
		return path == "bud" || strings.HasPrefix(path, "bud/")
	})
	genfs := genfs.New(db, fsys, log)
	parser := parser.New(genfs, module)
	injector := di.New(genfs, log, module, parser)
	genfs.FileGenerator("bud/internal/generator/generator.go", generator.New(log, module, parser))
	genfs.FileGenerator("bud/cmd/afs/main.go", afs.New(cfg, injector, log, module))

	// Generate AFS
	skips := []func(name string, isDir bool) bool{}
	// Skip hidden files and directories
	skips = append(skips, func(name string, isDir bool) bool {
		base := filepath.Base(name)
		return base[0] == '_' || base[0] == '.'
	})
	// Skip files we want to carry over
	skips = append(skips, func(name string, isDir bool) bool {
		switch name {
		case "bud/bud.db", "bud/afs", "bud/app":
			return true
		default:
			return false
		}
	})
	if err := dsync.To(genfs, module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// Build the afs binary
	if err := cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+afsBinPath, afsMainPath); err != nil {
		return err
	}

	// Reset the cache
	// TODO: optimize later
	if err := db.Reset(); err != nil {
		return err
	}

	// Start the file server listener
	afsLn, err := c.listenAFS(":0")
	if err != nil {
		return err
	}
	defer afsLn.Close()
	cmd.Env = append(cmd.Env, "BUD_AFS_URL="+afsLn.Addr().String())
	log.Debug("run: afs server is listening on http://" + afsLn.Addr().String())

	// Load the *os.File for afsLn
	afsFile, err := afsLn.File()
	if err != nil {
		return err
	}
	defer afsFile.Close()

	// Inject the file under the AFS prefix
	cmd.Inject("AFS", afsFile)

	// Start afs
	afsProcess, err := cmd.Start(ctx, module.Directory("bud", "afs"))
	if err != nil {
		return err
	}
	defer afsProcess.Close()

	remoteClient, err := remotefs.Dial(ctx, afsLn.Addr().String())
	if err != nil {
		return err
	}
	defer remoteClient.Close()

	// Generate the app
	// Skip over the afs files we just generated
	skips = append(skips, func(name string, isDir bool) bool {
		return isAFSPath(name)
	})
	// Sync the app files again with the remote filesystem
	if err := dsync.To(remoteClient, module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	// Build the application binary
	if err := cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+appBinPath, appMainPath); err != nil {
		return err
	}

	// Get the file descriptor for the web listener
	webFile, err := webLn.File()
	if err != nil {
		return err
	}
	defer webFile.Close()

	// Inject that file under the WEB prefix
	cmd.Inject("WEB", webFile)

	// Start the app
	appProcess, err := cmd.Start(ctx, filepath.Join("bud", "app"))
	if err != nil {
		return err
	}
	defer appProcess.Close()

	if !in.Watch {
		return appProcess.Wait()
	}

	// Setup the prompter
	prompter := c.prompter(webLn)

	// Watch for changes
	err = watcher.Watch(ctx, module.Directory(), catchError(prompter, func(events []watcher.Event) error {
		// Trigger reloading
		prompter.Reloading(events)
		// Inform the bud filesystem of the changes
		changes := make([]string, len(events))
		for i, event := range events {
			log.Debug("run: file path changed %q", event.Path)
			changes[i] = event.Path
		}
		// Refresh the cache
		// TODO: optimize later
		if err := db.Reset(); err != nil {
			return err
		}
		// Check if we can incrementally reload
		if canIncrementallyReload(events) {
			log.Debug("run: incrementally reloading")
			// Publish the frontend:update event
			c.Bus.Publish("frontend:update", nil)
			log.Debug("run: published event %q", "frontend:update")
			// Publish the app:ready event
			c.Bus.Publish("app:ready", nil)
			log.Debug("run: published event %q", "app:ready")
			prompter.SuccessReload()
			return nil
		}
		now := time.Now()
		log.Debug("run: restarting the process")
		if err := appProcess.Close(); err != nil {
			return err
		}
		c.Bus.Publish("backend:update", nil)
		log.Debug("run: published event %q", "backend:update")
		// Generate the app
		if err := dsync.To(remoteClient, module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}
		// Build the application binary
		if err := cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+appBinPath, appMainPath); err != nil {
			return err
		}
		// Restart the process
		p, err := appProcess.Restart(ctx)
		if err != nil {
			c.Bus.Publish("app:error", nil)
			log.Debug("run: published event %q", "app:error")
			return err
		}
		prompter.SuccessReload()
		log.Debug("restarted the process in %s", time.Since(now))
		appProcess = p
		return nil
	}))
	if err != nil {
		return err
	}

	// Close the final process. This process is most likely different than the
	// deferred process.
	if err := appProcess.Close(); err != nil {
		return err
	}

	return eg.Wait()
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
