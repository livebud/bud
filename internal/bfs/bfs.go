package bfs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/livebud/bud/internal/dsync"
	"github.com/livebud/bud/internal/exe"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/framework/afs"
	generator "github.com/livebud/bud/framework/generator2"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/virtual/vcache"
)

func Load(flag *framework.Flag, log log.Log, module *gomod.Module) (*FS, error) {
	// TODO: figure out how to re-enable caching
	fsys := budfs.New(vcache.Discard, module, log)
	parser := parser.New(fsys, module)
	injector := di.New(fsys, log, module, parser)
	// vm, err := v8.Load()
	// if err != nil {
	// 	return nil, err
	// }
	// svelteCompiler, err := svelte.Load(vm)
	// if err != nil {
	// 	return nil, err
	// }
	// transforms, err := transformrt.Load(svelte.NewTransformable(svelteCompiler))
	// if err != nil {
	// 	return nil, err
	// }
	exec := &exe.Template{
		Dir:    module.Directory(),
		Env:    flag.Env,
		Stderr: flag.Stderr,
		Stdin:  flag.Stdin,
		Stdout: flag.Stdout,
	}
	generator := generator.New(log, module, parser)
	fsys.FileGenerator("bud/internal/generator/generator.go", generator)
	fsys.FileGenerator("bud/cmd/afs/main.go", afs.New(exec, injector, log, module))
	// fsys.FileGenerator("bud/internal/app/main.go", app.New(injector, module, flag))
	// fsys.FileGenerator("bud/internal/web/web.go", web.New(module, parser))
	// fsys.FileGenerator("bud/internal/web/controller/controller.go", controller.New(injector, module, parser))
	// fsys.FileGenerator("bud/internal/web/view/view.go", view.New(module, transforms, flag))
	// fsys.FileGenerator("bud/internal/web/public/public.go", public.New(flag, module))
	// fsys.FileGenerator("bud/view/_ssr.js", ssr.New(module, transforms.SSR))
	// fsys.FileServer("bud/view", dom.New(module, transforms.DOM))
	// fsys.FileServer("bud/node_modules", dom.NodeModules(module))
	// fsys.FileGenerator("bud/internal/generator/transform/transform.go", transform.New(flag, injector, log, module, parser))
	// fsys.FileGenerator("bud/command/.generate/main.go", generator.New(fsys, flag, injector, log, module, parser))
	return &FS{exec, flag, fsys, generator, log, module, nil}, nil
}

type FS struct {
	exec       *exe.Template
	flag       *framework.Flag
	fsys       *budfs.FileSystem
	generators *generator.Generator
	log        log.Log
	module     *gomod.Module
	process    *remotefs.Process // Starts out nil
}

func (f *FS) Open(name string) (fs.File, error) {
	return f.fsys.Open(name)
}

// Skipper prevents certain files from being deleted during sync
var skipHidden = dsync.WithSkip(func(name string, isDir bool) bool {
	base := filepath.Base(name)
	return base[0] == '_' || base[0] == '.'
})

// Directories to expand
var expandDirs = [...]string{
	"bud/internal/generator",
	"bud/command/.generate",
}

func (f *FS) expand() error {
	for _, to := range expandDirs {
		if err := f.fsys.Sync(f.module, to); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}
	}
	return nil
}

// Directories to sync
var syncDirs = [...]string{
	"bud/command",
	"bud/internal",
	"bud/package",
}

// Sync delegates to either sync
func (f *FS) Sync(dirs ...string) error {
	fmt.Println("generating...")
	return f.Generate(context.Background(), f.flag, dirs...)
	// if len(dirs) == 0 {
	// 	return f.syncDefault()
	// }
	// return f.syncDirs(dirs...)
}

// syncDefault performs the sync used in `bud run`
func (f *FS) syncDefault() error {
	if err := f.expand(); err != nil {
		return err
	}
	for _, to := range syncDirs {
		if err := f.fsys.Sync(f.module, to, skipHidden); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}
	}
	return nil
}

// syncDirs syncs specific directories and is used in `bud generate`
func (f *FS) syncDirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := f.fsys.Sync(f.module, dir, skipHidden); err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return err
			}
		}
	}
	return nil
}

// REFACTOR
