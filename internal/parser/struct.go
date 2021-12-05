package parser

import (
	"fmt"
	"go/ast"
	"unicode"
)

// Struct struct
type Struct struct {
	file *File
	ts   *ast.TypeSpec
	node *ast.StructType
}

var _ Declaration = (*Struct)(nil)

// File returns the file containing the struct
func (stct *Struct) File() *File {
	return stct.file
}

// Package name
func (stct *Struct) Package() *Package {
	return stct.file.Package()
}

// Directory gets the directory
func (stct *Struct) Directory() string {
	return stct.file.Package().Directory()
}

// Name of the struct
func (stct *Struct) Name() string {
	return stct.ts.Name.Name
}

func (stct *Struct) Kind() Kind {
	return KindStruct
}

// Private returns true if the field is private
func (stct *Struct) Private() bool {
	return unicode.IsLower(rune(stct.ts.Name.Name[0]))
}

func (stct *Struct) Field(name string) *Field {
	if stct.node.Fields == nil {
		return nil
	}
	for _, field := range stct.node.Fields.List {
		if len(field.Names) == 0 {
			ident := getIdentifier(field.Type)
			if name == ident.Name {
				return &Field{
					stct:     stct,
					name:     ident.Name,
					node:     field,
					embedded: true,
				}
			}
		}
		for _, ident := range field.Names {
			if ident.Name == name {
				return &Field{
					stct: stct,
					name: ident.Name,
					node: field,
				}
			}
		}
	}
	return nil
}

// Fields function
func (stct *Struct) Fields() (fields []*Field) {
	if stct.node.Fields == nil {
		return fields
	}
	for _, field := range stct.node.Fields.List {
		if len(field.Names) == 0 {
			id := getIdentifier(field.Type)
			fields = append(fields, &Field{
				stct:     stct,
				name:     id.Name,
				node:     field,
				embedded: true,
			})
			continue
		}
		for _, name := range field.Names {
			fields = append(fields, &Field{
				stct: stct,
				name: name.Name,
				node: field,
			})
		}
	}
	return fields
}

// PublicFields returns all public fields
func (stct *Struct) PublicFields() (fields []*Field) {
	for _, field := range stct.Fields() {
		if field.Private() {
			continue
		}
		fields = append(fields, field)
	}
	return fields
}

// getIdentifier gets an identifier from an expression
func getIdentifier(x ast.Expr) (id *ast.Ident) {
	ast.Inspect(x, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.SelectorExpr:
			id = t.Sel
			return false
		case *ast.Ident:
			id = t
			return false
		}
		return true
	})
	if id == nil {
		// Shouldn't happen, but if it does, it's a bug to fix.
		panic(fmt.Errorf("parse: unable to get identifier from expression %T", x))
	}
	return id
}

// FieldAt gets the field at i
func (stct *Struct) FieldAt(nth int) (field *Field, err error) {
	if stct.node.Fields == nil {
		return nil, fmt.Errorf("struct %q in %q has no fields", stct.Name(), stct.file.Path())
	}
	i := 0
	for _, field := range stct.node.Fields.List {
		for _, name := range field.Names {
			if nth == i {
				return &Field{
					stct: stct,
					name: name.Name,
					node: field,
				}, nil
			}
			i++
		}
	}
	return nil, fmt.Errorf("struct %q in %q has no field at %d", stct.Name(), stct.file.Path(), nth)
}

// Field is a regular struct field
type Field struct {
	stct     *Struct
	name     string
	node     *ast.Field
	embedded bool
}

var _ Fielder = (*Field)(nil)

// File that contains this field
func (f *Field) File() *File {
	return f.stct.File()
}

// Name of the field
func (f *Field) Name() string {
	return f.name
}

// Embedded is true if the field is embedded
// TODO: all fields should have a name, we should have and additional boolean to
// determine if it's an embedded field or not.
// func (f *Field) Embedded() bool {
// 	return f.embedded
// }

// Private returns true if the field is private
func (f *Field) Private() bool {
	return unicode.IsLower(rune(f.name[0]))
}

// Type of the field
func (f *Field) Type() Type {
	return getType(f, f.node.Type)
}

// Definition gets the definition of the type
func (f *Field) Definition() (Declaration, error) {
	return Definition(f.Type())
}
