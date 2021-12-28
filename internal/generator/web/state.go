package web

import "gitlab.com/mnm/bud/internal/imports"

type State struct {
	Imports   []*imports.Import
	HasHot    bool
	HasPublic bool
	HasView   bool
}
