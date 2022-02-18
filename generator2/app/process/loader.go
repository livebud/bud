package process

import (
	"gitlab.com/mnm/bud/internal/bail"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/parser"
)

// Load state
func Load(module *gomod.Module, parser *parser.Parser) (*State, error) {
	loader := &loader{
		imports: imports.New(),
		parser:  parser,
		module:  module,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	imports *imports.Set
	module  *gomod.Module
	parser  *parser.Parser
}

// Load the command state
func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// Add initial imports
	l.imports.AddNamed("commander", "gitlab.com/mnm/bud/pkg/commander")
	l.imports.AddNamed("gomod", "gitlab.com/mnm/bud/pkg/gomod")
	state.Imports = l.imports.List()
	// TODO: finish state
	state = &State{
		Imports: l.imports.List(),
		Process: &Process{
			Name: "app",
			Slug: "app",
		},
	}
	return state, nil
}
