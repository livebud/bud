package pluginfs

import (
	"io/fs"
	"path"
	"strings"

	mergefs "github.com/yalue/merged_fs"
	"gitlab.com/mnm/bud/2/mod"
	"golang.org/x/sync/errgroup"
)

func Load(fsys fs.FS, module *mod.Module) (fs.FS, error) {
	plugins, err := loadPlugins(module)
	if err != nil {
		return nil, err
	}
	return merge(fsys, plugins), nil
}

// Load plugins
func loadPlugins(module *mod.Module) (plugins []*mod.Module, err error) {
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
	plugins = make([]*mod.Module, len(importPaths))
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

// Merge the filesystems into one
func merge(app fs.FS, plugins []*mod.Module) fs.FS {
	if len(plugins) == 0 {
		return app
	}
	next := app
	for _, plugin := range plugins {
		next = mergefs.NewMergedFS(next, plugin)
	}
	return next
}
