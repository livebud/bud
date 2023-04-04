package commands

import (
	_ "embed"
	"io/fs"

	"github.com/livebud/bud/package/imports"

	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/parser"
)

func New(module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{module, parser}
}

type Generator struct {
	module *gomod.Module
	parser *parser.Parser
}

//go:embed commands.gotext
var template string

var generator = gotemplate.MustParse("commands.gotext", template)

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

// Generate code from state
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func (g *Generator) Generate(state *State) ([]byte, error) {
	return Generate(state)
}

type State struct {
	Imports []*imports.Import
}

func (g *Generator) Load(fsys fs.FS) (*State, error) {
	state := new(State)
	imset := imports.New()
	imset.AddStd("context")
	imset.AddNamed("runtime", "github.com/livebud/bud/commands/runtime")

	// Discovered commands
	imset.Add(g.module.Import("command"))
	imset.Add(g.module.Import("command/addons"))
	imset.Add(g.module.Import("command/ps"))
	imset.Add(g.module.Import("command/ps/autoscale"))
	state.Imports = imset.List()
	return state, nil
}
