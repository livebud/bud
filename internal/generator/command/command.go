package command

import (
	_ "embed"
	"fmt"

	"gitlab.com/mnm/bud/internal/di"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command", template)

type Generator struct {
	Module   *mod.Module
	Injector *di.Injector
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) (err error) {
	// TODO: consider also building when only commands are present
	if err := gen.Exists(f, "bud/web/web.go"); err != nil {
		return err
	}
	// Load command state
	state, err := Load(g.Injector, g.Module)
	if err != nil {
		return err
	}
	// Generate our template
	code, err := generator.Generate(state)
	if err != nil {
		fmt.Println(err)
		return err
	}
	file.Write(code)
	return nil
}
