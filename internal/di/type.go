package di

import "gitlab.com/mnm/bud/go/mod"

type Type struct {
	Import string
	Type   string

	modFile *mod.File // Optional, defaults to project modfile
	name    string    // Optional, defaults to assumed name + type
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
	return finder.Find(t.modFile, t.Import, t.Type)
}
