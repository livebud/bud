package public

import (
	"gitlab.com/mnm/bud/internal/embed"
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/runtime/bud"
)

type State struct {
	Imports []*imports.Import
	Embeds  []*embed.File
	Flag    *bud.Flag
}
