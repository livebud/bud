package di

import "gitlab.com/mnm/bud/package/parser"

type Aliases map[Dependency]Dependency

type Dependency interface {
	ID() string
	ImportPath() string
	TypeName() string
	Find(Finder) (Declaration, error)
}

func getID(importPath, typeName string) string {
	return `"` + importPath + `".` + typeName
}

type Generator interface {
	WriteString(code string) (n int, err error)
	Identifier(importPath, name string) string
	Variable(importPath, name string) string
	MarkError(hasError bool)
}

type Variable struct {
	Import string      // Import path
	Type   string      // Type of the variable
	Name   string      // Name of the variable
	Kind   parser.Kind // Kind of type (struct, interface, etc.)
}

func (v *Variable) ID() string {
	return getID(v.Import, v.Type)
}

type External struct {
	Variable *Variable
	Key      string // Name to be used as a key in a struct
	Hoisted  bool   // True if this external was hoisted up
	FullType string // Type name including package name
}

type Declaration interface {
	ID() string
	Dependencies() []Dependency
	Generate(gen Generator, inputs []*Variable) (outputs []*Variable)
}

// Check if the field or variable is an interface
func isInterface(k parser.Kind) bool {
	return k == parser.KindInterface
}
