package view

import (
	"github.com/livebud/bud/internal/embed"
	"github.com/livebud/bud/internal/imports"
)

type State struct {
	Imports []*imports.Import
	Routes  []string
	Embeds  []*embed.File
}
