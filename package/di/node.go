package di

import (
	"fmt"
	"sort"
	"strings"

	"github.com/livebud/bud/internal/imports"
)

// node in the dependency injection graph
type Node struct {
	// Type and import of the dependency we were looking for
	Import string
	Type   string

	// Declaration that would instantiate this type. This will be nil if the node
	// is External
	Declaration Declaration
	// Dependencies that the declaration relies on to be able to instantiate
	Dependencies []*Node
	// External is true if the type matches an external dependency. External types
	// are passed in, not instantiated.
	External bool
	// Hoisted is true if the dependency has been hoisted up. Hoisted types are
	// passed in, not instantiated.
	Hoist bool
}

func (n *Node) ID() string {
	if n.Declaration != nil {
		return n.Declaration.ID()
	}
	return getID(n.Import, n.Type)
}

// Build a provider for the target import path
func (n *Node) Generate(imports *imports.Set, fnName, target string) *Provider {
	// Build context
	g := &generator{
		Seen:    map[string][]*Variable{},
		Code:    new(strings.Builder),
		Imports: imports,
		Target:  target,
	}
	// Wire everything up!
	outputs := g.Generate(n)
	// Add import, adjust type and name for the generated Load function
	for _, output := range outputs {
		output.Type = g.DataType(output.Import, output.Type)
	}
	// Add an error if we have one
	if len(outputs) > 0 {
		rightmost := outputs[len(outputs)-1]
		if !g.HasError && rightmost.Type == "error" {
			rightmost.Name = "nil"
		}
	}
	// Create the provider
	return &Provider{
		Name:        fnName,
		Target:      target,
		Imports:     g.Imports.List(),
		Externals:   sortExternals(g.Externals),
		Code:        g.Code.String(),
		Results:     outputs,
		externalMap: externalMap(g.Externals),
	}
}

func (node *Node) Print() string {
	out := "digraph G {\n"
	seen := map[string]bool{}
	out += "  " + node.print(seen)
	out += "\n}\n"
	return out
}

func (node *Node) format() string {
	return `'` + node.Import + `'.` + toTypeName(node.Type)
}

func (node *Node) print(seen map[string]bool) string {
	id := node.format()
	if seen[id] {
		return ""
	}
	seen[id] = true
	var outs []string
	for _, dep := range node.Dependencies {
		str := new(strings.Builder)
		label := dep.Type
		fmt.Fprintf(str, `%q -> %q`, dep.format(), id)
		if dep.External {
			label += " (external)"
		}
		fmt.Fprintf(str, ` [label=%q];`, label)
		outs = append(outs, str.String())
		subgraph := dep.print(seen)
		if subgraph == "" {
			continue
		}
		outs = append(outs, subgraph)
	}
	return strings.Join(outs, "\n  ")
}

// Helper function to turn *Web into *web.Web
func toDataType(packageName string, dataType string) string {
	if strings.Contains(dataType, ".") {
		return dataType
	}
	if strings.HasPrefix(dataType, "*") {
		return "*" + packageName + "." + strings.TrimPrefix(dataType, "*")
	}
	return packageName + "." + dataType
}

// Helper function to turn *web.Web into Web
func toTypeName(dataType string) string {
	parts := strings.SplitN(dataType, ".", 2)
	last := parts[len(parts)-1]
	return strings.TrimLeft(last, "[]*")
}

// Sort the variables by name so the order is always consistent.
func sortExternals(externals []*External) []*External {
	sort.Slice(externals, func(i, j int) bool {
		// Hoisted params always come before non-hoisted params
		if externals[i].Hoisted != externals[j].Hoisted {
			return externals[i].Hoisted && !externals[j].Hoisted
		}
		return externals[i].Variable.Name < externals[j].Variable.Name
	})
	return externals
}

// Turn the results into a map for faster provider lookup
func externalMap(exts []*External) map[string]string {
	m := make(map[string]string, len(exts))
	for _, ext := range exts {
		m[ext.Variable.Import+"."+ext.Variable.Type] = ext.Variable.Name
	}
	return m
}
