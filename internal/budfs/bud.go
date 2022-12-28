package budfs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

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
	"github.com/livebud/bud/package/virtual"
)

type FileSystem interface {
	Refresh(ctx context.Context, paths ...string) error
	Sync(ctx context.Context, writable virtual.FS, dirs ...string) error
	Close(ctx context.Context) error
}

func Load(cmd *shell.Command, flag *framework.Flag, module *gomod.Module, log log.Log) (FileSystem, error) {
	cache, err := dag.Load(module, module.Directory("bud/bud.db"))
	if err != nil {
		return nil, err
	}
	fsys := genfs.New(cache, module, log)
	parser := parser.New(fsys, module)
	injector := di.New(fsys, log, module, parser)
	fsys.FileGenerator("bud/internal/generator/generator.go", generator.New(log, module, parser))
	fsys.FileGenerator("bud/cmd/afs/main.go", afs.New(injector, log, module))
	return &fileSystem{cache, cmd, fsys, log, module}, nil
}

type fileSystem struct {
	cache  dag.Cache
	cmd    *shell.Command
	fsys   fs.FS
	log    log.Log
	module *gomod.Module
}

func (f *fileSystem) Refresh(ctx context.Context, paths ...string) error {
	return nil
}

func (f *fileSystem) Sync(ctx context.Context, writable virtual.FS, paths ...string) error {
	skips := []func(name string, isDir bool) bool{}
	// Skip hidden files and directories
	skips = append(skips, func(name string, isDir bool) bool {
		base := filepath.Base(name)
		return base[0] == '_' || base[0] == '.'
	})
	// Skip files that aren't within these directories
	if len(paths) > 0 {
		skips = append(skips, only(paths...))
	}
	// Reset the cache
	// TODO: optimize later
	if err := f.cache.Reset(); err != nil {
		return err
	}
	// Sync the files
	if err := dsync.To(f.fsys, f.module, "bud", dsync.WithSkip(skips...), dsync.WithLog(f.log)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// We may not need the binary if we're only generating files for afs. If
	// that's the case, we can end early
	if !needsAFSBinary(paths) {
		fmt.Println("doesnt need binary", paths)
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
	// cmd := f.cmd.Clone()
	// cmd.
	// cmd := sh.Extend()
	// Build the binary
	// Compile the args
	// args := append([]string{
	// 	"build",
	// 	"-mod=mod",
	// 	"-o=" + outPath,
	// }, flags...)
	// args = append(args, mainPath)
	// cmd := exec.CommandContext(ctx, "go", args...)
	// cmd.Env = append(b.Env,
	// 	"GOMODCACHE="+b.module.ModCache(),
	// )
	fmt.Println("Ok")
	return nil
}

func (f *fileSystem) Close(ctx context.Context) error {
	return nil
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

func only(dirs ...string) func(name string, isDir bool) bool {
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
		return !m[name] && !hasPrefix(name)
	}
}

func isWithin(parent, child string) bool {
	if parent == child {
		return true
	}
	return strings.HasPrefix(child, parent)
}
