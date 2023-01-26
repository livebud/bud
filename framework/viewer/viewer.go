package viewer

import (
	_ "embed"
	"io/fs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
)

//go:embed viewer.gotext
var template string

var generator = gotemplate.MustParse("framework/viewer/viewer.gotext", template)

// New viewer generator
func New(flag *framework.Flag, log log.Log, module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{flag, log, module, parser}
}

type Generator struct {
	flag   *framework.Flag
	log    log.Log
	module *gomod.Module
	parser *parser.Parser
}

func (g *Generator) GenerateFile(fsys genfs.FS, file *genfs.File) error {
	state, err := g.Load(fsys)
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

type State struct {
	Imports []*imports.Import
}

func (g *Generator) Load(fsys fs.FS) (*State, error) {
	state := new(State)
	imset := imports.New()
	imset.AddNamed("router", "github.com/livebud/bud/package/router")
	state.Imports = imset.List()
	return state, nil
}
