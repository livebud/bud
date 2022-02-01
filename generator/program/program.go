package program

import (
	_ "embed"
	"fmt"

	"gitlab.com/mnm/bud/pkg/budfs"
	"gitlab.com/mnm/bud/pkg/di"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/gen"
	"gitlab.com/mnm/bud/pkg/gomod"
)

//go:embed program.gotext
var template string

var generator = gotemplate.MustParse("program.gotext", template)

type Generator struct {
	BFS      budfs.FS
	Module   *gomod.Module
	Injector *di.Injector
}

type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	if err := gen.SkipUnless(g.BFS, "bud/command/command.go"); err != nil {
		return err
	}
	// Add the imports
	imports := imports.New()
	imports.AddStd("os", "errors", "context", "runtime", "path/filepath")
	// imports.AddStd("fmt")
	imports.AddNamed("console", "gitlab.com/mnm/bud/pkg/log/console")
	imports.AddNamed("gen", "gitlab.com/mnm/bud/pkg/gen")
	imports.AddNamed("budfs", "gitlab.com/mnm/bud/pkg/budfs")
	// imports.AddNamed("plugin", "gitlab.com/mnm/bud/plugin")
	imports.AddNamed("gomod", "gitlab.com/mnm/bud/pkg/gomod")
	imports.Add(g.Module.Import("bud/command"))
	provider, err := g.Injector.Wire(&di.Function{
		Name:   "loadCLI",
		Target: g.Module.Import("bud", "program"),
		Params: []di.Dependency{
			&di.Type{Import: "gitlab.com/mnm/bud/pkg/gomod", Type: "*Module"},
			// TODO: remove gen
			&di.Type{Import: "gitlab.com/mnm/bud/pkg/gen", Type: "*FileSystem"},
			&di.Type{Import: "gitlab.com/mnm/bud/pkg/budfs", Type: "FS"},
		},
		Results: []di.Dependency{
			&di.Type{Import: g.Module.Import("bud", "command"), Type: "*CLI"},
			&di.Error{},
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
	file.Write(code)
	return nil
}
