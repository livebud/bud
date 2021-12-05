package di

import (
	"gitlab.com/mnm/bud/go/mod"
	"gitlab.com/mnm/bud/internal/parser"
)

type Map = map[Dependency]Dependency

func New(modFile *mod.File, parser *parser.Parser, typeMap Map) *Injector {
	tm := map[string]Dependency{}
	for from, to := range typeMap {
		tm[from.ID()] = to
	}
	return &Injector{modFile, parser, tm}
}

type Injector struct {
	// Modfile of the module where project dependencies will be wired
	modFile *mod.File
	// Go parser
	parser *parser.Parser
	// Type aliasing
	typeMap map[string]Dependency
}

// Load the dependency graph, but don't generate any code. Load is intentionally
// low-level and used by higher-level APIs like Generate.
func (i *Injector) Load(fn *Function) (*Node, error) {
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
		node, err := i.load(externals, result)
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
func (i *Injector) load(externals map[string]bool, dep Dependency) (*Node, error) {
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
		child, err := i.load(externals, dep)
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
		return nil, err
	}
	return node.Generate(fn.Target), nil
}

// GenerateFile generates a provider function into string
func (i *Injector) Generate(fn *Function) (string, error) {
	provider, err := i.Wire(fn)
	if err != nil {
		return "", err
	}
	code := provider.Function(fn.Name)
	return code, nil
}

// GenerateFile generates a provider function into a Go file string
func (i *Injector) GenerateFile(fn *Function) (string, error) {
	provider, err := i.Wire(fn)
	if err != nil {
		return "", err
	}
	code := provider.File(fn.Name)
	return code, nil
}
