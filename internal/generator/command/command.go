package command

import (
	_ "embed"

	"gitlab.com/mnm/bud/internal/di"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command", template)

type Generator struct {
	Modfile  *mod.File
	Injector *di.Injector
}

type State struct {
	Imports []*imports.Import
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	// TODO: consider also building when only commands are present
	if err := gen.Exists(f, "bud/web/web.go"); err != nil {
		return err
	}
	// 1. Load all the commands
	// TODO: fill in
	// 2. DI a virtual *Command
	// TODO: fill in
	// 3. Generate
	imports := imports.New()
	imports.AddStd("context", "errors", "os")
	// imports.AddStd("fmt")
	imports.AddNamed("web", g.Modfile.ModulePath("bud/web"))
	imports.AddNamed("commander", "gitlab.com/mnm/bud/commander")
	imports.AddNamed("console", "gitlab.com/mnm/bud/log/console")
	imports.AddNamed("plugin", "gitlab.com/mnm/bud/plugin")
	imports.AddNamed("mod", "gitlab.com/mnm/bud/go/mod")
	imports.AddNamed("gen", "gitlab.com/mnm/bud/gen")

	// TODO: pull from DI
	imports.AddNamed("js", "gitlab.com/mnm/bud/js/v8")
	imports.AddNamed("controller1", "gitlab.com/mnm/bud/example/hn/bud/controller")
	imports.AddNamed("tailwind", "gitlab.com/mnm/bud-tailwind/transform/tailwind")
	imports.AddNamed("public", "gitlab.com/mnm/bud/example/hn/bud/public")
	imports.AddNamed("transform", "gitlab.com/mnm/bud/example/hn/bud/transform")
	imports.AddNamed("view", "gitlab.com/mnm/bud/example/hn/bud/view")
	imports.AddNamed("web", "gitlab.com/mnm/bud/example/hn/bud/web")
	imports.AddNamed("controller", "gitlab.com/mnm/bud/example/hn/controller")
	imports.AddNamed("hn", "gitlab.com/mnm/bud/example/hn/internal/hn")
	imports.AddNamed("hot", "gitlab.com/mnm/bud/hot")
	imports.AddNamed("router", "gitlab.com/mnm/bud/router")
	imports.AddNamed("svelte", "gitlab.com/mnm/bud/svelte")

	code, err := generator.Generate(State{
		Imports: imports.List(),
	})
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}

type Parser struct {
}

func (p *Parser) Parse() (*State, error) {
	return &State{}, nil
}
