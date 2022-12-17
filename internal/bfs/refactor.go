package bfs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/framework"
	generator "github.com/livebud/bud/framework/generator2"
	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/errs"
	"github.com/livebud/bud/internal/gobuild"
	"github.com/livebud/bud/package/budfs"
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
	appMainPath = "bud/internal/app/main.go"
	appBinPath  = "bud/app"
)

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

func (f *FS) Generate(ctx context.Context, flag *framework.Flag, paths ...string) error {
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
	if err := dsync.To(f.fsys, f.module, "bud", dsync.WithSkip(skips...), dsync.WithLog(f.log)); err != nil {
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
	builder := gobuild.New(f.module)
	builder.Env = f.flag.Env
	builder.Stdin = f.flag.Stdin
	builder.Stdout = f.flag.Stdout
	builder.Stderr = f.flag.Stderr
	// Check if we have a main file to build
	if _, err := fs.Stat(f.module, afsMainPath); err != nil {
		return err
	}
	// Build the afs binary
	if err := builder.Build(ctx, afsMainPath, afsBinPath); err != nil {
		return err
	}
	// Launch the application file server
	// TODO: Pass the flags through
	remotefs, err := remotefs.Start(ctx, f.exec, afsBinPath)
	if err != nil {
		return err
	}
	f.process = remotefs
	// Proxy all application generators through the remote filesystem
	state, err := f.generators.Load(f.fsys)
	if err != nil {
		return err
	}
	for _, gen := range state.Generators {
		switch gen.Type {
		case generator.DirGenerator:
			f.log.Debug("generate: adding remote dir generator %q", gen.Path)
			f.fsys.GenerateDir(gen.Path, func(fsys budfs.FS, dir *budfs.Dir) error {
				subfs, err := fs.Sub(remotefs, dir.Path())
				if err != nil {
					return err
				}
				return dir.Mount(subfs)
			})
		case generator.FileGenerator:
			f.log.Debug("generate: adding remote file generator %q", gen.Path)
			f.fsys.GenerateFile(gen.Path, func(fsys budfs.FS, file *budfs.File) error {
				code, err := fs.ReadFile(remotefs, file.Target())
				if err != nil {
					return err
				}
				file.Data = code
				return nil
			})
		case generator.FileServer:
			f.log.Debug("generate: adding remote file server %q", gen.Path)
			f.fsys.ServeFile(gen.Path, func(fsys budfs.FS, file *budfs.File) error {
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
	if err := dsync.To(f.fsys, f.module, "bud", dsync.WithSkip(skips...)); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}
	// Ensure we have an app file to build
	if _, err := fs.Stat(f.module, appMainPath); err != nil {
		return err
	}
	// Build the application binary
	if err := builder.Build(ctx, appMainPath, appBinPath); err != nil {
		return err
	}
	return nil
}

func (f *FS) Change(paths ...string) error {
	if f.process != nil {
		if err := f.process.Change(paths...); err != nil {
			return err
		}
	}
	f.fsys.Change(paths...)
	return nil
}

func (f *FS) Close() (err error) {
	if f.process != nil {
		err = errs.Join(err, f.process.Close())
	}
	err = errs.Join(err, f.fsys.Close())
	return err
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
