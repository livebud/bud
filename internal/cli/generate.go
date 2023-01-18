package cli

import (
	"context"
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/log"
)

const (
	// internal/generator
	afsGeneratorDir = "bud/internal/generator"

	// cmd/afs
	afsMainPath = "bud/cmd/afs/main.go"
	afsMainDir  = "bud/cmd/afs"
	afsBinPath  = "bud/afs"

	// cmd/app
	appMainPath = "bud/cmd/app/main.go"
	appBinPath  = "bud/app"
)

type Generate struct {
	Flag      *framework.Flag
	ListenAFS string
	ListenDev string
	Packages  []string
}

func (c *CLI) Generate(ctx context.Context, in *Generate) (err error) {
	// Load the logger if not already provided
	log, err := c.loadLog()
	if err != nil {
		return err
	}
	log = log.Field("method", "Generate").Field("package", "cli")

	// Find the module if not already provided
	module, err := c.findModule()
	if err != nil {
		return err
	}

	// Ensure the runtime is the same as the bud version
	if err := versions.AlignRuntime(ctx, module, versions.Bud); err != nil {
		return err
	}

	// Load and start the dev server
	bus := c.bus()

	// Listen on the dev port
	devLn, err := c.listenDev(in.ListenDev)
	if err != nil {
		return err
	}
	log.Debug("dev listening on https://" + devLn.Addr().String())

	// Load V8
	v8, err := c.loadV8()
	if err != nil {
		return err
	}

	// Start the dev server
	ds := c.devServer(bus, devLn, in.Flag, log, v8)
	go ds.Listen(ctx)

	// Load and start the filesystem
	db, err := c.openDB(log, module)
	if err != nil {
		return err
	}

	genfs := c.genFS(db, in.Flag, log, module)

	// Reset the database
	// TODO: optimize and remove in the future
	if err := db.Reset(); err != nil {
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
	if err := dsync.To(genfs, module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
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
	if _, err := fs.Stat(module, afsMainPath); err != nil {
		return err
	}

	// Build the afs binary
	cmd := c.command(module.Directory(), "go", "build", "-mod=mod", "-o="+afsBinPath, afsMainPath)
	if err := cmd.Run(); err != nil {
		return err
	}

	// Load and start bud/afs
	afsLn, err := c.listenAFS(in.ListenAFS)
	if err != nil {
		return err
	}
	log.Debug("afs listening on https://" + afsLn.Addr().String())

	afsFile, err := c.listenFileAFS(afsLn)
	if err != nil {
		return err
	}

	// Start the bud/afs process
	if _, err := c.startAFS(ctx, in.Flag, module, afsFile, devLn); err != nil {
		return err
	}

	// Load the remote client
	afsClient, err := c.dialAFS(ctx, afsLn)
	if err != nil {
		return err
	}

	// Sync the app files again with the remote filesystem
	skips = append(skips, appSkips...)
	log.Debug("syncing app")
	if err := dsync.To(afsClient, module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
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
	if _, err := fs.Stat(module, appMainPath); err != nil {
		return err
	}
	// Build the application binary
	cmd = c.command(module.Directory(), "go", "build", "-mod=mod", "-o="+appBinPath, appMainPath)
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
		isWithin(afsMainDir, fpath)
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
