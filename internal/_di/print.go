package di

import (
	"fmt"
	"strings"
)

type PrintInput struct {
	// Targets to load
	Dependencies []*Dependency `json:"dependencies,omitempty"`
	// External types
	Externals []*Dependency `json:"externals,omitempty"`
	// Hoist dependencies that don't depend on externals, turning them into
	// externals. This is to avoid initializing these inner deps every time.
	// Useful for per-request dependency injection.
	Hoist bool `json:"hoist,omitempty"`
}

// Print the dependencies to the dot format to display with GraphViz.
// I use: https://dreampuf.github.io/GraphvizOnline
func (i *Injector) Print(in *PrintInput) (string, error) {
	dep, err := i.Load(&LoadInput{
		Dependencies: in.Dependencies,
		Externals:    in.Externals,
		Hoist:        in.Hoist,
	})
	if err != nil {
		return "", err
	}
	return dep.Print(), nil
}

func (node *Node) Print() (out string) {
	out = "digraph G {\n"
	seen := map[string]bool{}
	out += "  " + print(node, seen)
	out += "\n}\n"
	return out
}

func print(node *Node, seen map[string]bool) string {
	id := node.ID()
	if seen[id] {
		return ""
	}
	seen[id] = true
	var outs []string
	for _, dep := range node.Dependencies {
		str := new(strings.Builder)
		label := dep.Original.Type
		fmt.Fprintf(str, `%q -> %q`, dep.ID(), id)
		if dep.External {
			label += " (external)"
		}
		fmt.Fprintf(str, ` [label=%q];`, label)
		outs = append(outs, str.String())
		subgraph := print(dep, seen)
		if subgraph == "" {
			continue
		}
		outs = append(outs, subgraph)
	}
	return strings.Join(outs, "\n  ")
}
