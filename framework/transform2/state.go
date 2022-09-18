package transform

import (
	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/di"
)

type State struct {
	Imports      []*imports.Import
	Transformers []*Transformer
	Provider     *di.Provider
}

type Transformer struct {
	Import     *imports.Import
	Path       string
	Pascal     string
	Transforms []*Transform
}

type Transform struct {
	From string // e.g. .svg
	To   string // e.g. .svelte
	Name string // e.g. SvgToSvelte
}
