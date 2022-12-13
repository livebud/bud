package afs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/gobuild"

	"github.com/livebud/bud/framework/afs"
	generator "github.com/livebud/bud/framework/generator2"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/exe"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/remotefs"
)

type Interface interface {
	// Generate handles generating and changing
	Generate(ctx context.Context, flag *framework.Flag, paths ...string) error
	Close() error
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
	appMainPath = "bud/cmd/app/main.go"
	appBinPath  = "bud/app"
)

func New(exec *exe.Template, module *gomod.Module, log log.Log) *AFS {
	// Setup the generator filesystem
	bfs := budfs.New(module, log)
	// Load the parser
	parser := parser.New(bfs, module)
	// Load the injector
	injector := di.New(bfs, log, module, parser)
	// Setup the initial afs file generators
	generators := generator.New(log, module, parser)
	afs := afs.New(exec, injector, log, module)
	bfs.FileGenerator(afsGeneratorPath, generators)
	bfs.FileGenerator(afsMainPath, afs)
	return &AFS{bfs, exec, generators, log, module, nil}
}

type AFS struct {
	bfs        *budfs.FileSystem
	exec       *exe.Template
	generators *generator.Generator
	log        log.Log
	module     *gomod.Module
	process    *remotefs.Process // Nil if there's no process running
}

var _ Interface = (*AFS)(nil)

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

func (a *AFS) Generate(ctx context.Context, flag *framework.Flag, paths ...string) error {
	fmt.Println("generating", paths, "needs_binary=", needsBinary(paths))
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
	// Sync the files
	if err := dsync.To(a.bfs, a.module, "bud", dsync.WithSkip(skips...)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// We may not need the binary if we're only generating files for afs. If
	// that's the case, we can end early
	if !needsBinary(paths) {
		return nil
	}
	// Create a go binary builder
	// TODO: Move builder into this file
	builder := gobuild.New(a.module)
	builder.Env = a.exec.Env
	builder.Stdin = a.exec.Stdin
	builder.Stdout = a.exec.Stdout
	builder.Stderr = a.exec.Stderr
	// Check if we have a main file to build
	if _, err := fs.Stat(a.module, afsMainPath); err != nil {
		return err
	}
	// Build the afs binary
	if err := builder.Build(ctx, afsMainPath, afsBinPath); err != nil {
		return err
	}
	// Launch the application file server
	// TODO: Pass the flags through
	remotefs, err := remotefs.Start(ctx, a.exec, afsBinPath)
	if err != nil {
		return err
	}
	a.process = remotefs
	// Proxy all application generators through the remote filesystem
	state, err := a.generators.Load(a.bfs)
	if err != nil {
		return err
	}
	for _, gen := range state.Generators {
		switch gen.Type {
		case generator.DirGenerator:
			a.log.Debug("generate: adding remote dir generator %q", gen.Path)
			a.bfs.GenerateDir(gen.Path, func(fsys budfs.FS, dir *budfs.Dir) error {
				subfs, err := fs.Sub(remotefs, dir.Path())
				if err != nil {
					return err
				}
				return dir.Mount(subfs)
			})
		case generator.FileGenerator:
			a.log.Debug("generate: adding remote file generator %q", gen.Path)
			a.bfs.GenerateFile(gen.Path, func(fsys budfs.FS, file *budfs.File) error {
				code, err := fs.ReadFile(remotefs, file.Target())
				if err != nil {
					return err
				}
				file.Data = code
				return nil
			})
		case generator.FileServer:
			a.log.Debug("generate: adding remote file server %q", gen.Path)
			a.bfs.ServeFile(gen.Path, func(fsys budfs.FS, file *budfs.File) error {
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
	if err := dsync.To(a.bfs, a.module, "bud", dsync.WithSkip(skips...)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// Ensure we have an app file to build
	if _, err := fs.Stat(a.module, appMainPath); err != nil {
		return err
	}
	// Build the application binary
	if err := builder.Build(ctx, appMainPath, appBinPath); err != nil {
		return err
	}
	return nil
}

func (a *AFS) Close() error {
	if a.process == nil {
		return nil
	}
	fmt.Println("closing the process")
	return a.process.Close()
}

func needsBinary(paths []string) bool {
	if len(paths) == 0 {
		return true
	}
	for _, path := range paths {
		if !isAFSPath(path) {
			return true
		}
	}
	return false
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
