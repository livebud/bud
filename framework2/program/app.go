package program

import (
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
)

// ForApp creates a program for the app
func ForApp(injector *di.Injector, module *gomod.Module) *Program {
	return &Program{
		injector: injector,
		function: &di.Function{
			Name:   "loadApp",
			Target: module.Import("bud", "program"),
			Params: []di.Dependency{
				&di.Type{Import: "gitlab.com/mnm/bud/pkg/gomod", Type: "*Module"},
			},
			Results: []di.Dependency{
				&di.Type{Import: module.Import("bud", ".app", "process"), Type: "*App"},
				&di.Error{},
			},
			Aliases: di.Aliases{
				di.ToType("gitlab.com/mnm/bud/pkg/js", "VM"):             di.ToType("gitlab.com/mnm/bud/pkg/js/v8client", "*Client"),
				di.ToType("gitlab.com/mnm/bud/runtime/view", "Renderer"): di.ToType("gitlab.com/mnm/bud/runtime/view", "*Server"),
			},
		},
	}
}
