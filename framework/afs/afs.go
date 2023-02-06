package afs

import (
	_ "embed"
	"fmt"
	"io/fs"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/gotemplate"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/genfs"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/log"
)

//go:embed afs.gotext
var template string

var generator = gotemplate.MustParse("framework/afs/afs.gotext", template)

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

func (g *Generator) GenerateFile(fsys genfs.FS, file *genfs.File) error {
	if _, err := fs.Stat(fsys, "bud/internal/generator/generator.go"); err != nil {
		return err
	}
	state, err := Load(g.injector, g.log, g.module)
	if err != nil {
		return fmt.Errorf("framework/afs: unable to load state. %w", err)
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
	defer l.Recover2(&err, "afs")
	state = new(State)
	state.Provider = l.loadProvider()
	l.imports.AddStd("context", "errors", "os")
	l.imports.AddNamed("commander", "github.com/livebud/bud/package/commander")
	l.imports.AddNamed("afsrt", "github.com/livebud/bud/framework/afs/afsrt")
	l.imports.AddNamed("console", "github.com/livebud/bud/package/log/console")
	l.imports.AddNamed("genfs", "github.com/livebud/bud/package/genfs")
	l.imports.AddNamed("parser", "github.com/livebud/bud/package/parser")
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadProvider() *di.Provider {
	jsVM := di.ToType("github.com/livebud/bud/package/js", "VM")
	// TODO: the public generator should be able to configure this
	provider, err := l.injector.Wire(&di.Function{
		Name:    "loadGeneratorFS",
		Imports: l.imports,
		Target:  l.module.Import("bud", "cmd", "afs"),
		Params: []*di.Param{
			{Import: "github.com/livebud/bud/package/log", Type: "Log"},
			{Import: "github.com/livebud/bud/package/gomod", Type: "*Module"},
			{Import: "github.com/livebud/bud/framework", Type: "*Flag"},
			{Import: "github.com/livebud/bud/package/genfs", Type: "*FileSystem"},
			{Import: "github.com/livebud/bud/package/di", Type: "*Injector"},
			{Import: "github.com/livebud/bud/package/parser", Type: "*Parser"},
			{Import: "context", Type: "Context"},
			{Import: "github.com/livebud/bud/package/budhttp", Type: "Client"},
		},
		Aliases: di.Aliases{
			jsVM: di.ToType("github.com/livebud/bud/package/budhttp", "Client"),
		},
		Results: []di.Dependency{
			di.ToType(l.module.Import("bud/internal/generator"), "FS"),
			&di.Error{},
		},
	})
	if err != nil {
		fmt.Println(l.module.Import("bud/internal/generator"))
		// Intentionally don't wrap this error, it gets swallowed up too easily
		l.Bail(fmt.Errorf("afs: unable to wire. %s", err))
	}
	// Add generated imports
	for _, imp := range provider.Imports {
		l.imports.AddNamed(imp.Name, imp.Path)
	}
	return provider
}
