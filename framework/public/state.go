package public

import (
	"github.com/livebud/bud/internal/embed"
	"github.com/livebud/bud/internal/imports"
)

type State struct {
	Imports []*imports.Import
	Files   []*File
}

type File struct {
	Path  string
	Route string
	Data  embed.Data
}
