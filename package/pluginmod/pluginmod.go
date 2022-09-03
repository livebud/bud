package pluginmod

import (
	"errors"
	"io/fs"
	"path"
	"strings"

	"github.com/livebud/bud/package/gomod"
)

func Glob(module *gomod.Module, dir string) (plugins []*gomod.Module, err error) {
	// Get all the bud plugins that start with bud-*
	modules, err := module.FindBy(func(req *gomod.Require) bool {
		// Plugins must be directly imported, they cannot come indirectly through
		// another dependency
		if req.Indirect {
			return false
		}
		return strings.HasPrefix(path.Base(req.Mod.Path), "bud-")
	})
	if err != nil {
		return nil, err
	}
	// Add the app module to the top of the list
	modules = append([]*gomod.Module{module}, modules...)
	for _, module := range modules {
		if _, err := fs.Stat(module, dir); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		plugins = append(plugins, module)
	}
	return plugins, nil
}
