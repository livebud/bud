package generate

import (
	_ "embed"

	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/pkg/gomod"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/budfs"
	"gitlab.com/mnm/bud/pkg/gen"
)

//go:embed generate.gotext
var template string

var generator = gotemplate.MustParse("generate", template)

type Generator struct {
	BFS      budfs.FS
	Injector *di.Injector
	Module   *gomod.Module
	Embed    bool
	Hot      bool
	Minify   bool
}

type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
	Embed    bool
	Hot      bool
	Minify   bool
}

func (g *Generator) GenerateFile(_ gen.F, file *gen.File) error {
	// Don't create a generate file if custom user-generators don't exist
	if err := gen.SkipUnless(g.BFS, "bud/generator/generator.go"); err != nil {
		return err
	}
	imports := imports.New()
	imports.AddStd("os")
	// imports.AddStd("fmt")
	imports.AddNamed("console", "gitlab.com/mnm/bud/pkg/log/console")
	imports.AddNamed("mod", "gitlab.com/mnm/bud/pkg/gomod")
	imports.AddNamed("budfs", "gitlab.com/mnm/bud/pkg/budfs")
	imports.AddNamed("generate", "gitlab.com/mnm/bud/generate")
	provider, err := g.Injector.Wire(&di.Function{
		Name:   "Load",
		Target: g.Module.Import("bud/generate"),
		Params: []di.Dependency{
			&di.Type{
				Import: "gitlab.com/mnm/bud/pkg/budfs",
				Type:   "FS",
			},
			&di.Type{
				Import: "gitlab.com/mnm/bud/pkg/gomod",
				Type:   "*Module",
			},
		},
		Results: []di.Dependency{
			&di.Type{
				Import: g.Module.Import("bud/generator"),
				Type:   "Generators",
			},
			&di.Error{},
		},
	})
	if err != nil {
		return err
	}
	for _, imp := range provider.Imports {
		imports.AddNamed(imp.Name, imp.Path)
	}
	code, err := generator.Generate(&State{
		Imports:  imports.List(),
		Provider: provider,
		Embed:    g.Embed,
		Minify:   g.Minify,
		Hot:      g.Hot,
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}
