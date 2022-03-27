package public

import (
	"gitlab.com/mnm/bud/internal/imports"
	"gitlab.com/mnm/bud/runtime/bud"
)

type State struct {
	Imports []*imports.Import
	Files   []*File
	Flag    *bud.Flag
}

type File struct {
	Path string
	Root string
	Data string
	Mode string
}
