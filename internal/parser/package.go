package parser

import (
	"fmt"
	"go/ast"
	"sort"

	"gitlab.com/mnm/bud/go/is"
	"gitlab.com/mnm/bud/go/mod"
)

// newPackage creates a new package
func newPackage(directory string, node *ast.Package, parser *Parser) *Package {
	pkg := &Package{
		node:      node,
		directory: directory,
		parser:    parser,
	}
	pkg.files = files(pkg)
	return pkg
}

// Package struct
type Package struct {
	directory string
	files     []*File
	node      *ast.Package
	parser    *Parser
}

// Name of the package
func (pkg *Package) Name() string {
	return pkg.node.Name
}

// Directory returns the directory of the package
func (pkg *Package) Directory() string {
	return pkg.directory
}

// Files returns a list of files
func (pkg *Package) Files() []*File {
	return pkg.files
}

// Modfile returns the modfile or fails
func (pkg *Package) Modfile() (*mod.File, error) {
	return pkg.parser.modfile(pkg.directory)
}

// Import returns the import path to this package
func (pkg *Package) Import() (path string, err error) {
	modfile, err := pkg.Modfile()
	if err != nil {
		return "", err
	}
	return modfile.ResolveImport(pkg.directory)
}

// ResolveDirectory resolves a directory from an import path
func (pkg *Package) ResolveDirectory(importPath string) (string, error) {
	modfile, err := pkg.Modfile()
	if err != nil {
		return "", err
	}
	return modfile.ResolveDirectory(importPath)
}

// ResolveImport resolves a directory from an import path
func (pkg *Package) ResolveImport(directory string) (string, error) {
	modfile, err := pkg.Modfile()
	if err != nil {
		return "", err
	}
	return modfile.ResolveImport(directory)
}

// ModuleDirectory returns the directory that contains go.mod.
func (pkg *Package) ModuleDirectory() (string, error) {
	modfile, err := pkg.Modfile()
	if err != nil {
		return "", err
	}
	return modfile.Directory(), nil
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

// Declaration interface
type Declaration interface {
	Name() string
	Package() *Package
	Directory() string
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
	err = fmt.Errorf("parser: Unable to find declaration for %q", name)
	var ts *ast.TypeSpec
	for _, file := range pkg.Files() {
		file := file
		ast.Inspect(pkg.node, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.TypeSpec:
				ts = n
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

// Directory for builtin is blank
func (b builtin) Directory() string {
	return ""
}

// Package for builtin is blank
func (b builtin) Package() *Package {
	return nil
}
