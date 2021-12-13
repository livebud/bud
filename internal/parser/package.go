package parser

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"sort"

	"gitlab.com/mnm/bud/go/is"
	"gitlab.com/mnm/bud/go/mod"
)

// newPackage creates a new package
func newPackage(dir string, module *mod.Module, node *ast.Package) *Package {
	absdir := filepath.Join(module.Directory(), dir)
	pkg := &Package{
		node:   node,
		dir:    absdir,
		module: module,
	}
	pkg.files = files(pkg)
	return pkg
}

// Package struct
type Package struct {
	dir    string
	files  []*File
	node   *ast.Package
	module *mod.Module
}

// Name of the package
func (pkg *Package) Name() string {
	return pkg.node.Name
}

// Directory returns the directory of the package
func (pkg *Package) Directory() string {
	return pkg.dir
}

// Files returns a list of files
func (pkg *Package) Files() []*File {
	return pkg.files
}

// Module returns the module or fails
func (pkg *Package) Module() *mod.Module {
	return pkg.module
}

// Import returns the import path to this package
func (pkg *Package) Import() (string, error) {
	return pkg.module.ResolveImport(pkg.dir)
}

// ResolveDirectory resolves a directory from an import path
func (pkg *Package) ResolveDirectory(importPath string) (string, error) {
	return pkg.Module().ResolveDirectory(importPath)
}

// ResolveImport resolves a directory from an import path
func (pkg *Package) ResolveImport(directory string) (string, error) {
	return pkg.Module().ResolveImport(directory)
}

// Files returns the Go files within the package
func files(pkg *Package) (files []*File) {
	for path, node := range pkg.node.Files {
		files = append(files, &File{
			pkg:  pkg,
			node: node,
			path: path,
		})
	}
	// Stable file sorting within the package
	sort.Slice(files, func(i, j int) bool {
		return files[i].path < files[j].path
	})
	return files
}

// Kind of declaration
type Kind uint8

const (
	KindBuiltin Kind = iota + 1
	KindStruct
	KindInterface
	KindAlias
)

// Declaration interface
type Declaration interface {
	Name() string
	Package() *Package
	Kind() Kind
}

// Functions returns all the functions in a package
func (pkg *Package) Functions() (fns []*Function) {
	for _, file := range pkg.Files() {
		fns = append(fns, file.Functions()...)
	}
	return fns
}

// PublicFunctions returns all public functions in the package
func (pkg *Package) PublicFunctions() (fns []*Function) {
	for _, file := range pkg.Files() {
		for _, fn := range file.Functions() {
			if fn.Private() {
				continue
			}
			fns = append(fns, fn)
		}
	}
	return fns
}

// PublicMethods returns all public methods in the package
func (pkg *Package) PublicMethods() (fns []*Function) {
	for _, file := range pkg.Files() {
		for _, fn := range file.Functions() {
			if fn.Private() {
				continue
			}
			if fn.Receiver() == nil {
				continue
			}
			fns = append(fns, fn)
		}
	}
	return fns
}

// Structs returns all the structs in a package
func (pkg *Package) Structs() (stcts []*Struct) {
	for _, file := range pkg.Files() {
		stcts = append(stcts, file.Structs()...)
	}
	return stcts
}

// Struct returns a struct by name
func (pkg *Package) Struct(name string) *Struct {
	for _, file := range pkg.Files() {
		if stct := file.Struct(name); stct != nil {
			return stct
		}
	}
	return nil
}

func (pkg *Package) Interface(name string) *Interface {
	for _, file := range pkg.Files() {
		if iface := file.Interface(name); iface != nil {
			return iface
		}
	}
	return nil
}

// Interfaces returns all the interfaces in the package
func (pkg *Package) Interfaces() (ifaces []*Interface) {
	for _, file := range pkg.Files() {
		ifaces = append(ifaces, file.Interfaces()...)
	}
	return ifaces
}

func (pkg *Package) Alias(name string) *Alias {
	for _, file := range pkg.Files() {
		if alias := file.Alias(name); alias != nil {
			return alias
		}
	}
	return nil
}

func (pkg *Package) Aliases() (aliases []*Alias) {
	for _, file := range pkg.Files() {
		aliases = append(aliases, file.Aliases()...)
	}
	return aliases
}

// var errIsBuiltin = errors.New("definition is a built-in type")

// // ErrIsBuiltin checks if the error is builtin
// func ErrIsBuiltin(err error) bool {
// 	return errors.Is(err, errIsBuiltin)
// }

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
