package web

import (
	_ "embed"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
)

//go:embed web.gotext
var template string

var generator = gotemplate.MustParse("framework/web/web.gotext", template)

// Generate the web server from state
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New(module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{module, parser}
}

type Generator struct {
	module *gomod.Module
	parser *parser.Parser
}

func (g *Generator) GenerateFile(fsys *budfs.FS, file *budfs.File) error {
	state, err := Load(fsys, g.module, g.parser)
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
