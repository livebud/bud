package parser

import (
	"go/ast"
	"unicode"
)

type Alias struct {
	file *File
	ts   *ast.TypeSpec
}

var _ Declaration = (*Alias)(nil)

func (a *Alias) File() *File {
	return a.file
}

func (a *Alias) Name() string {
	return a.ts.Name.Name
}

// Private returns true if the field is private
func (a *Alias) Private() bool {
	return unicode.IsLower(rune(a.ts.Name.Name[0]))
}

func (a *Alias) Package() *Package {
	return a.file.Package()
}

func (a *Alias) Type() Type {
	return getType(a, a.ts.Type)
}

// Definition goes to the aliases definition
func (a *Alias) Definition() (Declaration, error) {
	return Definition(a.Type())
}
