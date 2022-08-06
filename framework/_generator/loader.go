package generator

import (
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
)

func Load(module *gomod.Module) (*State, error) {
	imports := imports.New()
	imports.AddStd("os", "net/rpc")
	imports.AddNamed("goplugin", "github.com/livebud/bud/package/goplugin")
	imports.AddNamed("console", "github.com/livebud/bud/package/log/console")
	// imports.AddNamed("tailwind", module.Import("generate/tailwind"))
	return &State{
		Imports: imports.List(),
	}, nil
}

type loader struct {
	module *gomod.Module
}

func (l *loader) Load() (*State, error) {
	return nil, nil
}
