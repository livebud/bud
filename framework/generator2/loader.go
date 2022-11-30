package generator

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/package/finder"
	"github.com/livebud/bud/package/log"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
	"github.com/matthewmueller/gotext"
)

func Load(fsys budfs.FS, injector *di.Injector, log log.Log, module *gomod.Module, parser *parser.Parser) (*State, error) {
	return (&loader{
		injector: injector,
		log:      log,
		module:   module,
		parser:   parser,
		imports:  imports.New(),
	}).Load(fsys)
}

type loader struct {
	injector *di.Injector
	log      log.Log
	module   *gomod.Module
	parser   *parser.Parser
	imports  *imports.Set
	bail.Struct
}

func (l *loader) Load(bfs budfs.FS) (state *State, err error) {
	defer l.Recover2(&err, "generator")
	state = new(State)
	state.Generators = l.loadGenerators(bfs)
	state.Generators = append(state.Generators, l.loadCoreGenerators()...)
	if len(state.Generators) == 0 {
		return nil, fs.ErrNotExist
	}
	l.imports.AddStd("io/fs")
	l.imports.AddNamed("budfs", "github.com/livebud/bud/package/budfs")
	l.imports.AddNamed("gomod", "github.com/livebud/bud/package/gomod")
	l.imports.AddNamed("log", "github.com/livebud/bud/package/log")
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadGenerators(bfs budfs.FS) (generators []*CodeGenerator) {
	generatorDirs, err := finder.Find(bfs, "{generator/**.go,bud/internal/generator/*/**.go}", func(path string, isDir bool) (entries []string) {
		if !isDir && valid.GoFile(path) {
			entries = append(entries, filepath.Dir(path))
		}
		return entries
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
		// Ensure the package has a Generator and a Register command
		// matches the accepted signature
		if s := pkg.Struct("Generator"); s == nil {
			l.log.Debug("framework/generator: skipping package because there's no Generator struct")
			continue
		} else if s.Method("Register") == nil {
			l.log.Debug("framework/generator: skipping package because Generator has no Register function")
			continue
		}
		imp := &imports.Import{
			Name: l.imports.Add(importPath),
			Path: importPath,
		}
		rootlessGenerator := strings.TrimPrefix(generatorDir, "generator/")
		generators = append(generators, &CodeGenerator{
			Import: imp,
			Type:   DirGenerator,
			Path:   rootlessGenerator,
			Camel:  gotext.Camel(rootlessGenerator),
		})
	}
	return generators
}

func (l *loader) loadCoreGenerators() (generators []*CodeGenerator) {
	return []*CodeGenerator{
		&CodeGenerator{
			Import: &imports.Import{
				Name: l.imports.Add("github.com/livebud/bud/framework/app2"),
				Path: "github.com/livebud/bud/framework/app2",
			},
			Path:  "bud/cmd/app/main.go",
			Type:  FileGenerator,
			Camel: gotext.Camel("app"),
		},
		&CodeGenerator{
			Import: &imports.Import{
				Name: l.imports.Add("github.com/livebud/bud/framework/public"),
				Path: "github.com/livebud/bud/framework/public",
			},
			Path:  "bud/internal/web/public/public.go",
			Type:  FileGenerator,
			Camel: gotext.Camel("public"),
		},
	}
}

// func (l *loader) loadProvider(generators []*GeneratorState) *di.Provider {
// 	structFields := make([]*di.StructField, len(generators))
// 	for i, generator := range generators {
// 		structFields[i] = &di.StructField{
// 			Name:   generator.Pascal,
// 			Import: generator.Import.Path,
// 			Type:   "*Generator",
// 		}
// 	}
// 	provider, err := l.injector.Wire(&di.Function{
// 		Name:    "loadGenerators",
// 		Target:  l.module.Import("bud/command/generate"),
// 		Imports: l.imports,
// 		Params: []*di.Param{
// 			{Import: "github.com/livebud/bud/package/log", Type: "Log"},
// 			{Import: "github.com/livebud/bud/package/gomod", Type: "*Module"},
// 			{Import: "context", Type: "Context"},
// 		},
// 		Results: []di.Dependency{
// 			&di.Struct{
// 				Import: l.module.Import("bud/command/generate"),
// 				Type:   "*Generator",
// 				Fields: structFields,
// 			},
// 			&di.Error{},
// 		},
// 	})
// 	if err != nil {
// 		l.Bail(err)
// 	}
// 	// Add generated imports
// 	for _, imp := range provider.Imports {
// 		l.imports.AddNamed(imp.Name, imp.Path)
// 	}
// 	return provider
// }
