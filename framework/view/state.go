package view

import (
	"github.com/livebud/bud/internal/embed"
	"github.com/livebud/bud/package/imports"
)

type State struct {
	Imports []*imports.Import
	Routes  []string
	Embeds  []*embed.File
}
