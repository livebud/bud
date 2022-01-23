package di

import (
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/imports"
)

type Function2 struct {
	Name    string
	ModFile *mod.File
	Hoist   bool
	Params  []*Param
	Results []Result
}

// Param contains external data types
type Param struct {
}

// Result contains the dependencies we provide
type Result interface {
}

type Struct2 struct {
	Name   string
	Fields []*StructField
}

type StructField2 struct {
	Name string
	Type *Type
}

type Type struct {
	Import string // Import path (e.g. myapp.com/web)
	Type   string // Data type (e.g. *Web)
}

// ID of the dependency as a string
func (t *Type) ID() string {
	return `"` + t.Import + `".` + t.Type
}

// Error type
type ErrorType struct{}

func (i *Injector) Wire(fn *Function2) (*Node2, error) {

	return nil, nil
}

type Node2 struct {
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

func (n *Node2) Generate() *Provider2 {
	return nil
}

type Provider2 struct {
}

func (p *Provider2) Imports() (imports []*imports.Import) {
	return imports
}

// Caller generates a caller with the following params
func (p *Provider2) Caller(params ...string) (string, error) {
	return "", nil
}

// Generate the wired graph into inline code
func (p *Provider2) Generate() string {
	return ""
}

// GenerateFile generates the wired graph into a separate file
func (p *Provider2) GenerateFile() string {
	return ""
}
