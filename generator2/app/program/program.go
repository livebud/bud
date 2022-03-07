package program

import (
	"context"
	_ "embed"
	"fmt"

	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/di"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/gomod"
)

//go:embed program.gotext
var template string

var generator = gotemplate.MustParse("program.gotext", template)

func New(injector *di.Injector, module *gomod.Module) *Generator {
	return &Generator{injector, module}
}

type Generator struct {
	injector *di.Injector
	module   *gomod.Module
}

type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
}

func (g *Generator) GenerateDir(ctx context.Context, f overlay.F, dir *overlay.Dir) error {
	dir.GenerateFile("program.go", func(ctx context.Context, f overlay.F, file *overlay.File) error {
		// if err := gen.SkipUnless(g.BFS, "bud/command/command.go"); err != nil {
		// 	return err
		// }
		// Add the imports
		imports := imports.New()
		imports.AddStd("errors", "context", "runtime", "path/filepath")
		imports.AddNamed("console", "gitlab.com/mnm/bud/pkg/log/console")
		imports.AddNamed("gomod", "gitlab.com/mnm/bud/pkg/gomod")
		imports.Add(g.module.Import("bud/.app/process"))
		provider, err := g.injector.Wire(&di.Function{
			Name:   "loadApp",
			Target: g.module.Import("bud", "program"),
			Params: []di.Dependency{
				&di.Type{Import: "gitlab.com/mnm/bud/pkg/gomod", Type: "*Module"},
			},
			Results: []di.Dependency{
				&di.Type{Import: g.module.Import("bud", ".app", "process"), Type: "*App"},
				&di.Error{},
			},
			Aliases: di.Aliases{
				di.ToType("gitlab.com/mnm/bud/pkg/js", "VM"):             di.ToType("gitlab.com/mnm/bud/pkg/js/v8client", "*Client"),
				di.ToType("gitlab.com/mnm/bud/runtime/view", "Renderer"): di.ToType("gitlab.com/mnm/bud/runtime/view", "*Server"),
			},
		})
		if err != nil {
			// Don't wrap on purpose, this error gets swallowed up easily
			return fmt.Errorf("program: unable to wire > %s", err)
		}
		for _, im := range provider.Imports {
			imports.AddNamed(im.Name, im.Path)
		}
		code, err := generator.Generate(State{
			Imports:  imports.List(),
			Provider: provider,
		})
		if err != nil {
			return err
		}
		file.Data = code
		return nil
	})
	return nil
}
