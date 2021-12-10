package command

import (
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/bail"
)

// Load state
func Load(modFile *mod.File) (*State, error) {
	loader := &loader{
		modFile: modFile,
	}
	return loader.Load()
}

type loader struct {
	bail.Struct
	modFile *mod.File
	// 	imports *imports.Set
}

func (l *loader) Load() (state *State, err error) {
	defer l.Recover(&err)
	state = new(State)
	// state.Files = l.loadFiles(node.Prerequisites())
	// state.Imports = l.loadImports()
	return state, nil
}
