package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"

	"gitlab.com/mnm/bud/go/is"
)

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
	default:
		// Shouldn't happen, but if it does, it's a bug to fix.
		panic(fmt.Errorf("parse: unhandled expression type %T in %q", f.File().Path(), t))
	}
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

// ImportPath tries going to the type's definition
func ImportPath(t Type) (path string, err error) {
	ip, ok := t.(importPath)
	if !ok {
		return "", fmt.Errorf("parser: type %q doesn't implement ImportPath", t)
	}
	return ip.ImportPath()
}

// Optional importPath interface
type importPath interface {
	ImportPath() (path string, err error)
}

// TypeName returns the name of the type
func TypeName(t Type) string {
	tn, ok := t.(typeName)
	if !ok {
		return ""
	}
	name := tn.Name()
	if name == "error" {
		return "err"
	}
	return name
}

type typeName interface {
	Name() string
}

// Type fn
type Type interface {
	String() string
	node() ast.Expr
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
	if is.Builtin(t.n.Name) {
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
	modfile, err := pkg.Modfile()
	if err != nil {
		return nil, err
	}
	dir, err := modfile.ResolveDirectory(imp)
	if err != nil {
		return nil, err
	}
	newPkg, err := pkg.parser.Parse(dir)
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

// printExpr prints an expression
func printExpr(expr ast.Expr) string {
	var buf bytes.Buffer
	fset := token.NewFileSet()
	_ = printer.Fprint(&buf, fset, expr)
	return buf.String()
}
