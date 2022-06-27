package app

import (
	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
)

type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
	Flag     *framework.Flag
}
