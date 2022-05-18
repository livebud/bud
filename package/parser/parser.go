package parser

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/fs"
	"path"
	"unicode"

	"github.com/livebud/bud/internal/goimport"
	"github.com/livebud/bud/package/gomod"
)

// New Go parser.
func New(fsys fs.FS, module *gomod.Module) *Parser {
	return &Parser{
		fsys:     fsys,
		importer: goimport.New(fsys),
		module:   module,
	}
}

// Parser for parsing Go code.
type Parser struct {
	fsys     fs.FS
	importer *goimport.Importer
	module   *gomod.Module
}

// Parse a dir containing Go files.
func (p *Parser) Parse(dir string) (*Package, error) {
	imported, err := p.Import(dir)
	if err != nil {
		return nil, err
	}
	parsedPackage := &ast.Package{
		Name:  imported.Name,
		Files: make(map[string]*ast.File),
	}
	fset := token.NewFileSet()
	// Parse each valid Go file
	for _, filename := range imported.GoFiles {
		filename = path.Join(dir, filename)
		code, err := fs.ReadFile(p.fsys, filename)
		if err != nil {
			return nil, err
		}
		parsedFile, err := parser.ParseFile(fset, filename, code, parser.DeclarationErrors)
		if err != nil {
			return nil, err
		}
		parsedPackage.Files[filename] = parsedFile
	}
	pkg := newPackage(dir, p, p.module, parsedPackage)
	return pkg, nil
}

// Import the package, taking into account build tags and file name conventions
func (p *Parser) Import(dir string) (*build.Package, error) {
	return p.importer.Import(dir)
}

func isPrivate(name string) bool {
	return unicode.IsLower(rune(name[0]))
}
