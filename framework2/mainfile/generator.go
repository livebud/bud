package mainfile

import (
	_ "embed"

	"gitlab.com/mnm/bud/internal/gotemplate"
	"gitlab.com/mnm/bud/internal/imports"
)

//go:embed main.gotext
var template string

var generator = gotemplate.MustParse("main.gotext", template)

// State for the generator
type State struct {
	Imports []*imports.Import
}

// Generate a main file
func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}
