package generator

import (
	"github.com/matthewmueller/gotext"
	"gitlab.com/mnm/bud/internal/imports"
)

type State struct {
	Imports    []*imports.Import
	Generators []*DirGenerator
}

type DirGenerator struct {
	Path   string
	Import *imports.Import
}

func (d *DirGenerator) Camel() string {
	return gotext.Camel(d.Import.Name)
}
