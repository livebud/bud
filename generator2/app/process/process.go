package process

import (
	_ "embed"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

//go:embed process.gotext
var template string

var generator = gotemplate.MustParse("process.gotext", template)

func New(module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{module, parser}
}

type Generator struct {
	module *gomod.Module
	parser *parser.Parser
}

func (g *Generator) GenerateDir(f overlay.F, dir *overlay.Dir) error {
	dir.GenerateFile("process.go", func(f overlay.F, file *overlay.File) error {
		// Load process state
		state, err := Load(g.module, g.parser)
		if err != nil {
			return err
		}
		// Generate our template
		file.Data, err = generator.Generate(state)
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}
