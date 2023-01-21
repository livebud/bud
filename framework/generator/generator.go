package generator

import (
	_ "embed"
	"fmt"
	"io/fs"

	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
)

//go:embed generator.gotext
var template string

var generator = gotemplate.MustParse("framework/generator/generator.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New(log log.Log, module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{log, module, parser}
}

type Generator struct {
	log    log.Log
	module *gomod.Module
	parser *parser.Parser
}

// GenerateFile connects to the remotefs and mounts the remote directory.
func (g *Generator) GenerateFile(fsys genfs.FS, file *genfs.File) error {
	state, err := g.Load(fsys)
	if err != nil {
		return fmt.Errorf("framework/generator: unable to load. %w", err)
	}
	code, err := Generate(state)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}

// Load the generators
func (g *Generator) Load(fsys fs.FS) (*State, error) {
	loader := &loader{g.log, g.module, g.parser, imports.New()}
	return loader.Load(fsys)
}
