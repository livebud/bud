package cli

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/internal/shell"
	"github.com/livebud/bud/package/budhttp/budsvr"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/socket"
)

const (
	// internal/generator
	afsGeneratorDir = "bud/internal/generator"

	// cmd/afs
	afsMainPath = "bud/cmd/afs/main.go"
	afsMainDir  = "bud/cmd/afs"
	afsBinPath  = "bud/afs"

	// cmd/app
	appMainPath = "bud/internal/app/main.go"
	appBinPath  = "bud/app"
)

type Generate struct {
	Closer    *once.Closer
	Flag      *framework.Flag
	ListenAFS string
	ListenDev string
	Packages  []string

	// Optionally passed in
	bus    pubsub.Client
	db     *dag.DB
	module *gomod.Module
	log    log.Log
	afsLn  socket.Listener
	devLn  socket.Listener

	// Cached for subsequent calls
	v8        *v8.VM
	devServer *budsvr.Server
	genfs     *genfs.FileSystem
	afsFile   *os.File
	afsClient *remotefs.Client
}

func (c *CLI) Generate(ctx context.Context, in *Generate) (err error) {
	// Load the logger if not already provided
	if in.log == nil {
		in.log, err = c.loadLog()
		if err != nil {
			return err
		}
	}
	log := in.log.Field("method", "Generate").Field("package", "cli")

	// When Generate is called directly, it owns the Closer
	if in.Closer == nil {
		log.Debug("owns the closer")
		in.Closer = new(once.Closer)
		defer func() {
			log.Debug("closing")
			in.Closer.Close()
			log.Debug("closed")
		}()
	}

	// Load module if not already provided
	if in.module == nil {
		in.module, err = c.findModule()
		if err != nil {
			return err
		}
	}

	// Load and start the dev server
	if in.bus == nil {
		in.bus = c.newBus()
	}
	if in.devLn == nil {
		in.devLn, err = c.listenDev(in.ListenDev)
		if err != nil {
			return err
		}
		in.Closer.Add(in.devLn.Close)
	}
	if in.v8 == nil {
		in.v8, err = c.loadV8()
		if err != nil {
			return err
		}
		in.Closer.Add(func() error {
			in.v8.Close()
			return nil
		})
	}
	if in.devServer == nil {
		in.devServer = c.devServer(in.bus, in.devLn, in.Flag, in.log, in.v8)
		go in.devServer.Listen(ctx)
		in.Closer.Add(in.devServer.Close)
	}

	// Load and start the filesystem
	if in.db == nil {
		in.db, err = c.openDB(in.log, in.module)
		if err != nil {
			return err
		}
		in.Closer.Add(in.db.Close)
	}

	if in.genfs == nil {
		in.genfs = c.genfs(in.db, in.Flag, in.log, in.module)
	}

	// Reset the database
	// TODO: optimize and remove in the future
	if err := in.db.Reset(); err != nil {
		return err
	}

	// Skip files that aren't within these directories
	var skips []func(name string, isDir bool) bool
	if len(in.Packages) > 0 {
		if needsAFSBinary(in.Packages) {
			in.Packages = append(in.Packages, afsGeneratorDir, afsMainPath)
		}
		skips = append(skips, only(log, in.Packages...))
	}

	// Sync genfs to the filesystem
	skips = append(skips, afsSkips...)
	log.Debug("syncing afs")
	if err := dsync.To(in.genfs, in.module, "bud", dsync.WithSkip(skips...), dsync.WithLog(in.log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	log.Debug("synced afs")

	// We may not need the binary if we're only generating files for afs. If
	// that's the case, we can end early
	if !needsAFSBinary(in.Packages) {
		return nil
	}

	// Check if we have a main file to build
	if _, err := fs.Stat(in.module, afsMainPath); err != nil {
		return err
	}

	// Build the afs binary
	cmd := c.newCommand(in.module, "go", "build", "-mod=mod", "-o="+afsBinPath, afsMainPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Load and start bud/afs
	if in.afsLn == nil {
		in.afsLn, err = c.listenAFS(in.ListenAFS)
		if err != nil {
			return err
		}
		in.Closer.Add(in.afsLn.Close)
	}
	if in.afsFile == nil {
		in.afsFile, err = c.listenFileAFS(in.afsLn)
		if err != nil {
			return err
		}
		in.Closer.Add(in.afsFile.Close)
	}
	afsCommand, err := c.loadCommandAFS(in.module, in.afsFile, in.devLn)
	if err != nil {
		return err
	}
	afsProcess, err := shell.Start(ctx, afsCommand)
	if err != nil {
		return err
	}
	in.Closer.Add(afsProcess.Close)

	// Load the remote client
	if in.afsClient == nil {
		in.afsClient, err = c.dialAFS(ctx, in.afsLn)
		if err != nil {
			return err
		}
		in.Closer.Add(in.afsClient.Close)
	}

	// Sync the app files again with the remote filesystem
	skips = append(skips, appSkips...)
	log.Debug("syncing app")
	if err := dsync.To(in.afsClient, in.module, "bud", dsync.WithSkip(skips...), dsync.WithLog(in.log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	log.Debug("synced app")

	// Don't rebuild the app binary if we generate a specific path that doesn't
	// include the app binary
	if !needsAppBinary(in.Packages) {
		return nil
	}

	// Ensure we have an app file to build
	if _, err := fs.Stat(in.module, appMainPath); err != nil {
		return err
	}
	// Build the application binary
	cmd = c.newCommand(in.module, "go", "build", "-mod=mod", "-o="+appBinPath, appMainPath)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

var afsSkips = []func(name string, isDir bool) bool{
	// Skip hidden files and directories
	func(name string, isDir bool) bool {
		base := filepath.Base(name)
		return base[0] == '_' || base[0] == '.'
	},
	// Skip files we want to carry over
	func(name string, isDir bool) bool {
		switch name {
		case "bud/bud.db", "bud/afs", "bud/app":
			return true
		default:
			return false
		}
	},
}

var appSkips = []func(name string, isDir bool) bool{
	// Skip over the afs files we just generated
	func(name string, isDir bool) bool {
		return isAFSPath(name)
	},
}

func only(log log.Log, dirs ...string) func(name string, isDir bool) bool {
	m := map[string]bool{}
	for _, dir := range dirs {
		for dir != "." {
			m[dir] = true
			dir = path.Dir(dir)
		}
	}
	hasPrefix := func(name string) bool {
		for _, dir := range dirs {
			if isWithin(dir, name) {
				return true
			}
		}
		return false
	}
	return func(name string, isDir bool) bool {
		shouldSkip := !m[name] && !hasPrefix(name)
		if shouldSkip {
			log.Debug("budfs: skipping %q", name)
		}
		return shouldSkip
	}
}

func isAFSPath(fpath string) bool {
	return fpath == afsBinPath ||
		isWithin(afsGeneratorDir, fpath) ||
		isWithin(afsMainDir, fpath) ||
		fpath == "bud/cmd" // TODO: remove once we move app over to cmd/app/main.go
}

func isWithin(parent, child string) bool {
	if parent == child {
		return true
	}
	return strings.HasPrefix(child, parent)
}

func needsAFSBinary(paths []string) bool {
	if len(paths) == 0 {
		return true
	}
	for _, path := range paths {
		// Anything outside of afsGeneratorDir and afsMainDir need the binary
		if !isWithin(afsGeneratorDir, path) && !isWithin(afsMainDir, path) {
			return true
		}
	}
	return false
}

func needsAppBinary(paths []string) bool {
	if len(paths) == 0 {
		return true
	}
	for _, path := range paths {
		if path == appBinPath {
			return true
		}
	}
	return false
}
