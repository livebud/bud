package parser

import (
	"fmt"
	"go/ast"
	"strconv"

	"gitlab.com/mnm/bud/internal/imports"
)

// File struct
type File struct {
	pkg  *Package
	path string
	node *ast.File
}

// Package returns the package type
func (f *File) Package() *Package {
	return f.pkg
}

// Path returns the file path
func (f *File) Path() string {
	return f.path
}

// Imports fn
func (f *File) Imports() (map[string]string, error) {
	out := map[string]string{}
	for _, imp := range f.node.Imports {
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			return nil, err
		}
		if imp.Name == nil {
			name := imports.AssumedName(path)
			out[name] = path
			continue
		}
		out[imp.Name.Name] = path
	}
	return out, nil
}

// Imports fn
// TODO: swap with above
// func (f *File) Imports() []*imports.Import {
// 	imports := imports.New()
// 	for _, imp := range f.node.Imports {
// 		path := strings.Trim(imp.Path.Value, `"`)
// 		if imp.Name == nil {
// 			imports.Add(path)
// 			continue
// 		}
// 		imports.AddNamed(imp.Name.Name, path)
// 	}
// 	return imports.List()
// }

// Import finds the import path of the package containing this file
func (f *File) Import() (path string, err error) {
	return f.pkg.Import()
}

// ImportPath finds an import path by import name
func (f *File) ImportPath(name string) (path string, err error) {
	for _, imp := range f.node.Imports {
		if imp.Name == nil {
			path, err = strconv.Unquote(imp.Path.Value)
			if err != nil {
				return "", err
			}
			if name == imports.AssumedName(path) {
				return path, nil
			}
			continue
		}
		if name == imp.Name.Name {
			return strconv.Unquote(imp.Path.Value)
		}
	}
	return "", fmt.Errorf("parser: unable to find import path for %q in %q", name, f.path)
}

// ImportName finds an import name by import path
func (f *File) ImportName(path string) (name string, err error) {
	targetPath := strconv.Quote(path)
	for _, imp := range f.node.Imports {
		if imp.Path.Value != targetPath {
			continue
		}
		if imp.Name == nil {
			return imports.AssumedName(path), nil
		}
		return imp.Name.Name, nil
	}
	return "", fmt.Errorf("parser: unable to find import name for %q in %q", path, f.path)
}

// Functions returns all the functions in the file
func (f *File) Functions() (fns []*Function) {
	for _, decl := range f.node.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		fns = append(fns, &Function{
			pkg:  f.pkg,
			file: f,
			node: fn,
		})
	}
	return fns
}

// Structs returns all the structs in a file
func (f *File) Structs() (stcts []*Struct) {
	for _, decl := range f.node.Decls {
		node, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range node.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			stct, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			stcts = append(stcts, &Struct{
				file: f,
				ts:   ts,
				node: stct,
			})
		}
	}
	return stcts
}

// Struct returns a struct by name
func (f *File) Struct(name string) *Struct {
	for _, stct := range f.Structs() {
		if stct.Name() == name {
			return stct
		}
	}
	return nil
}

// Interface returns an interface by name
func (f *File) Interface(name string) *Interface {
	for _, iface := range f.Interfaces() {
		if iface.Name() == name {
			return iface
		}
	}
	return nil
}

// Interfaces returns all the interfaces in a file
func (f *File) Interfaces() (ifaces []*Interface) {
	for _, decl := range f.node.Decls {
		node, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range node.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			iface, ok := ts.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}
			ifaces = append(ifaces, &Interface{
				file: f,
				ts:   ts,
				node: iface,
			})
		}
	}
	return ifaces
}

func (f *File) Alias(name string) *Alias {
	for _, alias := range f.Aliases() {
		if alias.Name() == name {
			return alias
		}
	}
	return nil
}

func (f *File) Aliases() (aliases []*Alias) {
	for _, decl := range f.node.Decls {
		node, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range node.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok || ts.Assign == 0 {
				continue
			}
			aliases = append(aliases, &Alias{
				file: f,
				ts:   ts,
			})
		}
	}
	return aliases
}
