package program

import (
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/pkg/di"
)

type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
}
