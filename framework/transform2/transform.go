package transform

import (
	_ "embed"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
	"github.com/livebud/bud/package/remotefs"
)

//go:embed transform.gotext
var template string

var generator = gotemplate.MustParse("framework/transform/transform.gotext", template)

// Generate the transform file
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

// New transform generator
func New(flag *framework.Flag, injector *di.Injector, log log.Log, module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{flag, injector, log, module, parser, nil}
}

type Generator struct {
	flag     *framework.Flag
	injector *di.Injector
	log      log.Log
	module   *gomod.Module
	parser   *parser.Parser

	process *remotefs.Process
}

func (g *Generator) GenerateFile(fsys budfs.FS, file *budfs.File) error {
	g.log.Debug("framework/transform: generating the main.go service containing the generators")
	state, err := Load(fsys, g.injector, g.log, g.module, g.parser)
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
