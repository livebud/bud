package cli

import (
	"context"
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud"
	"github.com/livebud/bud/framework/afs"
	"github.com/livebud/bud/framework/generator"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/once"
	"github.com/livebud/bud/internal/sh"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/socket"
	"github.com/livebud/bud/package/virtual"
)

func (c *CLI) budfs(cfg *bud.Config, cmd *sh.Command, db *dag.DB, fileLn socket.Listener, log log.Log, module *gomod.Module) *budfs {
	fsys := virtual.Exclude(module, func(path string) bool {
		return path == "bud" || strings.HasPrefix(path, "bud/")
	})
	genfs := genfs.New(db, fsys, log)
	parser := parser.New(genfs, module)
	injector := di.New(genfs, log, module, parser)
	genfs.FileGenerator("bud/internal/generator/generator.go", generator.New(log, module, parser))
	genfs.FileGenerator("bud/cmd/afs/main.go", afs.New(cfg, injector, log, module))
	return &budfs{
		cmd,
		new(once.Closer),
		cfg,
		db,
		fileLn,
		genfs,
		log,
		module,
		nil,
	}
}

type budfs struct {
	cmd      *sh.Command
	closer   *once.Closer
	config   *bud.Config
	cache    *dag.DB
	fileLn   socket.Listener
	genfs    fs.FS
	log      log.Log
	module   *gomod.Module
	remotefs *remotefs.Process // starts as nil
}

func (b *budfs) Refresh(ctx context.Context, paths ...string) error {
	// TODO: optimize later
	if err := b.cache.Reset(); err != nil {
		return err
	}
	return nil
}

func (b *budfs) only(packages []string) func(name string, isDir bool) bool {
	m := map[string]bool{}
	for _, dir := range packages {
		for dir != "." {
			m[dir] = true
			dir = path.Dir(dir)
		}
	}
	hasPrefix := func(name string) bool {
		for _, pkg := range packages {
			if isWithin(pkg, name) {
				return true
			}
		}
		return false
	}
	return func(name string, isDir bool) bool {
		shouldSkip := !m[name] && !hasPrefix(name)
		if shouldSkip {
			b.log.Debug("budfs: skipping %q", name)
		}
		return shouldSkip
	}
}

func (b *budfs) Sync(ctx context.Context, packages ...string) error {
	log := b.log
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
	// Skip files that aren't within these directories
	if len(packages) > 0 {
		skips = append(skips, b.only(packages))
	}
	// Reset the cache
	// TODO: optimize later
	if err := b.cache.Reset(); err != nil {
		return err
	}
	// Sync the files
	if err := dsync.To(b.genfs, b.module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// We may not need the binary if we're only generating files for afs. If
	// that's the case, we can end early
	if !needsAFSBinary(packages) {
		return nil
	}
	// Check if we have a main file to build
	if _, err := fs.Stat(b.module, afsMainPath); err != nil {
		return err
	}
	// Build the afs binary
	if err := b.cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+afsBinPath, afsMainPath); err != nil {
		return err
	}
	remotefs, err := remotefs.Start(ctx, b.cmd, b.fileLn, "bud/afs", b.config.Flags()...)
	if err != nil {
		return err
	}
	b.remotefs = remotefs
	// b.merged = mergefs.Merge(b.genfs, remotefs)
	// Skip over the afs files we just generated
	skips = append(skips, func(name string, isDir bool) bool {
		return isAFSPath(name)
	})
	// Sync the app files again with the remote filesystem
	if err := dsync.To(remotefs, b.module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// Don't rebuild the app binary if we generate a specific path that doesn't
	// include the app binary
	if !needsAppBinary(packages) {
		return nil
	}
	// Ensure we have an app file to build
	if _, err := fs.Stat(b.module, appMainPath); err != nil {
		return err
	}
	// Build the application binary
	if err := b.cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+appBinPath, appMainPath); err != nil {
		return err
	}
	return nil
}

func (b *budfs) Close() error {
	return b.closer.Close()
}

func isWithin(parent, child string) bool {
	if parent == child {
		return true
	}
	return strings.HasPrefix(child, parent)
}

const (
	// internal/generator
	afsGeneratorPath = "bud/internal/generator/generator.go"
	afsGeneratorDir  = "bud/internal/generator"

	// cmd/afs
	afsMainPath = "bud/cmd/afs/main.go"
	afsMainDir  = "bud/cmd/afs"
	afsBinPath  = "bud/afs"

	// cmd/app
	appMainPath = "bud/internal/app/main.go"
	appBinPath  = "bud/app"
)

func isAFSPath(fpath string) bool {
	return fpath == afsBinPath ||
		isWithin(afsGeneratorDir, fpath) ||
		isWithin(afsMainDir, fpath) ||
		fpath == "bud/cmd" // TODO: remove once we move app over to cmd/app/main.go
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
