package generator

import (
	_ "embed"
	"errors"

	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/pkg/gen"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

//go:embed generator.gotext
var template string

var generator = gotemplate.MustParse("generator.gotext", template)

func New(genFS gen.FS, module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{genFS, module, parser}
}

type Generator struct {
	genFS  gen.FS
	module *gomod.Module
	parser *parser.Parser
}

func (g *Generator) GenerateFile(f gen.F, file *gen.File) error {
	// Load command state
	state, err := g.Load()
	if err != nil {
		return err
	}
	// Generate our template
	code, err := generator.Generate(state)
	if err != nil {
		return err
	}
	file.Write(code)
	return nil
}

func (g *Generator) Load() (*State, error) {
	loader := &loader{Generator: g}
	return loader.Load()
}

type loader struct {
	bail.Struct
	*Generator
}

func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	l.Bail(errors.New("not implemented yet"))
	return state, nil
}

//
