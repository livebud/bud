package generate

import "github.com/livebud/bud/internal/imports"

type State struct {
	Imports    []*imports.Import
	Generators []*stateGenerator
}

type stateGenerator struct {
	ImportName string
	Path       string
	Pascal     string
}
