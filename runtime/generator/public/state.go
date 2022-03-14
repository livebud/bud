package public

import "gitlab.com/mnm/bud/internal/imports"

type State struct {
	Imports []*imports.Import
	Embed   bool
	Files   []*File
}

type File struct {
	Path string
	Root string
}
