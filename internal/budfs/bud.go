package budfs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/mergefs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/errs"

	"github.com/livebud/bud/framework/afs"
	generator "github.com/livebud/bud/framework/generator"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/shell"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/virtual"
)

type FileSystem interface {
	fs.ReadDirFS
	Refresh(ctx context.Context, paths ...string) error
	Sync(ctx context.Context, writable virtual.FS, dirs ...string) error
	Close(ctx context.Context) error
}

func Load(budln net.Listener, cmd *shell.Command, flag *framework.Flag, module *gomod.Module, log log.Log) (FileSystem, error) {
	// Load the cache
	cache, err := dag.Load(log, module.Directory("bud/bud.db"))
	if err != nil {
		return nil, fmt.Errorf("bud: unable to load cache. %w", err)
	}
	fsys := virtual.Exclude(module, func(path string) bool {
		return path == "bud" || strings.HasPrefix(path, "bud/")
	})
	gfs := genfs.New(cache, fsys, log)
	parser := parser.New(gfs, module)
	injector := di.New(gfs, log, module, parser)
	generators := generator.New(log, module, parser)
	gfs.FileGenerator("bud/internal/generator/generator.go", generators)
	gfs.FileGenerator("bud/cmd/afs/main.go", afs.New(flag, injector, log, module))
	merged := gfs
	return &fileSystem{cache, cmd, flag, gfs, generators, log, merged, module, nil}, nil
}

type fileSystem struct {
	cache      dag.Cache
	cmd        *shell.Command
	flag       *framework.Flag
	genfs      fs.ReadDirFS
	generators *generator.Generator
	log        log.Log
	merged     fs.FS
	module     *gomod.Module

	remotefs *remotefs.Process // starts as nil
}

func (f *fileSystem) Refresh(ctx context.Context, paths ...string) error {
	// TODO: optimize later
	if err := f.cache.Reset(); err != nil {
		fmt.Println("Unable to reset", paths, err)
		return err
	}
	return nil
}

func (f *fileSystem) Sync(ctx context.Context, writable virtual.FS, paths ...string) error {
	log := f.log
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
	if len(paths) > 0 {
		skips = append(skips, only(log, paths...))
	}
	// Reset the cache
	// TODO: optimize later
	if err := f.cache.Reset(); err != nil {
		return err
	}
	// Sync the files
	if err := dsync.To(f.genfs, f.module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// We may not need the binary if we're only generating files for afs. If
	// that's the case, we can end early
	if !needsAFSBinary(paths) {
		return nil
	}
	// Check if we have a main file to build
	if _, err := fs.Stat(f.module, afsMainPath); err != nil {
		return err
	}
	// Build the afs binary
	if err := f.cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+afsBinPath, afsMainPath); err != nil {
		return err
	}
	remotefs, err := (*remotefs.Command)(f.cmd).Start(ctx, "bud/afs", f.flag.Flags()...)
	if err != nil {
		return err
	}
	f.remotefs = remotefs
	f.merged = mergefs.Merge(f.genfs, remotefs)
	// Skip over the afs files we just generated
	skips = append(skips, func(name string, isDir bool) bool {
		return isAFSPath(name)
	})
	// Sync the app files again with the remote filesystem
	if err := dsync.To(remotefs, f.module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// Don't rebuild the app binary if we generate a specific path that doesn't
	// include the app binary
	if !needsAppBinary(paths) {
		return nil
	}
	// Ensure we have an app file to build
	if _, err := fs.Stat(f.module, appMainPath); err != nil {
		return err
	}
	// Build the application binary
	if err := f.cmd.Run(ctx, "go", "build", "-mod=mod", "-o="+appBinPath, appMainPath); err != nil {
		return err
	}
	return nil
}

// Open implements fs.FS by proxying to genfs. This is needed for bud server.
func (f *fileSystem) Open(name string) (fs.File, error) {
	return f.merged.Open(name)
}

// ReadDir implements fs.ReadDirFS by proxying to genfs. This is needed for
// bud server.
func (f *fileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(f.merged, name)
}

func (f *fileSystem) Close(ctx context.Context) (err error) {
	errs.Join(f.cache.Close())
	// Close the remote filesystem if we opened it
	if f.remotefs != nil {
		err = errs.Join(err, f.remotefs.Close())
	}
	return err
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

func isWithin(parent, child string) bool {
	if parent == child {
		return true
	}
	return strings.HasPrefix(child, parent)
}
