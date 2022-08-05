package generate

import (
	"fmt"
	"io/fs"
	"strings"

	"github.com/livebud/bud/internal/bail"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/internal/scan"
	"github.com/livebud/bud/internal/valid"
	"github.com/livebud/bud/package/di"
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

// Load the command state
func (l *loader) Load(fsys fs.FS) (state *State, err error) {
	defer l.Recover2(&err, "generate")
	state = new(State)
	state.Generators = l.loadGenerators(fsys)
	if len(state.Generators) == 0 {
		return nil, fmt.Errorf("generate: error loading. %w", fs.ErrNotExist)
	}
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadGenerators(fsys fs.FS) (generators []*stateGenerator) {
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
		path = strings.TrimPrefix(path, "generator/")
		name := l.imports.Add(l.module.Import(path))
		generators = append(generators, &stateGenerator{
			ImportName: name,
			Path:       path,
			Pascal:     gotext.Pascal(path),
		})
	}
	return generators
}

func (l *loader) loadProvider() *di.Provider {
	return nil
}
