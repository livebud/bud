package web

import "gitlab.com/mnm/bud/internal/imports"

type State struct {
	Imports   []*imports.Import
	Actions   []*Action
	HasHot    bool
	HasPublic bool
	HasView   bool
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
