package di

import (
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
)

type Type struct {
	Import string
	Type   string

	module *gomod.Module // Optional, defaults to project module
	kind   parser.Kind   // Kind of type (e.g. struct, interface, etc.)
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

func ToType(importPath, dataType string) *Type {
	return &Type{
		Import: importPath,
		Type:   dataType,
	}
}
