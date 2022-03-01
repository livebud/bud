package pluginfs

import (
	"context"
	"path"
	"strings"

	"gitlab.com/mnm/bud/internal/fscache"
	"gitlab.com/mnm/bud/package/fs"
	"gitlab.com/mnm/bud/package/merged"
	"gitlab.com/mnm/bud/pkg/gomod"
	"golang.org/x/sync/errgroup"
)

type Option = func(o *option)

type option struct {
	fsCache *fscache.Cache // can be nil
}

// WithFSCache uses a custom mod cache instead of the default
func WithFSCache(cache *fscache.Cache) func(o *option) {
	return func(opt *option) {
		opt.fsCache = cache
	}
}

func Load(module *gomod.Module, options ...Option) (fs.FS, error) {
	opt := &option{
		fsCache: nil,
	}
	plugins, err := loadPlugins(module)
	if err != nil {
		return nil, err
	}
	merged := merged.Merge(append([]fs.FS{module}, plugins...)...)
	return &FS{
		opt:    opt,
		merged: merged,
	}, nil
}

// Load plugins
func loadPlugins(module *gomod.Module) (plugins []fs.FS, err error) {
	modfile := module.File()
	var importPaths []string
	for _, req := range modfile.Requires() {
		// The last path in the module path needs to start with "bud-"
		if !strings.HasPrefix(path.Base(req.Mod.Path), "bud-") {
			continue
		}
		importPaths = append(importPaths, req.Mod.Path)
	}
	// Concurrently resolve directories
	plugins = make([]fs.FS, len(importPaths))
	eg := new(errgroup.Group)
	for i, importPath := range importPaths {
		i, importPath := i, importPath
		eg.Go(func() error {
			module, err := module.Find(importPath)
			if err != nil {
				return err
			}
			plugins[i] = module
			return nil
		})
	}
	// Wait for modules to finish resolving
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return plugins, nil
}

type FS struct {
	opt    *option
	merged *merged.FS
}

func (f *FS) Open(name string) (fs.File, error) {
	if f.opt.fsCache == nil {
		return f.merged.Open(name)
	}
	return f.cachedOpen(f.opt.fsCache, name)
}

func (f *FS) OpenContext(ctx context.Context, name string) (fs.File, error) {
	// TODO: support caching
	return f.merged.OpenContext(ctx, name)
}

func (f *FS) cachedOpen(fmap *fscache.Cache, name string) (fs.File, error) {
	if fmap.Has(name) {
		return fmap.Open(name)
	}
	file, err := f.merged.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	vfile, err := fscache.From(file)
	if err != nil {
		return nil, err
	}
	fmap.Set(name, vfile)
	return fmap.Open(name)
}
