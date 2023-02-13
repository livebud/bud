package parser

import "go/ast"

type TypeSpec struct {
	file *File
	node *ast.TypeSpec
}

var _ Declaration = (*TypeSpec)(nil)

func (t *TypeSpec) Name() string {
	return t.node.Name.Name
}

func (t *TypeSpec) Private() bool {
	return isPrivate(t.Name())
}

func (t *TypeSpec) File() *File {
	return t.file
}

func (t *TypeSpec) Package() *Package {
	return t.file.Package()
}

func (t *TypeSpec) Kind() Kind {
	return KindTypeSpec
}

func (t *TypeSpec) Type() Type {
	return getType(t, t.node.Type)
}
