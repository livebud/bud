package di

import (
	"errors"
	"fmt"
	"os"

	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/parser"
)

type Declaration interface {
	ID() string
	Dependencies() []*Dependency
	Generate(gen *Generator, inputs []*Variable) (outputs []*Variable)
}

type Variable struct {
	Import string // Import path
	Type   string // Type of the variable
	Name   string // Name of the variable
}

type External struct {
	*Variable
	Key     string // Name to be used as a key in a struct
	Hoisted bool   // True if this external was hoisted up
}

// The default searcher just searches the dependency's import path
var defaultSearcher = func(importPath string) []string {
	return []string{importPath}
}

// New dependency injector
func New(parser *parser.Parser) *Injector {
	return &Injector{
		Parser:   parser,
		Searcher: defaultSearcher,
	}
}

type Injector struct {
	// Parser for parsing Go code
	Parser *parser.Parser

	// Searcher function defines where to search for dependencies.
	// Usually the searchPath is the same as a dependency's import,
	// but it doesn't have to be. For example,
	//
	//    package log
	//    import "duo/log"
	//    func New() log.Log {}
	//
	// In this case where we instantiate `log.Log` is in a different package from
	// where it's defined.
	//
	// By default, SearchPaths contains only the dependency's Import path.
	Searcher Searcher
}

type Searcher = func(importPath string) (searchPaths []string)

type Type struct {
	// Import Path
	Import string `json:"import,omitempty"`
	// Type of the target (e.g. *Web)
	Type string `json:"type,omitempty"`
	// Module file containing this type
	ModFile *mod.File
}

func (t *Type) Provider() (Declaration, error) {
	return nil, nil
}

type StructType struct {
	// Module file containing this type
	ModFile *mod.File
	// Import Path
	Import string `json:"import,omitempty"`
	// Type of the target (e.g. *Web)
	Type   string `json:"type,omitempty"`
	Fields []*StructTypeField
}

type StructTypeField struct {
	Name   string // Field or parameter name
	Import string // Import path
	Type   string // Field type
}

func (s *StructType) Provider() (Declaration, error) {

	// return &Struct{
	// 	NeedsRef: s.
	// }, nil
	return nil, nil
}

// Dependency struct
type Dependency struct {
	// Import Path
	Import string `json:"import,omitempty"`
	// Type of the target (e.g. *Web)
	Type string `json:"type,omitempty"`
	// Module file for this import
	// Currently nil for externals (TODO: re-consider)
	ModFile *mod.File
}

// ID of the dependency as a string
func (d *Dependency) ID() string {
	return `"` + d.Import + `".` + d.Type
}

// Declaration that would provide this dependency
func (d *Dependency) Declaration(parser *parser.Parser, searcher Searcher) (Declaration, error) {
	searchPaths := searcher(d.Import)
	return d.findDeclaration(parser, searchPaths, []string{})
}

func (d *Dependency) findDeclaration(parser *parser.Parser, searchPaths []string, foundPaths []string) (Declaration, error) {
	// No more search paths, we're unable to find this dependency
	if len(searchPaths) == 0 {
		if len(foundPaths) == 0 {
			return nil, fmt.Errorf("di: unable to find dependency %q.%s", d.Import, d.Type)
		} else {
			return nil, fmt.Errorf("di: unclear how to provide %q.%s", d.Import, d.Type)
		}
	}
	// Resolve the absolute directory based on the import
	dir, err := d.ModFile.ResolveDirectory(searchPaths[0])
	if err != nil {
		// If the directory doesn't exist, search the next search path
		if errors.Is(err, os.ErrNotExist) {
			return d.findDeclaration(parser, searchPaths[1:], foundPaths)
		}
		return nil, err
	}
	// Parse the package
	pkg, err := parser.Parse(dir)
	if err != nil {
		return nil, err
	}
	// Look through the functions
	for _, fn := range pkg.Functions() {
		decl, err := tryFunction(fn, d)
		if err != nil {
			if err == ErrNoMatch {
				continue
			}
			return nil, err
		}
		return decl, nil
	}
	// Look through the structs
	for _, stct := range pkg.Structs() {
		decl, err := tryStruct(stct, d)
		if err != nil {
			if err == ErrNoMatch {
				continue
			}
			return nil, err
		}
		return decl, nil
	}
	// Add to foundPaths and try again
	foundPaths = append(foundPaths, dir)
	// Search the next search path
	decl, err := d.findDeclaration(parser, searchPaths[1:], foundPaths)
	if err != nil {
		return nil, err
	}
	return decl, nil
}

// Node in the dependency injection graph
type Node struct {
	// Original dependency we were looking for. Careful, the Import may not be
	// correct if the search path found this dependency in a different package.
	Original *Dependency
	// Declaration that would instantiate this type. This will be nil if the node
	// is External
	Declaration Declaration
	// Dependencies that the declaration relies on to be able to instantiate
	Dependencies []*Node
	// External types are passed in, not instantiated. This is true if the type
	// matches an external dependency or this dependency is hoisted up.
	External bool
}

// ID returns the declaration's ID or the dependency ID (for externals).
func (node *Node) ID() string {
	id := node.Original.ID()
	if node.Declaration != nil {
		id = node.Declaration.ID()
	}
	return id
}
