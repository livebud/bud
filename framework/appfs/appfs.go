package appfs

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

//go:embed appfs.gotext
var template string

var generator = gotemplate.MustParse("framework/appfs/appfs.gotext", template)

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
		return fmt.Errorf("framework/appfs: unable to load state %w", err)
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
	defer l.Recover2(&err, "appfs")
	state = new(State)
	state.Provider = l.loadProvider()
	l.imports.AddStd("context", "errors", "os")
	l.imports.AddNamed("commander", "github.com/livebud/bud/package/commander")
	l.imports.AddNamed("appfsrt", "github.com/livebud/bud/framework/appfs/appfsrt")
	l.imports.AddNamed("console", "github.com/livebud/bud/package/log/console")
	l.imports.AddNamed("budfs", "github.com/livebud/bud/package/budfs")
	l.imports.AddNamed("parser", "github.com/livebud/bud/package/parser")
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadProvider() *di.Provider {
	provider, err := l.injector.Wire(&di.Function{
		Name:    "loadGeneratorFS",
		Imports: l.imports,
		Target:  l.module.Import("bud", "cmd", "appfs"),
		Params: []*di.Param{
			{Import: "github.com/livebud/bud/package/log", Type: "Log"},
			{Import: "github.com/livebud/bud/package/gomod", Type: "*Module"},
			{Import: "github.com/livebud/bud/framework", Type: "*Flag"},
			{Import: "github.com/livebud/bud/package/budfs", Type: "*FileSystem"},
			{Import: "github.com/livebud/bud/package/di", Type: "*Injector"},
			{Import: "github.com/livebud/bud/package/parser", Type: "*Parser"},
			{Import: "context", Type: "Context"},
		},
		Results: []di.Dependency{
			di.ToType(l.module.Import("bud/internal/generator"), "FS"),
			&di.Error{},
		},
	})
	if err != nil {
		// Intentionally don't wrap this error, it gets swallowed up too easily
		l.Bail(fmt.Errorf("appfs: unable to wire. %s", err))
	}
	// Add generated imports
	for _, imp := range provider.Imports {
		l.imports.AddNamed(imp.Name, imp.Path)
	}
	return provider
}
