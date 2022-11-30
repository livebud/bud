package generator

import (
	_ "embed"
	"fmt"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/di"
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

func New(bfs *budfs.FileSystem, flag *framework.Flag, injector *di.Injector, log log.Log, module *gomod.Module, parser *parser.Parser) *Generator {
	return &Generator{bfs, flag, injector, log, module, parser}
}

type Generator struct {
	bfs      *budfs.FileSystem
	flag     *framework.Flag
	injector *di.Injector
	log      log.Log
	module   *gomod.Module
	parser   *parser.Parser
}

// GenerateFile connects to the remotefs and mounts the remote directory.
func (g *Generator) GenerateFile(fsys budfs.FS, file *budfs.File) error {
	state, err := Load(fsys, g.injector, g.log, g.module, g.parser)
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
