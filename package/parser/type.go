package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/livebud/bud/internal/gois"
)

// Type fn
type Type interface {
	String() string
	node() ast.Expr
}

// Get the expression
// https://golang.org/ref/spec#Types
func getType(f Fielder, x ast.Expr) Type {
	switch t := x.(type) {
	case *ast.Ident:
		return &IdentType{f, t}
	case *ast.SelectorExpr:
		return &SelectorType{f, t}
	case *ast.ArrayType:
		return &ArrayType{f, t}
	case *ast.StructType:
		return &StructType{f, t}
	case *ast.StarExpr:
		return &StarType{f, t}
	case *ast.FuncType:
		return &FuncType{f, t}
	case *ast.InterfaceType:
		return &InterfaceType{f, t}
	case *ast.SliceExpr:
		return &SliceExpr{f, t}
	case *ast.MapType:
		return &MapType{f, t}
	case *ast.ChanType:
		return &ChanType{f, t}
	case *ast.Ellipsis:
		return &EllipsisType{f, t}
	default:
		// Shouldn't happen, but if it does, it's a bug to fix.
		panic(fmt.Errorf("parse: unhandled expression type %T in %q", t, f.File().Path()))
	}
}

// Optional inner interface
type inner interface {
	Inner() Type
}

// Innermost returns the innermost type
// e.g. []*ast.Package becomes ast.Package
// e.g. []*string becomes string
func Innermost(t Type) Type {
	u, ok := t.(inner)
	if !ok {
		return t
	}
	return Innermost(u.Inner())
}

// IsBuiltin returns true if the type is built into Go
func IsBuiltin(t Type) bool {
	return gois.Builtin(TypeName(t))
}

// Qualify adds a package to a type
// e.g. []*Package becomes []*ast.Package
func Qualify(t Type, qualifier string) Type {
	q, ok := t.(qualify)
	if !ok {
		return t
	}
	return q.Qualify(qualifier)
}

// Optional qualify interface
type qualify interface {
	Qualify(qualifier string) Type
}

// Unqualify removes a package from a type
// e.g. []*ast.Package becomes []*Package
func Unqualify(t Type) Type {
	u, ok := t.(unqualify)
	if !ok {
		return t
	}
	return u.Unqualify()
}

// Optional unqualify interface
type unqualify interface {
	Unqualify() Type
}

// Requalify changes the type qualifier
// e.g. []*v8.VM becomes []*js.VM
// e.g. []*string becomes []*string
func Requalify(t Type, replace string) Type {
	st, ok := Innermost(t).(*SelectorType)
	if !ok {
		return t
	}
	// Optimize the common case, where what we're requalifying has the same name.
	id, ok := st.n.X.(*ast.Ident)
	if ok && id.Name == replace {
		return t
	}
	return Qualify(Unqualify(t), replace)
}

// IsImportType checks if the type matches the import type
func IsImportType(t Type, importPath, name string) (bool, error) {
	t = Innermost(t)
	v, ok := t.(*SelectorType)
	if !ok {
		return false, nil
	}
	if v.Name() != name {
		return false, nil
	}
	im, err := v.ImportPath()
	if err != nil {
		return false, err
	}
	if im != importPath {
		return false, nil
	}
	return true, nil
}

// Definition tries going to the type's definition
func Definition(t Type) (Declaration, error) {
	d, ok := t.(definition)
	if !ok {
		return nil, fmt.Errorf("parser: type %q doesn't implement Definition", t)
	}
	return d.Definition()
}

// Optional definition interface
type definition interface {
	Definition() (Declaration, error)
}

// Optional importPath interface
type importPath interface {
	ImportPath() (path string, err error)
}

// ImportPath tries going to the type's definition
func ImportPath(t Type) (path string, err error) {
	ip, ok := t.(importPath)
	if !ok {
		return "", fmt.Errorf("parser: type %q doesn't implement ImportPath", t)
	}
	return ip.ImportPath()
}

// TypeName returns the name of the type
func TypeName(t Type) string {
	tn, ok := t.(typeName)
	if !ok {
		return ""
	}
	return tn.Name()
}

type typeName interface {
	Name() string
}

// FullName does it's best to resolve the full name of a type
func FullName(t Type) string {
	s := new(strings.Builder)
	// Try pulling the import path
	if i, ok := t.(importPath); ok {
		imp, err := i.ImportPath()
		if err != nil {
			// Fallback to the type string
			return t.String()
		}
		s.WriteString(strconv.Quote(imp))
	}
	// Try pulling the name
	if n, ok := t.(typeName); ok {
		name := n.Name()
		s.WriteString("." + name)
		return s.String()
	}
	// Fallback to the type string
	return t.String()
}

