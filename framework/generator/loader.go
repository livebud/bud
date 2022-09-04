package generator

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/scan"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
	"github.com/matthewmueller/gotext"
)

func Load(fsys fs.FS, module *gomod.Module, parser *parser.Parser) (*State, error) {
	loader := &loader{
		imports: imports.New(),
		module:  module,
		parser:  parser,
	}
	return loader.Load(fsys)
}

type loader struct {
	bail.Struct
	imports *imports.Set
	module  *gomod.Module
	parser  *parser.Parser
}

func (l *loader) Load(fsys fs.FS) (state *State, err error) {
	defer l.Recover2(&err, "generator")
	state = new(State)
	state.Generators = l.loadGenerators(fsys)
	if len(state.Generators) == 0 {
		return nil, fmt.Errorf("generator: no custom generators: %w", fs.ErrNotExist)
	}
	state.Imports = l.loadImports(state.Generators)
	return state, nil
}

func (l *loader) loadGenerators(fsys fs.FS) (generators []*Gen) {
	relDirs, err := scan.List(fsys, "generator", func(de fs.DirEntry) bool {
		if de.IsDir() {
			return valid.Dir(de.Name())
		} else {
			return valid.GoFile(de.Name())
		}
	})
	if err != nil {
		l.Bail(err)
	}
	for _, relDir := range relDirs {
		importPath := l.module.Import(relDir)
		pkg, err := l.parser.Parse(relDir)
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
		generators = append(generators, &Gen{
			Import: imp,
			Path:   relDir,
			Pascal: gotext.Pascal(strings.TrimPrefix(relDir, "generator/")),
		})
	}
	return generators
}

func (l *loader) loadImports(generators []*Gen) []*imports.Import {
	if len(generators) == 0 {
		return nil
	}
	l.imports.AddNamed("budfs", "github.com/livebud/bud/package/budfs")
	return l.imports.List()
}
