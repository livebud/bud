package di

type LoadInput struct {
	// Targets to load
	Dependencies []*Dependency `json:"dependencies,omitempty"`
	// External types
	Externals []*Dependency `json:"externals,omitempty"`
	// Hoist dependencies that don't depend on externals, turning them into
	// externals. This is to avoid initializing these inner deps every time.
	// Useful for per-request dependency injection.
	Hoist bool `json:"hoist,omitempty"`

	// Filled in during Load
	externals map[string]bool
}

// Load the dependency graph, but don't generate any code. Load is intentionally
// low-level and used by higher-level APIs like Generate.
func (i *Injector) Load(in *LoadInput) (*Node, error) {
	in.externals = map[string]bool{}
	for _, external := range in.Externals {
		in.externals[external.ID()] = true
	}
	node, err := i.load(in, in.Dependencies[0])
	if err != nil {
		return nil, err
	}
	if in.Hoist {
		node = Hoist(node)
	}
	return node, err
}

// Load the dependencies recursively. This produces a dependency graph of nodes.
func (i *Injector) load(in *LoadInput, dep *Dependency) (*Node, error) {
	// Handle external nodes
	if in.externals[dep.ID()] {
		return &Node{
			External: true,
			Original: dep,
		}, nil
	}

	// Find the declaration that would instantiate this dependency
	decl, err := i.Find(dep)
	if err != nil {
		return nil, err
	}
	node := &Node{
		Original:    dep,
		Declaration: decl,
	}
	// Get the Declaration's dependencies
	deps := decl.Dependencies()
	// Find and load the dependencies
	for _, dep := range deps {
		child, err := i.load(in, dep)
		if err != nil {
			return nil, err
		}
		node.Dependencies = append(node.Dependencies, child)
	}
	return node, nil
}
