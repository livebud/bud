package app

import (
	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/package/di"
	"github.com/livebud/bud/package/imports"
)

type State struct {
	Imports  []*imports.Import
	Provider *di.Provider
	Flag     *framework.Flag
}
