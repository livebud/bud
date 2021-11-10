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
