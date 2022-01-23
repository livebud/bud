package di

import (
	"gitlab.com/mnm/bud/2/mod"
	"gitlab.com/mnm/bud/2/parser"
)

type Type struct {
	Import string
	Type   string

	module *mod.Module // Optional, defaults to project module
	name   string      // Optional, defaults to assumed name + type
	kind   parser.Kind // Kind of type (e.g. struct, interface, etc.)
}

var _ Dependency = (*Type)(nil)

func (t *Type) ID() string {
	return getID(t.Import, t.Type)
}

func (t *Type) ImportPath() string {
	return t.Import
}

func (t *Type) TypeName() string {
	return t.Type
}

// Find a declaration that provides this type
func (t *Type) Find(finder Finder) (Declaration, error) {
	return finder.Find(t.module, t)
}
