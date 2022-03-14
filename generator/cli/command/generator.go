package command

import (
	_ "embed"

	"gitlab.com/mnm/bud/internal/gotemplate"
)

//go:embed command2.gotext
var template string

var generator = gotemplate.MustParse("command.gotext", template)

func Generate(state *State) ([]byte, error) {
	// if state.Command == nil {
	// 	return nil, fmt.Errorf("command: generator must have a root command")
	// }
	return generator.Generate(state)
}
