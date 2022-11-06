package web

import "github.com/livebud/bud/internal/imports"

type State struct {
	Imports   []*imports.Import
	Resources []*Resource

	// TODO: remove below
	ShowWelcome bool
}

// Resource is a web package that will register its routes
type Resource struct {
	Import *imports.Import
	Path   string
	Camel  string
}
