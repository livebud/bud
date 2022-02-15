package overlay

import (
	"io/fs"

	"gitlab.com/mnm/bud/package/mergefs"

	"gitlab.com/mnm/bud/package/conjure"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/pluginfs"
)

// Load the overlay filesystem
func Load(module *gomod.Module) (FS, error) {
	pluginFS, err := pluginfs.Load(module)
	if err != nil {
		return nil, err
	}
	cfs := conjure.New()
	fsys := mergefs.Merge(cfs, pluginFS)
	return fsys, nil
}

type FS = fs.FS
