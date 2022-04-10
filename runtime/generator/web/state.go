package web

import "gitlab.com/mnm/bud/internal/imports"

type State struct {
	Imports []*imports.Import

	Actions   []*Action
	HasPublic bool
	HasView   bool

	// Show the welcome page
	ShowWelcome bool
}

type Def struct {
	Type string
	Name string
}

type Action struct {
	Method   string
	Route    string
	CallName string
}
