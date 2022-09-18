package controller

import (

	// Embed templates

	_ "embed"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
)

//go:embed controller.gotext
var template string

var generator = gotemplate.MustParse("framework/controller/controller.gotext", template)

// Generate the controller template from state
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

// New controller generator
func New(injector *di.Injector, module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{injector, module, parser}
}

// Generator for controllers
type Generator struct {
	injector *di.Injector
	module   *gomod.Module
	parser   *parser.Parser
}

func (g *Generator) GenerateFile(fsys budfs.FS, file *budfs.File) error {
	state, err := Load(fsys, g.injector, g.module, g.parser)
	if err != nil {
		return err
	}
	code, err := Generate(state)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}
