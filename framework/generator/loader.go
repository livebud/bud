package generator

import (
	"io/fs"
	"path"
	"path/filepath"
	"strings"

	"github.com/livebud/bud/package/finder"
	"github.com/livebud/bud/package/log"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
	"github.com/matthewmueller/gotext"
)

type loader struct {
	log     log.Log
	module  *gomod.Module
	parser  *parser.Parser
	imports *imports.Set
	bail.Struct
}

func (l *loader) Load(fsys fs.FS) (state *State, err error) {
	defer l.Recover2(&err, "generator")
	state = new(State)
	state.Generators = l.loadGenerators(fsys)
	state.Generators = append(state.Generators, l.loadCoreGenerators()...)
	if len(state.Generators) == 0 {
		return nil, fs.ErrNotExist
	}
	l.imports.AddStd("io/fs")
	l.imports.AddNamed("genfs", "github.com/livebud/bud/package/genfs")
	l.imports.AddNamed("gomod", "github.com/livebud/bud/package/gomod")
	l.imports.AddNamed("log", "github.com/livebud/bud/package/log")
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadGenerators(fsys fs.FS) (generators []*CodeGenerator) {
	generatorDirs, err := finder.Find(fsys, "{generator/**.go,bud/internal/generator/*/**.go}", func(path string, isDir bool) (entries []string) {
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
		} else if s.Method("GenerateDir") == nil {
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
			Path:   path.Join("bud", "generator", rootlessGenerator),
			Camel:  gotext.Camel(rootlessGenerator),
		})
	}
	return generators
}

var coreGenerators = []struct {
	Import string
	Path   string
	Type   Type
}{
	{
		Import: "github.com/livebud/bud/framework/app",
		Path:   "bud/cmd/app/main.go",
		Type:   FileGenerator,
	},
	{
		Import: "github.com/livebud/bud/framework/web",
		Path:   "bud/internal/web/web.go",
		Type:   FileGenerator,
	},
	{
		Import: "github.com/livebud/bud/framework/controller",
		Path:   "bud/internal/web/controller/controller.go",
		Type:   FileGenerator,
	},
	{
		Import: "github.com/livebud/bud/framework/view",
		Path:   "bud/internal/web/view/view.go",
		Type:   FileGenerator,
	},
	{
		Import: "github.com/livebud/bud/framework/public",
		Path:   "bud/internal/web/public/public.go",
		Type:   FileGenerator,
	},
	{
		Import: "github.com/livebud/bud/framework/view/ssr",
		Path:   "bud/view/_ssr.js",
		Type:   FileGenerator,
	},
	{
		Import: "github.com/livebud/bud/framework/view/dom",
		Path:   "bud/view",
		Type:   FileServer,
	},
	{
		Import: "github.com/livebud/bud/framework/view/nodemodules",
		Path:   "bud/node_modules",
		Type:   FileServer,
	},
	// TODO: transform
}

func (l *loader) loadCoreGenerators() (generators []*CodeGenerator) {
	for _, generator := range coreGenerators {
		generators = append(generators, &CodeGenerator{
			Import: &imports.Import{
				Name: l.imports.Add(generator.Import),
				Path: generator.Import,
			},
			Path:  generator.Path,
			Type:  generator.Type,
			Camel: gotext.Camel(path.Base(generator.Import)),
		})
	}
	return generators
}