// StarType struct
type StarType struct {
	f Fielder
	n *ast.StarExpr
}

// Inner type
func (t *StarType) Inner() Type {
	return getType(t.f, t.n.X)
}

func (t *StarType) String() string {
	x := t.Inner()
	s := x.String()
	return "*" + s
}

func (t *StarType) Name() string {
	x := t.Inner()
	return TypeName(x)
}

// ImportPath returns the import path if there is one
func (t *StarType) ImportPath() (path string, err error) {
	return ImportPath(t.Inner())
}

// expr type
func (t *StarType) node() ast.Expr {
	return t.n
}

// Definition returns the type definition
func (t *StarType) Definition() (Declaration, error) {
	return Definition(t.Inner())
}

// Qualify fn
func (t *StarType) Qualify(qualifier string) Type {
	inner := Qualify(t.Inner(), qualifier)
	return &StarType{
		f: t.f,
		n: &ast.StarExpr{
			X: inner.node(),
		},
	}
}

// Unqualify returns the local type
func (t *StarType) Unqualify() Type {
	unqualified := Unqualify(t.Inner())
	return &StarType{
		f: t.f,
		n: &ast.StarExpr{
			X: unqualified.node(),
		},
	}
}

// IdentType struct
type IdentType struct {
	f Fielder
	n *ast.Ident
}

func (t *IdentType) String() string {
	return t.n.Name
}

// ImportPath returns the import path if there is one
func (t *IdentType) ImportPath() (path string, err error) {
	// Builtins aren't imported
	if gois.Builtin(t.n.Name) {
		return "", nil
	}
	return t.f.File().Import()
}

// expr type
func (t *IdentType) node() ast.Expr {
	return t.n
}

// Qualify fn
func (t *IdentType) Qualify(qualifier string) Type {
	return &SelectorType{
		f: t.f,
		n: &ast.SelectorExpr{
			X:   ast.NewIdent(qualifier),
			Sel: t.n,
		},
	}
}

// Unqualify identifier is the identifier
func (t *IdentType) Unqualify() Type {
	// Nothing to do
	return t
}

// Definition returns the type definition
func (t *IdentType) Definition() (Declaration, error) {
	pkg := t.f.File().Package()
	return pkg.definition(t.n.Name)
}

func (t *IdentType) Name() string {
	return t.n.Name
}

// SelectorType struct
type SelectorType struct {
	f Fielder
	n *ast.SelectorExpr
}

func (t *SelectorType) Name() string {
	return t.n.Sel.Name
}

func (t *SelectorType) String() string {
	x := getType(t.f, t.n.X)
	s := x.String()
	return s + "." + t.n.Sel.Name
}

// ImportPath returns the import path if there is one
func (t *SelectorType) ImportPath() (path string, err error) {
	pkg, ok := t.n.X.(*ast.Ident)
	if !ok {
		return "", fmt.Errorf("parser: unknown selector type %T", t.n.X)
	}
	return t.f.File().ImportPath(pkg.Name)
}

// expr type
func (t *SelectorType) node() ast.Expr {
	return t.n
}

// Qualify fn
func (t *SelectorType) Qualify(qualifier string) Type {
	// Nothing to do
	return t
}

// Unqualify returns the selector type
func (t *SelectorType) Unqualify() Type {
	return &IdentType{
		f: t.f,
		n: t.n.Sel,
	}
}

// Definition returns the type definition
func (t *SelectorType) Definition() (Declaration, error) {
	decl, err := t.definition()
	if err != nil {
		return nil, fmt.Errorf("parser: unable to find declaration for %s in %q. %w", FullName(t), t.f.File().Path(), err)
	}
	return decl, nil
}

// definition tries finding the definition for the selector type
func (t *SelectorType) definition() (Declaration, error) {
	left, ok := t.n.X.(*ast.Ident)
	if !ok {
		return nil, fmt.Errorf("parser: unknown selector type %T", t.n.X)
	}
	file := t.f.File()
	pkg := file.Package()
	imp, err := file.ImportPath(left.Name)
	if err != nil {
		return nil, err
	}
	current := pkg.Module()
	module, err := current.FindIn(pkg.FS(), imp)
	if err != nil {
		return nil, err
	}
	fsys := pkg.FS()
	if current.Directory() != module.Directory() {
		// Module is the FS, since we're outside the application dir now
		fsys = module
	}
	dir, err := module.ResolveDirectoryIn(fsys, imp)
	if err != nil {
		return nil, err
	}
	rel, err := filepath.Rel(module.Directory(), dir)
	if err != nil {
		return nil, err
	}
	newPkg, err := New(fsys, module).Parse(rel)
	if err != nil {
		return nil, err
	}
	return newPkg.definition(t.n.Sel.Name)
}

