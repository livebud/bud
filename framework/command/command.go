package command

import (
	_ "embed"
	"fmt"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("framework/command/command.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}

func New(flag *framework.Flag, injector *di.Injector, log log.Log, module *gomod.Module) *Generator {
	return &Generator{flag, injector, log, module}
}

type Generator struct {
	flag     *framework.Flag
	injector *di.Injector
	log      log.Log
	module   *gomod.Module
}

func (g *Generator) GenerateFile(fsys budfs.FS, file *budfs.File) error {
	state, err := Load(g.injector, g.log, g.module)
	if err != nil {
		return fmt.Errorf("framework/command: unable to load state %w", err)
	}
	code, err := Generate(state)
	if err != nil {
		return err
	}
	file.Data = code
	return nil
}

func Load(injector *di.Injector, log log.Log, module *gomod.Module) (*State, error) {
	loader := &loader{
		injector: injector,
		imports:  imports.New(),
		module:   module,
	}
	return loader.Load()
}

type loader struct {
	injector *di.Injector
	module   *gomod.Module

	imports *imports.Set
	bail.Struct
}

type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
}

func (l *loader) Load() (state *State, err error) {
	defer l.Recover2(&err, "command")
	state = new(State)
	// state.Provider = l.loadProvider()
	state.Imports = l.imports.List()
	return state, nil
}
