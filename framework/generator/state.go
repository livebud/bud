package generator

import (
	"github.com/livebud/bud/internal/imports"
)

type State struct {
	Imports        []*imports.Import
	FileGenerators []*CodeGenerator
	FileServers    []*CodeGenerator
	GenerateDirs   []*CodeGenerator
	ServeFiles     []*CodeGenerator
}

type Type string

type CodeGenerator struct {
	Import *imports.Import
	Path   string // Path that triggers the generator (e.g. "bud/cmd/app/main.go")
	Camel  string
}
