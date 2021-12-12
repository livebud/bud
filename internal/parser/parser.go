package parser

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"path/filepath"

	"gitlab.com/mnm/bud/go/mod"
)

// New Go parser.
func New(modFinder *mod.Finder) *Parser {
	return &Parser{
		modFinder: modFinder,
	}
}

// Parser for parsing Go code.
type Parser struct {
	modFinder *mod.Finder
}

// Parse a dir containing Go files.
func (p *Parser) Parse(dir string) (*Package, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	// Imports the package, taking into account build tags and file name
	// conventions
	imported, err := build.ImportDir(dir, build.ImportMode(0))
	if err != nil {
		return nil, err
	}
	parsedPackage := &ast.Package{
		Name:  imported.Name,
		Files: make(map[string]*ast.File),
	}
	fset := token.NewFileSet()
	// Parse each valid Go file
	for _, path := range imported.GoFiles {
		filename := filepath.Join(dir, path)
		parsedFile, err := parser.ParseFile(fset, filename, nil, parser.DeclarationErrors)
		if err != nil {
			return nil, err
		}
		parsedPackage.Files[path] = parsedFile
	}
	pkg := newPackage(dir, p.modFinder, parsedPackage, p)
	return pkg, nil
}
