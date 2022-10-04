package web

import "github.com/livebud/bud/internal/imports"

type State struct {
	Imports   []*imports.Import
	Resources []*Resource

	// TODO: remove below
	Actions     []*Action
	HasView     bool
	ShowWelcome bool
}

// Resource is a web package that will register its routes
type Resource struct {
	Import *imports.Import
	Path   string
	Camel  string
}

// TODO: remove action
type Action struct {
	Method   string
	Route    string
	CallName string
}
