package web

import (
	"context"
	_ "embed"
	"io/fs"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

//go:embed web.gotext
var template string

var generator = gotemplate.MustParse("web", template)

type Generator struct {
	FS     fs.FS
	Module *gomod.Module
	Parser *parser.Parser
}

func (g *Generator) GenerateFile(ctx context.Context, _ overlay.F, file *overlay.File) error {
	state, err := Load(g.FS, g.Module, g.Parser)
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
