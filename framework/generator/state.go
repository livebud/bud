package generator

import (
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
)

type State struct {
	Imports    []*imports.Import
	Generators []*UserGenerator
	Provider   *di.Provider
}

type UserGenerator struct {
	Import *imports.Import
	Path   string
	Pascal string
}
