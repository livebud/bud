package generator

import (
	_ "embed"
	"fmt"
	"io/fs"

	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/gotemplate"
	"github.com/livebud/bud/package/imports"
	"github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/parser"
)

//go:embed generator.gotext
var template string

var generator = gotemplate.MustParse("framework/generator/generator.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

const defaultGlob = `{generator/**.go,bud/internal/generator/*/*.go}`

type Selector struct {
	Import string
	Type   string
}

var emptySelectors = map[string]Selector{}

func New(log log.Log, module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{log, module, parser, defaultGlob, emptySelectors, emptySelectors, emptySelectors}
}

type Generator struct {
	log            log.Log
	module         *gomod.Module
	parser         *parser.Parser
	Glob           string
	FileGenerators map[string]Selector
	FileServers    map[string]Selector
	DirGenerators  map[string]Selector
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
	loader := &loader{g, imports.New()}
	return loader.Load(fsys)
}
