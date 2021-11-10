package command

import (
	_ "embed"
	"fmt"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command", template)

type Generator struct {
	Modfile mod.File
}

type State struct {
	Imports []*imports.Import
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	// TODO: consider also building when only commands are present
	if err := gen.Exists(f, "bud/web/web.go"); err != nil {
		fmt.Println("error?", err)
		return err
	}
	imports := imports.New()
	imports.AddStd("context")
	imports.AddNamed("web", g.Modfile.ModulePath("bud/web"))
	imports.AddNamed("commander", "gitlab.com/mnm/bud/commander")
	imports.AddNamed("console", "gitlab.com/mnm/bud/log/console")
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
