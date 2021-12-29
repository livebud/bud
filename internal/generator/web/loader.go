package web

import (
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/ldflag"
	"gitlab.com/mnm/bud/vfs"
)

func Load(module *mod.Module, f gen.F) (*State, error) {
	loader := &loader{
		imports: imports.New(),
		module:  module,
	}
	return loader.Load(f)
}

type loader struct {
	bail.Struct
	imports *imports.Set
	module  *mod.Module
}

// Load the command state
func (l *loader) Load(f gen.F) (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// Add initial imports
	l.imports.AddStd("net", "net/http", "context")
	l.imports.AddNamed("hot", "gitlab.com/mnm/bud/hot")
	l.imports.AddNamed("middleware", "gitlab.com/mnm/bud/middleware")
	l.imports.AddNamed("web", "gitlab.com/mnm/bud/web")
	l.imports.AddNamed("router", l.module.Import("bud/router"))
	l.imports.AddNamed("public", l.module.Import("bud/public"))
	l.imports.AddNamed("view", l.module.Import("bud/view"))
	// Load the conditionals
	state.HasHot = ldflag.Hot()
	exists := vfs.Exists(f,
		"bud/router",
		"bud/public",
		"bud/view",
	)
	_ = exists
	// state.Command = l.loadRoot("command")
	// Load the imports
	state.Imports = l.imports.List()
	return state, nil
}
