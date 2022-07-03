package di

// Hoist the nodes that don't depend on the external nodes and turn these
// nodes into external nodes. This allows for dependencies that don't depend
// on externals to be initialized once, rather than each time the generated
// function is called.
//
// Start with hoisting true, but if we encounter any external along the way, the
// hoisting of all children becomes false
func Hoist(root *Node) *Node {
	// Hoisting only applies to ancestor dependencies
	for _, result := range root.Dependencies {
		for _, dep := range result.Dependencies {
			hoist(dep)
		}
	}
	return root
}

func hoist(node *Node) (shouldHoist bool) {
	// Default to hoisting
	shouldHoist = true
	// Dependencies that rely on an external node cannot be hoisted.
	if node.External && !node.Hoist {
		return false
	}
	// Loop over the inputs. If any input is non-hoistable, this node is becomes
	// non-hoistable. Order of the conditional matters here. We intentionally call
	// hoist(dep) before shouldHoist because we don't want the algorithm skipping
	// over the subtrees.
	for _, dep := range node.Dependencies {
		shouldHoist = hoist(dep) && shouldHoist
	}
	// If shouldHoist is true, we externalize the node.
	node.Hoist = shouldHoist
	return shouldHoist
}
