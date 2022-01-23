package parser

import (
	"go/ast"
	"path/filepath"
	"sort"

	"gitlab.com/mnm/bud/go/mod"
)

// newPackage creates a new package
func newPackage(dir string, module *mod.Module, node *ast.Package) *Package {
	dir = filepath.Join(module.Directory(), dir)
	pkg := &Package{
		node:   node,
		dir:    dir,
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
	return pkg.module.ResolveDirectory(importPath)
}

// ResolveImport resolves a directory from an import path
func (pkg *Package) ResolveImport(directory string) (string, error) {
	return pkg.module.ResolveImport(directory)
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
type Kind string

const (
	KindBuiltin   Kind = "builtin"
	KindStruct         = "struct"
	KindInterface      = "interface"
	KindAlias          = "alias"
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
