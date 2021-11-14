package parser

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/mnm/bud/go/mod"
)

// New parser. Packages are cached between parses
func New(modfinder *mod.Finder) *Parser {
	return &Parser{
		cache:     newCache(),
		modfinder: modfinder,
	}
}

// Parser for parsing Go code.
type Parser struct {
	cache     *cache
	modfinder *mod.Finder
}

// Parse a dir containing Go files.
func (p *Parser) Parse(dir string) (*Package, error) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	pkgs, err := p.parseDir(fset, dir, ignoreFilter, parser.DeclarationErrors)
	if err != nil {
		return nil, err
	}
	if len(pkgs) == 0 {
		return nil, fmt.Errorf("ast: Parse couldn't find Go code in %q", dir)
	}
	firstPkg := firstPackage(pkgs)
	pkg := newPackage(dir, firstPkg, p)
	return pkg, nil
}

// parseDir is originally from parser.ParserDir. We just added caching.
// TODO: parse files concurrently.
func (p *Parser) parseDir(fset *token.FileSet, path string, filter func(fs.FileInfo) bool, mode parser.Mode) (pkgs map[string]*ast.Package, first error) {
	list, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	pkgs = make(map[string]*ast.Package)
	for _, d := range list {
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".go") {
			continue
		}
		if filter != nil {
			info, err := d.Info()
			if err != nil {
				return nil, err
			}
			if !filter(info) {
				continue
			}
		}
		filename := filepath.Join(path, d.Name())
		// First check if the cache has this package
		src, ok := p.cache.Get(filename)
		if !ok {
			src, err = parser.ParseFile(fset, filename, nil, mode)
		}
		if nil == err {
			name := src.Name.Name
			pkg, found := pkgs[name]
			if !found {
				pkg = &ast.Package{
					Name:  name,
					Files: make(map[string]*ast.File),
				}
				pkgs[name] = pkg
			}
			pkg.Files[filename] = src
			p.cache.Set(filename, src)
		} else if first == nil {
			first = err
		}
	}
	return
}

func (p *Parser) modfile(dir string) (*mod.File, error) {
	return p.modfinder.Find(dir)
}

// Filter for ignoring Go code
func ignoreFilter(fi os.FileInfo) bool {
	name := fi.Name()
	return filepath.Ext(name) == ".go" &&
		!strings.HasSuffix(name, "_test.go")
}

// firstPackage returns the first package we can find
func firstPackage(pkgs map[string]*ast.Package) *ast.Package {
	name := ""
	for name = range pkgs {
		break
	}
	return pkgs[name]
}
