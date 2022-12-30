package budfs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/internal/errs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/afs"
	generator "github.com/livebud/bud/framework/generator2"
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

func Load(cmd *shell.Command, flag *framework.Flag, module *gomod.Module, log log.Log) (FileSystem, error) {
	// Load the cache
	cache, err := dag.Load(module, module.Directory("bud/bud.db"))
	if err != nil {
		return nil, err
	}
	fsys := virtual.Exclude(module, func(path string) bool {
		return path == "bud" || strings.HasPrefix(path, "bud/")
	})
	gfs := genfs.New(cache, fsys, log)
	parser := parser.New(gfs, module)
	injector := di.New(gfs, log, module, parser)
	generators := generator.New(log, module, parser)
	gfs.FileGenerator("bud/internal/generator/generator.go", generators)
	gfs.FileGenerator("bud/cmd/afs/main.go", afs.New(injector, log, module))
	return &fileSystem{cache, cmd, gfs, generators, log, module, nil}, nil
}

type fileSystem struct {
	cache      dag.Cache
	cmd        *shell.Command
	fsys       fs.ReadDirFS
	generators *generator.Generator
	log        log.Log
	module     *gomod.Module

	remotefs *remotefs.Process // starts as nil
}

func (f *fileSystem) Refresh(ctx context.Context, paths ...string) error {
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
	if err := dsync.To(f.fsys, f.module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
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
	remotefs, err := (*remotefs.Command)(f.cmd).Start(ctx, "bud/afs")
	if err != nil {
		return err
	}
	f.remotefs = remotefs
	// Proxy all application generators through the remote filesystem
	state, err := f.generators.Load(f.fsys)
	if err != nil {
		return err
	}
	// Create an new genfs that uses the remote filesystem. A fresh genfs is
	// created everytime to handle generator changes.
	gfs := genfs.New(f.cache, virtual.Tree{}, log)
	for _, gen := range state.Generators {
		switch gen.Type {
		case generator.DirGenerator:
			log.Debug("generate: adding remote dir generator %q", gen.Path)
			gfs.GenerateDir(gen.Path, func(fsys genfs.FS, dir *genfs.Dir) error {
				// subfs, err := fs.Sub(remotefs, dir.Path())
				// if err != nil {
				// 	return err
				// }
				return fmt.Errorf("TODO: implement dir.Mount(subfs)")
				// return dir.Mount(subfs)
			})
		case generator.FileGenerator:
			log.Debug("generate: adding remote file generator %q", gen.Path)
			gfs.GenerateFile(gen.Path, func(fsys genfs.FS, file *genfs.File) error {
				code, err := fs.ReadFile(remotefs, file.Target())
				if err != nil {
					return err
				}
				file.Data = code
				return nil
			})
		case generator.FileServer:
			log.Debug("generate: adding remote file server %q", gen.Path)
			gfs.ServeFile(gen.Path, func(fsys genfs.FS, file *genfs.File) error {
				code, err := fs.ReadFile(remotefs, file.Target())
				if err != nil {
					return err
				}
				file.Data = code
				return nil
			})
		default:
			return fmt.Errorf("afs: unknown type of generator %q", gen.Type)
		}
	}
	// Skip over the afs files we just generated
	skips = append(skips, func(name string, isDir bool) bool {
		return isAFSPath(name)
	})
	// Sync the app files again with the remote filesystem
	if err := dsync.To(gfs, f.module, "bud", dsync.WithSkip(skips...), dsync.WithLog(log)); err != nil {
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
	return f.fsys.Open(name)
}

// ReadDir implements fs.ReadDirFS by proxying to genfs. This is needed for
// bud server.
func (f *fileSystem) ReadDir(name string) ([]fs.DirEntry, error) {
	return f.fsys.ReadDir(name)
}

func (f *fileSystem) Close(ctx context.Context) (err error) {
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
		isWithin(afsMainDir, fpath)
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
