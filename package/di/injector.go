package di

import (
	"fmt"
	"io/fs"

	"github.com/livebud/bud/internal/imports"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/parser"
)

func New(fsys fs.FS, module *gomod.Module, parser *parser.Parser) *Injector {
	return &Injector{
		fsys:   fsys,
		module: module,
		parser: parser,
	}
}

type Injector struct {
	// Filesystem to look for files
	fsys fs.FS
	// Module where project dependencies will be wired
	module *gomod.Module
	// Go parser
	parser *parser.Parser
}

// Load the dependency graph, but don't generate any code. Load is intentionally
// low-level and used by higher-level APIs like Generate.
func (i *Injector) Load(fn *Function) (*Node, error) {
	// Validate the function
	if err := fn.Validate(); err != nil {
		return nil, err
	}
	// Setup the aliases
	aliases := map[string]Dependency{}
	for from, to := range fn.Aliases {
		aliases[from.ID()] = to
	}
	externals := map[string]bool{}
	for _, param := range fn.Params {
		id := param.ID()
		externals[id] = true
	}
	root := &Node{
		Import:      fn.Target,
		Type:        fn.Name,
		External:    false,
		Declaration: fn,
	}
	// Load the dependencies
	for _, result := range fn.Results {
		node, err := i.load(externals, aliases, result)
		if err != nil {
			return nil, err
		}
		root.Dependencies = append(root.Dependencies, node)
	}
	if fn.Hoist {
		root = Hoist(root)
	}
	return root, nil
}

// Load the dependencies recursively. This produces a dependency graph of nodes.
func (i *Injector) load(externals map[string]bool, aliases map[string]Dependency, dep Dependency) (*Node, error) {
	// Replace dep with mapped type alias if we have one
	if alias, ok := aliases[dep.ID()]; ok {
		dep = alias
	}
	// Handle external nodes
	importPath := dep.ImportPath()
	typeName := dep.TypeName()
	id := dep.ID()
	if externals[id] {
		return &Node{
			Import:   importPath,
			Type:     typeName,
			External: true,
		}, nil
	}
	// Find the declaration that would instantiate this dependency
	decl, err := dep.Find(i)
	if err != nil {
		return nil, err
	}
	node := &Node{
		Import:      importPath,
		Type:        typeName,
		Declaration: decl,
	}
	// Get the Declaration's dependencies
	deps := decl.Dependencies()
	// Find and load the dependencies
	for _, dep := range deps {
		child, err := i.load(externals, aliases, dep)
		if err != nil {
			return nil, err
		}
		node.Dependencies = append(node.Dependencies, child)
	}
	return node, nil
}

// Wire up the provider function into provider state. The Provider has some
// helper functions that are useful when passed into a template.
func (i *Injector) Wire(fn *Function) (*Provider, error) {
	node, err := i.Load(fn)
	if err != nil {
		return nil, fmt.Errorf("di: unable to wire %q.%s function. %w", fn.Target, fn.Name, err)
	}
	if fn.Imports == nil {
		fn.Imports = imports.New()
	}
	return node.Generate(fn.Imports, fn.Name, fn.Target), nil
}

// GenerateFile generates a provider function into string
func (i *Injector) Generate(fn *Function) (string, error) {
	provider, err := i.Wire(fn)
	if err != nil {
		return "", err
	}
	code := provider.Function()
	return code, nil
}

// GenerateFile generates a provider function into a Go file string
func (i *Injector) GenerateFile(fn *Function) (string, error) {
	provider, err := i.Wire(fn)
	if err != nil {
		return "", err
	}
	code := provider.File()
	return code, nil
}
