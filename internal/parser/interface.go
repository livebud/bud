package parser

import (
	"go/ast"
)

// Interface struct
type Interface struct {
	file *File
	ts   *ast.TypeSpec
	node *ast.InterfaceType
}

var _ Declaration = (*Interface)(nil)

// File returns the file containing the struct
func (iface *Interface) File() *File {
	return iface.file
}

// Package name
func (iface *Interface) Package() *Package {
	return iface.file.Package()
}

// Directory gets the directory
func (iface *Interface) Directory() string {
	return iface.file.Package().Directory()
}

// Name of the struct
func (iface *Interface) Name() string {
	return iface.ts.Name.Name
}

func (iface *Interface) Method(name string) *InterfaceMethod {
	if iface.node.Methods == nil {
		return nil
	}
	for _, method := range iface.node.Methods.List {
		if len(method.Names) != 1 {
			continue
		} else if method.Names[0].Name != name {
			continue
		}
		funcType, ok := method.Type.(*ast.FuncType)
		if !ok {
			continue
		}
		return &InterfaceMethod{
			iface: iface,
			name:  method.Names[0].Name,
			node:  funcType,
		}
	}
	return nil
}

func (iface *Interface) Methods() (methods []*InterfaceMethod) {
	if iface.node.Methods == nil {
		return methods
	}
	for _, method := range iface.node.Methods.List {
		if len(method.Names) != 1 {
			continue
		}
		funcType, ok := method.Type.(*ast.FuncType)
		if !ok {
			continue
		}
		methods = append(methods, &InterfaceMethod{
			iface: iface,
			name:  method.Names[0].Name,
			node:  funcType,
		})
	}
	return methods
}

type InterfaceMethod struct {
	iface *Interface
	name  string
	node  *ast.FuncType
}

func (im *InterfaceMethod) File() *File {
	return im.iface.File()
}

func (im *InterfaceMethod) Name() string {
	return im.name
}

func (im *InterfaceMethod) Params() (fields []*Param) {
	// Handle no params
	params := im.node.Params
	if params == nil {
		return
	}
	// List of fields
	for _, field := range params.List {
		for _, name := range field.Names {
			if len(field.Names) == 0 {
				fields = append(fields, &Param{
					parent: im,
					node:   field,
				})
				continue
			}
			fields = append(fields, &Param{
				parent: im,
				name:   name.Name,
				node:   field,
			})
		}
	}
	return fields
}

func (im *InterfaceMethod) Results() (fields []*Result) {
	// Handle no results
	results := im.node.Results
	if results == nil {
		return
	}
	// List of fields
	of := len(results.List)
	for i, field := range results.List {
		if len(field.Names) == 0 {
			fields = append(fields, &Result{
				parent: im,
				node:   field,
				n:      i + 1,
				of:     of,
			})
			continue
		}
		for j, name := range field.Names {
			fields = append(fields, &Result{
				parent: im,
				name:   name.Name,
				node:   field,
				n:      i + j + 1,
				of:     of,
			})
		}
	}
	return fields
}
