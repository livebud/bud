package generator

import (
	"fmt"
	"io/fs"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/scan"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/gomod"
	"github.com/matthewmueller/gotext"
)

func Load(fsys fs.FS, module *gomod.Module) (*State, error) {
	loader := &loader{
		imports: imports.New(),
		module:  module,
	}
	return loader.Load(fsys)
}

type loader struct {
	bail.Struct
	imports *imports.Set
	module  *gomod.Module
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
	paths, err := scan.List(fsys, "generator", func(de fs.DirEntry) bool {
		if de.IsDir() {
			return valid.Dir(de.Name())
		} else {
			return valid.GoFile(de.Name())
		}
	})
	if err != nil {
		l.Bail(err)
	}
	for _, path := range paths {
		importPath := l.module.Import(path)
		imp := &imports.Import{
			Name: l.imports.Add(importPath),
			Path: importPath,
		}
		generators = append(generators, &Gen{
			Import: imp,
			Path:   path,
			Pascal: gotext.Pascal(path),
		})
	}
	return generators
}

func (l *loader) loadImports(generators []*Gen) []*imports.Import {
	if len(generators) == 0 {
		return nil
	}
	l.imports.AddNamed("overlay", "github.com/livebud/bud/package/overlay")
	return l.imports.List()
}
