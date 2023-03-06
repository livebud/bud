package web

import "github.com/livebud/bud/package/imports"

type State struct {
	Imports   []*imports.Import
	Resources []*Resource
}

// Resource is a web package that will register its routes
type Resource struct {
	Import *imports.Import
	Camel  string
}
