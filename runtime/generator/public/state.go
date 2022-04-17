package public

import (
	"github.com/livebud/bud/internal/embed"
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/runtime/bud"
)

type State struct {
	Imports []*imports.Import
	Embeds  []*embed.File
	Flag    *bud.Flag
}
