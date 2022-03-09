package command

import (
	_ "embed"

	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed command.gotext
var template string

var generator = gotemplate.MustParse("command.gotext", template)

func Generate(state *State) ([]byte, error) {
	return generator.Generate(state)
}
