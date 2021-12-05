package di

type ID interface {
	ID() string
}

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
	Import string // Import path
	Type   string // Type of the variable
	Name   string // Name of the variable
}

type External struct {
	*Variable
	Key     string // Name to be used as a key in a struct
	Hoisted bool   // True if this external was hoisted up
}

type Declaration interface {
	ID() string
	Dependencies() []Dependency
	Generate(gen Generator, inputs []*Variable) (outputs []*Variable)
}
