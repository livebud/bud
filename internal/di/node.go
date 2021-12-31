package di

import (
	"fmt"
	"strings"

	"gitlab.com/mnm/bud/internal/imports"
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
	// External types are passed in, not instantiated. This is true if the type
	// matches an external dependency or this dependency is hoisted up.
	External bool

	// Unhoistable marks that a node cannot be hoisted. The only case we found so
	// far is if the node is the original function result.
	unhoistable bool
}

func (n *Node) ID() string {
	if n.Declaration != nil {
		return n.Declaration.ID()
	}
	return getID(n.Import, n.Type)
}

// Build a provider for the target import path
func (n *Node) Generate(fnName, target string) *Provider {
	// Build context
	g := &generator{
		Seen:    map[string][]*Variable{},
		Code:    new(strings.Builder),
		Imports: imports.New(),
		Target:  target,
	}
	// Wire everything up!
	outputs := g.Generate(n)
	// Add import, adjust type and name for the generated Load function
	for _, output := range outputs {
		output.Type = g.DataType(output.Import, output.Type)
	}
	// Add an error if we have one
	rightmost := outputs[len(outputs)-1]
	if !g.HasError && rightmost.Type == "error" {
		rightmost.Name = "nil"
	}
	// Create the provider
	return &Provider{
		Name:      fnName,
		Target:    target,
		Imports:   g.Imports.List(),
		Externals: g.Externals,
		Code:      g.Code.String(),
		Results:   outputs,
	}
}

type generator struct {
	Seen       map[string][]*Variable
	Target     string
	Imports    *imports.Set
	Externals  []*External
	Code       *strings.Builder
	HasContext bool
	HasError   bool
}

func (g *generator) Generate(node *Node, params ...*Variable) []*Variable {
	id := node.ID()
	if outputs, ok := g.Seen[id]; ok {
		return outputs
	}
	var results []*Variable
	if node.External {
		external := g.External(node)
		results = append(results, external.Variable)
		g.Seen[id] = results
		return results
	}
	for _, dep := range node.Dependencies {
		outputs := g.Generate(dep, params...)
		results = append(results, outputs[0])
	}
	outputs := node.Declaration.Generate(g, results)
	g.Seen[id] = outputs
	return outputs
}

// Helper to mark a dependency as external returning a variable to that external
// value
func (g *generator) External(n *Node) *External {
	importPath := n.Import
	dataType := n.Type
	ex := &External{
		Key:     toTypeName(dataType),
		Hoisted: n.Declaration != nil,
		Variable: &Variable{
			Import: importPath,
			Name:   g.Variable(importPath, dataType),
			Type:   g.DataType(importPath, dataType),
		},
	}
	g.Externals = append(g.Externals, ex)
	return ex
}

// Helper to create an identifier variable based on the import and type name.
// This function will also add an import automatically if the importPath doesn't
// match our target path.
func (g *generator) Identifier(importPath, typeName string) string {
	name := strings.TrimLeft(typeName, "*[]")
	if g.Target != importPath {
		pkg := g.Imports.Add(importPath)
		return toDataType(pkg, name)
	}
	return name
}

// Helper to create an data type (e.g. *web.Web) based on the import path and
// data type. This function will also add an import automatically if the
// importPath doesn't match our target path.
func (g *generator) DataType(importPath, dataType string) string {
	if importPath == "" {
		return dataType
	}
	if g.Target != importPath {
		pkg := g.Imports.Add(importPath)
		return toDataType(pkg, dataType)
	}
	return dataType
}

func (g *generator) WriteString(code string) (n int, err error) {
	return g.Code.WriteString(code)
}

// Helper to create an variable name based on the import and data type.
//
// Since we're just generating a variable, the variable may not require an
// import. But we still want to reserve the import name in case we do need to
// import this path in the future.
//
// For example, given the `web`, `log`, and `console` packages:
//
//   package web
//   func Load(log log.Log) *Web {}
//   type Web struct {}
//
//   package log
//   type Log interface {
//     Log(msg string)
//   }
//
//   package console
//   type Console struct {}
//   func (c *Console) Log(msg string)
//
// If we generate into the `genweb` package for `*Web`, the `log` package isn't
// actually referenced.
//
//   package genweb
//   import "web"
//   import "console"
//   func Load() *web.Web {
//     consoleConsole := &console.Console{}
//     webWeb := Load(consoleConsole)
//     return webWeb
//   }
//
func (g *generator) Variable(importPath, typeName string) string {
	if typeName == "error" {
		return "err"
	}
	name := strings.TrimLeft(typeName, "*[]")
	pkg := g.Imports.Reserve(importPath)
	return pkg + name
}

func (g *generator) MarkError(hasError bool) {
	g.HasError = hasError
}

func (node *Node) Print() string {
	out := "digraph G {\n"
	seen := map[string]bool{}
	out += "  " + node.print(seen)
	out += "\n}\n"
	return out
}

func (node *Node) format() string {
	return `"` + node.Import + `".` + toTypeName(node.Type)
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
	return strings.TrimRight(last, "[]*")
}
