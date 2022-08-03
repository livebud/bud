package generate

import (
	"io/fs"

	"github.com/livebud/bud/internal/imports"
)

func Load(fsys fs.FS) (*State, error) {
	imports := imports.New()
	imports.AddStd("os", "net/rpc")
	imports.AddNamed("goplugin", "github.com/livebud/bud/package/goplugin")
	imports.AddNamed("console", "github.com/livebud/bud/package/log/console")
	return &State{
		Imports: imports.List(),
	}, nil
}
