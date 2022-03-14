package action

import (

	// Embed templates
	"context"
	_ "embed"
	"io/fs"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/di"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

//go:embed action.gotext
var template string

var generator = gotemplate.MustParse("action.gotext", template)

type Generator struct {
	FS       fs.FS
	Injector *di.Injector
	Module   *gomod.Module
	Parser   *parser.Parser
}

func (g *Generator) GenerateFile(ctx context.Context, _ overlay.F, file *overlay.File) error {
	state, err := Load(g.FS, g.Injector, g.Module, g.Parser)
	if err != nil {
		return err
	}
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