// ArrayType struct
type ArrayType struct {
	f Fielder
	n *ast.ArrayType
}

// Inner type
func (t *ArrayType) Inner() Type {
	return getType(t.f, t.n.Elt)
}

func (t *ArrayType) Name() string {
	x := t.Inner()
	return TypeName(x)
}

// TODO: handle len in [len]elt
func (t *ArrayType) String() string {
	inner := t.Inner()
	s := inner.String()
	return "[]" + s
}

// ImportPath returns the import path if there is one
func (t *ArrayType) ImportPath() (path string, err error) {
	return ImportPath(t.Inner())
}

// expr type
func (t *ArrayType) node() ast.Expr {
	return t.n
}

// Qualify fn
func (t *ArrayType) Qualify(qualifier string) Type {
	elt := Qualify(t.Inner(), qualifier)
	return &ArrayType{
		f: t.f,
		n: &ast.ArrayType{
			Len: t.n.Len,
			Elt: elt.node(),
		},
	}
}

// Unqualify returns the type if you were referring to it within the same
// package
func (t *ArrayType) Unqualify() Type {
	inner := t.Inner()
	shifted := Unqualify(inner)
	return &ArrayType{
		f: t.f,
		n: &ast.ArrayType{
			Len: t.n.Len,
			Elt: shifted.node(),
		},
	}
}

// Definition returns the type definition
func (t *ArrayType) Definition() (Declaration, error) {
	inner := t.Inner()
	return Definition(inner)
}

// StructType struct
type StructType struct {
	f Fielder
	n *ast.StructType
}

var _ Type = (*StructType)(nil)

// String fn
func (t *StructType) String() string {
	return printExpr(t.n)
}

// expr fn
func (t *StructType) node() ast.Expr {
	return t.n
}

// FuncType struct
type FuncType struct {
	f Fielder
	n *ast.FuncType
}

var _ Type = (*FuncType)(nil)

// String fn
func (t *FuncType) String() string {
	return printExpr(t.n)
}

// expr fn
func (t *FuncType) node() ast.Expr {
	return t.n
}

// InterfaceType struct
type InterfaceType struct {
	f Fielder
	n *ast.InterfaceType
}

var _ Type = (*InterfaceType)(nil)

// String fn
func (t *InterfaceType) String() string {
	return printExpr(t.n)
}

// expr fn
func (t *InterfaceType) node() ast.Expr {
	return t.n
}

// SliceExpr struct
type SliceExpr struct {
	f Fielder
	n *ast.SliceExpr
}

var _ Type = (*SliceExpr)(nil)

// String fn
func (t *SliceExpr) String() string {
	return printExpr(t.n)
}

// expr fn
func (t *SliceExpr) node() ast.Expr {
	return t.n
}

// MapType struct
type MapType struct {
	f Fielder
	n *ast.MapType
}

var _ Type = (*MapType)(nil)

// String fn
func (t *MapType) String() string {
	return printExpr(t.n)
}

// expr fn
func (t *MapType) node() ast.Expr {
	return t.n
}

// ChanType struct
type ChanType struct {
	f Fielder
	n *ast.ChanType
}

var _ Type = (*ChanType)(nil)

// String fn
func (t *ChanType) String() string {
	return printExpr(t.n)
}

// expr fn
func (t *ChanType) node() ast.Expr {
	return t.n
}

// Ellipsis struct
type EllipsisType struct {
	f Fielder
	n *ast.Ellipsis
}

var _ Type = (*EllipsisType)(nil)

// String fn
func (t *EllipsisType) String() string {
	return printExpr(t.n)
}

// String fn
func (t *EllipsisType) node() ast.Expr {
	return t.n
}

// Inner type
func (t *EllipsisType) Inner() Type {
	return getType(t.f, t.n.Elt)
}

func (t *EllipsisType) Name() string {
	x := t.Inner()
	return TypeName(x)
}

// ImportPath returns the import path if there is one
func (t *EllipsisType) ImportPath() (path string, err error) {
	return ImportPath(t.Inner())
}

// Definition returns the declaration
func (t *EllipsisType) Definition() (Declaration, error) {
	return Definition(t.Inner())
}

// printExpr prints an expression
// TODO: benchmark, we use type.String() a lot and this might be slow
func printExpr(expr ast.Expr) string {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	_ = printer.Fprint(&buf, fset, expr)
	return buf.String()
}
