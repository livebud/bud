package public

import (
	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/embed"
	"github.com/livebud/bud/internal/imports"
)

type State struct {
	Imports []*imports.Import
	Embeds  []*embed.File
	Flag    *framework.Flag
}
