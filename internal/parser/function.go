package parser

import (
	"go/ast"
	"strings"
	"unicode"
)

// Function struct
type Function struct {
	pkg  *Package
	file *File
	node *ast.FuncDecl
}

// Package returns the package containing this function
func (fn *Function) Package() *Package {
	return fn.pkg
}

// File returns file containing this function
func (fn *Function) File() *File {
	return fn.file
}

// Private checks if the function is private or public
func (fn *Function) Private() bool {
	return unicode.IsLower(rune(fn.node.Name.Name[0]))
}

// Name of the function
func (fn *Function) Name() string {
	return fn.node.Name.Name
}

// Receiver returns the receiver field, if any
func (fn *Function) Receiver() *Receiver {
	if fn.node.Recv == nil {
		return nil
	}
	if len(fn.node.Recv.List) == 0 {
		return nil
	}
	field := fn.node.Recv.List[0]
	names := field.Names
	if len(names) == 0 {
		return &Receiver{
			fn:   fn,
			node: field,
		}
	}
	return &Receiver{
		fn:   fn,
		name: names[0].Name,
		node: field,
	}
}

// Params returns parameters
func (fn *Function) Params() (fields []*Param) {
	// Handle no params
	params := fn.node.Type.Params
	if params == nil {
		return
	}
	// List of fields
	for _, field := range params.List {
		for _, name := range field.Names {
			if len(field.Names) == 0 {
				fields = append(fields, &Param{
					fn:   fn,
					node: field,
				})
				continue
			}
			fields = append(fields, &Param{
				fn:   fn,
				name: name.Name,
				node: field,
			})
		}
	}
	return fields
}

// Results returns parameters
func (fn *Function) Results() (fields []*Result) {
	// Handle no results
	results := fn.node.Type.Results
	if results == nil {
		return
	}
	// List of fields
	of := len(results.List)
	for i, field := range results.List {
		if len(field.Names) == 0 {
			fields = append(fields, &Result{
				fn:   fn,
				node: field,
				n:    i + 1,
				of:   of,
			})
			continue
		}
		for j, name := range field.Names {
			fields = append(fields, &Result{
				fn:   fn,
				name: name.Name,
				node: field,
				n:    i + j + 1,
				of:   of,
			})
		}
	}
	return fields
}

// Signature returns the function signature
func (fn *Function) Signature() string {
	out := new(strings.Builder)
	out.WriteString("func")
	recv := fn.Receiver()
	if recv != nil {
		out.WriteString(" (")
		out.WriteString(recv.String())
		out.WriteString(")")
	}
	out.WriteString(" ")
	out.WriteString(fn.Name())
	out.WriteString("(")
	for i, param := range fn.Params() {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(param.String())
	}
	out.WriteString(")")
	results := fn.Results()
	if len(results) > 0 {
		out.WriteString(" (")
		for i, result := range fn.Results() {
			if i > 0 {
				out.WriteString(", ")
			}
			out.WriteString(result.String())
		}
		out.WriteString(")")
	}
	return out.String()
}

// Receiver is a function input
type Receiver struct {
	fn   *Function
	name string
	node *ast.Field
}

var _ Fielder = (*Receiver)(nil)

// File that contains this field
func (f *Receiver) File() *File {
	return f.fn.File()
}

// Name of the field
func (f *Receiver) Name() string {
	return f.name
}

// Private returns true if the field is private
func (f *Receiver) Private() bool {
	return unicode.IsLower(rune(f.name[0]))
}

// Type of the field
func (f *Receiver) Type() Type {
	return getType(f, f.node.Type)
}

// Definition gets the definition of the type
func (f *Receiver) Definition() (Declaration, error) {
	return Definition(f.Type())
}

func (f *Receiver) String() string {
	return fieldString(f)
}

// Param is a function input
type Param struct {
	fn   *Function
	name string
	node *ast.Field
}

var _ Fielder = (*Param)(nil)

// File that contains this field
func (f *Param) File() *File {
	return f.fn.File()
}

// Name of the field
func (f *Param) Name() string {
	return f.name
}

// Private returns true if the field is private
func (f *Param) Private() bool {
	return unicode.IsLower(rune(f.name[0]))
}

// Type of the field
func (f *Param) Type() Type {
	return getType(f, f.node.Type)
}

// Definition gets the definition of the type
func (f *Param) Definition() (Declaration, error) {
	return Definition(f.Type())
}

func (f *Param) String() string {
	return fieldString(f)
}

// Result is a function output
type Result struct {
	fn   *Function
	name string
	node *ast.Field

	// Used to implement First() and Last()
	n, of int
}

var _ Fielder = (*Result)(nil)

// First result
func (f *Result) First() bool {
	return f.n == 1
}

// Last result
func (f *Result) Last() bool {
	return f.n == f.of
}

// File that contains this field
func (f *Result) File() *File {
	return f.fn.File()
}

// Name of the field
func (f *Result) Name() string {
	return f.name
}

// Named returns true if the result has a name
func (f *Result) Named() bool {
	return f.name != ""
}

// Private returns true if the field is private
func (f *Result) Private() bool {
	return unicode.IsLower(rune(f.name[0]))
}

// Type of the field
func (f *Result) Type() Type {
	return getType(f, f.node.Type)
}

// IsError returns true if the field is an error
// TODO: eventually check if field implements error
func (f *Result) IsError() bool {
	id, ok := f.node.Type.(*ast.Ident)
	if !ok {
		return false
	}
	return id.Name == "error"
}

func (f *Result) String() string {
	return fieldString(f)
}
