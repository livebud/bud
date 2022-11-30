package generator

import (
	"github.com/livebud/bud/internal/imports"
)

type State struct {
	Imports    []*imports.Import
	Generators []*CodeGenerator
}

type Type string

const (
	DirGenerator  Type = "DirGenerator"
	FileGenerator Type = "FileGenerator"
	FileServer    Type = "FileServer"
)

type CodeGenerator struct {
	Import *imports.Import
	Type   Type   // Type of generator
	Path   string // Path that triggers the generator (e.g. "bud/cmd/app/main.go")
	Camel  string
}
