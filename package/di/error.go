package di

import "github.com/livebud/bud/package/parser"

// Error type
type Error struct {
}

var _ Dependency = (*Error)(nil)

func (*Error) ID() string {
	return "error"
}

func (*Error) ImportPath() string {
	return ""
}

func (*Error) TypeName() string {
	return "error"
}

func (e *Error) Find(Finder) (Declaration, error) {
	return e, nil
}

func (*Error) Dependencies() (deps []Dependency) {
	return deps
}

func (*Error) Generate(gen Generator, inputs []*Variable) (outputs []*Variable) {
	return append(outputs, &Variable{
		Import: "", // error doesn't have an import
		Name:   "err",
		Type:   "error",
		Kind:   parser.KindInterface,
	})
}
