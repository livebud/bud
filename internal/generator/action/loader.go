package action

import (
	"errors"
	"fmt"
	"io/fs"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/di"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/internal/parser"
)

func Load(injector *di.Injector, module *mod.Module, parser *parser.Parser) (*State, error) {
	loader := &loader{
		imports: imports.New(),
		module:  module,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	imports  *imports.Set
	injector *di.Injector
	module   *mod.Module
	parser   *parser.Parser
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	state.Controllers = l.loadControllers()
	// TODO: move these
	l.imports.AddStd("net/http")
	// l.imports.AddStd("fmt")
	l.imports.AddNamed("view", "gitlab.com/mnm/bud/view")
	l.imports.AddNamed("action", l.module.Import("action"))
	state.Imports = l.imports.List()
	return state, nil
}

func (l *loader) loadControllers() (controllers []*Controller) {
	des, err := fs.ReadDir(l.module, "action")
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			l.Bail(err)
		}
		return controllers
	}
	fmt.Println(len(des))
	for _, de := range des {
		fmt.Println(de.Name())
	}
	return controllers
}
