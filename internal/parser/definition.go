package parser

import (
	"fmt"
	"go/ast"

	"gitlab.com/mnm/bud/go/is"
)

// Definition looks a local definition up by name
// TODO: support more type definitions
func (pkg *Package) definition(name string) (decl Declaration, err error) {
	if is.Builtin(name) {
		return builtin(name), nil
	}
	err = fmt.Errorf("parser: unable to find declaration for %q in %q", name, pkg.Name())
	var ts *ast.TypeSpec
	for _, file := range pkg.Files() {
		file := file
		ast.Inspect(pkg.node, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.TypeSpec:
				if n.Assign == 0 {
					ts = n
					return true
				}
				decl = &Alias{
					file: file,
					ts:   n,
				}
				err = nil
				return false
			case *ast.StructType:
				if ts == nil || ts.Name.Name != name {
					return true
				}
				decl = &Struct{
					file: file,
					ts:   ts,
					node: n,
				}
				err = nil
				return false
			case *ast.InterfaceType:
				if ts == nil || ts.Name.Name != name {
					return true
				}
				decl = &Interface{
					file: file,
					ts:   ts,
					node: n,
				}
				err = nil
				return false
			// TODO: support const and var
			case *ast.ValueSpec:
			}
			return true
		})
	}
	return decl, err
}

// builtin declaration
type builtin string

// Name is the built-in type
func (b builtin) Name() string {
	return string(b)
}

func (b builtin) Kind() Kind {
	return KindBuiltin
}

// Directory for builtin is blank
func (b builtin) Directory() string {
	return ""
}

// Package for builtin is blank
func (b builtin) Package() *Package {
	return nil
}
