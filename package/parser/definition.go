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
			// TODO: handle type declarations (e.g. type A string)
			if n.Assign == 0 {
				ts = n
				return true
			} else if n.Name.Name != name {
				return true
			}
			alias := &Alias{
				file: file,
				ts:   n,
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
