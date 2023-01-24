package di

import (
	"strings"

	"github.com/livebud/bud/package/parser"
)

func tryTypeAlias(alias *parser.Alias, dataType string) (*typeAlias, error) {
	if alias.Private() {
		return nil, ErrNoMatch
	}
	if strings.TrimPrefix(dataType, "*") != alias.Name() {
		return nil, ErrNoMatch
	}
	aliasType := alias.Type()
	decl, err := parser.Definition(aliasType)
	if err != nil {
		return nil, err
	}
	importPath, err := decl.Package().Import()
	if err != nil {
		return nil, err
	}
	to := &Type{
		Import: importPath,
		Type:   decl.Name(),
		kind:   decl.Kind(),
		module: decl.Package().Module(),
	}
	return &typeAlias{
		Import: importPath,
		Name:   alias.Name(),
		Type:   to,
	}, nil
}

// typeAlias is a declaration that can provide a dependency
type typeAlias struct {
	Import string
	Name   string
	Type   *Type
}

var _ Declaration = (*typeAlias)(nil)

func (a *typeAlias) ID() string {
	return `'` + a.Import + `'.` + a.Name
}

func (a *typeAlias) Dependencies() []Dependency {
	return []Dependency{a.Type}
}

func (a *typeAlias) Generate(gen Generator, inputs []*Variable) (outputs []*Variable) {
	return inputs
}
