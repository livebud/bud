package generator

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/scan"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
	"github.com/matthewmueller/gotext"
)

func Load(fsys budfs.FS, injector *di.Injector, module *gomod.Module, parser *parser.Parser) (*State, error) {
	return (&loader{
		injector: injector,
		module:   module,
		parser:   parser,
		imports:  imports.New(),
	}).Load(fsys)
}

type loader struct {
	injector *di.Injector
	module   *gomod.Module
	parser   *parser.Parser
	imports  *imports.Set
	bail.Struct
}

func (l *loader) Load(bfs budfs.FS) (state *State, err error) {
	defer l.Recover2(&err, "generator")
	state = new(State)
	// TODO: replace with a watch function, this is currently wasteful if we're
	// not going to use the results
	if files, err := fs.Glob(bfs, "generator/**.go"); err != nil {
		return nil, err
	} else if len(files) == 0 {
		return nil, fs.ErrNotExist
	}
	state.Generators = l.loadGenerators(bfs)
	if len(state.Generators) == 0 {
		return nil, fs.ErrNotExist
	}
	state.Provider = l.loadProvider(state.Generators)
	l.imports.AddStd("context", "os", "errors")
	l.imports.AddNamed("gomod", "github.com/livebud/bud/package/gomod")
	l.imports.AddNamed("log", "github.com/livebud/bud/package/log")
	l.imports.AddNamed("filter", "github.com/livebud/bud/package/log/filter")
	l.imports.AddNamed("console", "github.com/livebud/bud/package/log/console")
	l.imports.AddNamed("commander", "github.com/livebud/bud/package/commander")
	l.imports.AddNamed("remotefs", "github.com/livebud/bud/package/remotefs")
	l.imports.AddNamed("budfs", "github.com/livebud/bud/package/budfs")
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadGenerators(bfs budfs.FS) (generators []*UserGenerator) {
	generatorDirs, err := scan.List(bfs, "generator", func(de fs.DirEntry) bool {
		if de.IsDir() {
			return valid.Dir(de.Name())
		} else {
			return valid.GoFile(de.Name())
		}
	})
	if err != nil {
		l.Bail(err)
	}
	for _, generatorDir := range generatorDirs {
		importPath := l.module.Import(generatorDir)
		pkg, err := l.parser.Parse(generatorDir)
		if err != nil {
			l.Bail(err)
		}
		// Ensure the package has a generator
		// TODO: ensure the package has a GenerateDir function that
		// matches the accepted signature
		if s := pkg.Struct("Generator"); s == nil {
			l.Bail(fmt.Errorf("no Generator struct in %q", importPath))
		} else if s.Method("GenerateDir") == nil {
			l.Bail(fmt.Errorf("no (*Generator).GenerateDir(...) method in %q", importPath))
		}
		imp := &imports.Import{
			Name: l.imports.Add(importPath),
			Path: importPath,
		}
		rootlessGenerator := strings.TrimPrefix(generatorDir, "generator/")
		generators = append(generators, &UserGenerator{
			Import: imp,
			Path:   rootlessGenerator,
			Pascal: gotext.Pascal(rootlessGenerator),
		})
	}
	return generators
}

func (l *loader) loadProvider(generators []*UserGenerator) *di.Provider {
	structFields := make([]*di.StructField, len(generators))
	for i, generator := range generators {
		structFields[i] = &di.StructField{
			Name:   generator.Pascal,
			Import: generator.Import.Path,
			Type:   "*Generator",
		}
	}
	provider, err := l.injector.Wire(&di.Function{
		Name:    "loadGenerators",
		Target:  l.module.Import("bud/internal/generator"),
		Imports: l.imports,
		Params: []*di.Param{
			{Import: "github.com/livebud/bud/package/log", Type: "Interface"},
			{Import: "github.com/livebud/bud/package/gomod", Type: "*Module"},
			{Import: "context", Type: "Context"},
		},
		Results: []di.Dependency{
			&di.Struct{
				Import: l.module.Import("bud/internal/generator"),
				Type:   "*Generator",
				Fields: structFields,
			},
			&di.Error{},
		},
	})
	if err != nil {
		l.Bail(err)
	}
	// Add generated imports
	for _, imp := range provider.Imports {
		l.imports.AddNamed(imp.Name, imp.Path)
	}
	return provider
}
