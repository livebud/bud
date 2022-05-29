// Package plugin provides a virtual filesystem that merges your local
// application directory structure with all required modules that start with
// bud-.
//
// When there are conflicts, bud prefers your local files over module
// files. When there are conflicts between modules, bud prioritizes modules that
// are alphanumerically higher. For example: github.com/livebud/bud-tailwind <
// github.com/livebud/bud-preflight.
//
// The plugin system doesn't recursively load or merge indirect dependencies.
package pluginfs

import (
	"path"
	"strings"

	"io/fs"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/merged"
	"golang.org/x/sync/errgroup"
)

// Load the virtual filesytem
func Load(module *gomod.Module) (fs.FS, error) {
	plugins, err := loadPlugins(module)
	if err != nil {
		return nil, err
	}
	merged := merged.Merge(append([]fs.FS{module}, plugins...)...)
	return &FS{
		merged: merged,
	}, nil
}

// Load plugins
func loadPlugins(module *gomod.Module) (plugins []fs.FS, err error) {
	modfile := module.File()
	var importPaths []string
	for _, req := range modfile.Requires() {
		// Plugins must be directly imported, they cannot come indirectly through
		// another dependency
		if req.Indirect {
			continue
		}
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
	merged fs.FS
}

func (f *FS) Open(name string) (fs.File, error) {
	return f.merged.Open(name)
}
