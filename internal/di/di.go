package di

import (
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
	Searcher func(importPath string) (searchPaths []string)
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
