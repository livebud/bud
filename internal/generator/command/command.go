package command

import (
	_ "embed"
	"fmt"

	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command.gotext", template)

type Generator struct {
	Module *mod.Module
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	// TODO: consider also building when only commands are present
	if err := gen.SkipUnless(f, "bud/web/web.go"); err != nil {
		return err
	}
	// Load command state
	state, err := Load(g.Module)
	if err != nil {
		return err
	}
	// Generate our template
	code, err := generator.Generate(state)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// fmt.Println(string(code))
	file.Write(code)
	return nil
}

// func (g *Generator) generateDI(f gen.F, file *gen.File) error {

// }
