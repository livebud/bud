package parser

import (
	"fmt"
	"go/ast"

	"github.com/livebud/bud/internal/gois"
)

// Definition looks a local definition up by name
// TODO: support more type definitions
func (pkg *Package) definition(name string) (decl Declaration, err error) {
	if gois.Builtin(name) {
		return builtin(name), nil
	}
	err = fmt.Errorf("parser: unable to find declaration for %q in %q", name, pkg.Name())
	var file *File
	var ts *ast.TypeSpec
	ast.Inspect(pkg.node, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.File:
			file = pkg.File(n.Name.Name)
			return true
		case *ast.TypeSpec:
			ts = n
			if n.Name.Name != name {
				return true
			}
			// Handle type spec
			if n.Assign == 0 {
				switch n.Type.(type) {
				// Continue on if we find a struct or interface
				case *ast.StructType, *ast.InterfaceType:
					return true
				default:
					// Otherwise it's a generic type spec
					decl = &TypeSpec{
						file: file,
						node: n,
					}
					err = nil
					return false
				}
			}
			// Handle type aliases
			alias := &Alias{
				file: file,
				node: n,
			}
			resolved, err2 := alias.Definition()
			if err2 != nil {
				err = err2
				return false
			}
			alias.kind = resolved.Kind()
			decl = alias
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
	return decl, err
}
