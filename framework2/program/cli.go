package program

import (
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
)

// ForCLI creates a program for the project CLI
func ForCLI(injector *di.Injector, module *gomod.Module) *Program {
	return &Program{
		injector: injector,
		function: &di.Function{
			Name:   "loadCLI",
			Target: module.Import("bud/.cli/program"),
			Params: []di.Dependency{
				di.ToType("gitlab.com/mnm/bud/pkg/di", "*Injector"),
				di.ToType("gitlab.com/mnm/bud/pkg/gomod", "*Module"),
				di.ToType("gitlab.com/mnm/bud/package/overlay", "*FileSystem"),
				di.ToType("gitlab.com/mnm/bud/pkg/parser", "*Parser"),
			},
			Results: []di.Dependency{
				di.ToType(module.Import("bud/.cli/command"), "*CLI"),
				&di.Error{},
			},
		},
	}
}
