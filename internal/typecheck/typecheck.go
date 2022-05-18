package typecheck

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"path"
	"path/filepath"

	"github.com/livebud/bud/internal/goimport"
	"github.com/livebud/bud/package/gomod"
)

func New(fsys fs.FS, module *gomod.Module) *Checker {
	return &Checker{fsys, goimport.New(fsys), module}
}

type Checker struct {
	fsys     fs.FS
	importer *goimport.Importer
	module   *gomod.Module
}

func (c *Checker) Check(dir string) error {
	pkg, err := c.importer.Import(dir)
	if err != nil {
		return err
	}
	var files []*ast.File
	fset := token.NewFileSet()
	for _, filename := range pkg.GoFiles {
		filename = path.Join(dir, filename)
		code, err := fs.ReadFile(c.fsys, filename)
		if err != nil {
			return err
		}
		file, err := parser.ParseFile(fset, filename, code, parser.DeclarationErrors)
		if err != nil {
			return err
		}
		files = append(files, file)
	}
	importer := &importer{fset, c.fsys, c.importer, c.module}
	tc := types.Config{Importer: importer}
	tpkg := types.NewPackage(pkg.ImportPath, pkg.Name)
	checker := types.NewChecker(&tc, fset, tpkg, nil)
	return checker.Files(files)
}

type importer struct {
	fset     *token.FileSet
	fsys     fs.FS
	importer *goimport.Importer
	module   *gomod.Module
}

var _ types.Importer = (*importer)(nil)

func (i *importer) Import(importPath string) (*types.Package, error) {
	return nil, fmt.Errorf("typecheck: unused")
}
func (i *importer) ImportFrom(importPath, _ string, mode types.ImportMode) (*types.Package, error) {
	// Resolve an import into Go files
	// Note: this is duplicated from parser's (*SelectorType).Definition()
	current := i.module
	module, err := current.FindIn(i.fsys, importPath)
	if err != nil {
		return nil, err
	}
	fsys := i.fsys
	if current.Directory() != module.Directory() {
		// Module is the FS, since we're outside the application dir now
		fsys = module
	}
	dir, err := module.ResolveDirectoryIn(fsys, importPath)
	if err != nil {
		return nil, err
	}
	relDir, err := filepath.Rel(module.Directory(), dir)
	if err != nil {
		return nil, err
	}
	buildPkg, err := goimport.New(fsys).Import(relDir)
	if err != nil {
		return nil, err
	}
	// Parse the Go files
	var files []*ast.File
	fset := token.NewFileSet()
	for _, filename := range buildPkg.GoFiles {
		filename = path.Join(relDir, filename)
		code, err := fs.ReadFile(fsys, filename)
		if err != nil {
			return nil, err
		}
		file, err := parser.ParseFile(fset, filename, code, parser.DeclarationErrors)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	// Run the type-checker on the parsed files
	tc := types.Config{Importer: i}
	pkg, err := tc.Check(importPath, fset, files, nil)
	if err != nil {
		return nil, err
	}
	return pkg, nil
}
