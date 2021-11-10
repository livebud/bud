package di

import (
	"strings"

	"gitlab.com/mnm/bud/internal/imports"
)

type GenerateInput struct {
	// Targets to load
	Dependencies []*Dependency `json:"dependencies,omitempty"`
	// External types
	Externals []*Dependency `json:"externals,omitempty"`
	// Hoist dependencies that don't depend on externals, turning them into
	// externals. This is to avoid initializing these inner deps every time.
	// Useful for per-request dependency injection.
	Hoist bool `json:"hoist,omitempty"`
	// Target import path of the generated code
	Target string `json:"target,omitempty"`
}

// Generate a provider
func (i *Injector) Generate(in *GenerateInput) (*Provider, error) {
	node, err := i.Load(&LoadInput{
		Dependencies: in.Dependencies,
		Externals:    in.Externals,
		Hoist:        in.Hoist,
	})
	if err != nil {
		return nil, err
	}
	provider := node.Generate(in.Target)
	return provider, nil
}

// Generator is the context needed during Generate
type Generator struct {
	Seen       map[string][]*Variable
	Target     string
	Imports    *imports.Set
	Externals  []*External
	Code       *strings.Builder
	HasContext bool
	HasError   bool
}

// Generate a provider recursively. Types are qualified and imports are added
// based on the target import path.
func (n *Node) Generate(target string) *Provider {
	// Generator context
	gen := &Generator{
		Seen:    map[string][]*Variable{},
		Code:    new(strings.Builder),
		Imports: imports.New(),
		Target:  target,
	}
	// Wire everything up!
	outputs := n.generate(gen)
	// Add import, adjust type and name for the generated Load function
	for _, output := range outputs {
		output.Type = gen.DataType(output.Import, output.Type)
	}
	// Add an error if we have one
	if gen.HasError {
		outputs = append(outputs, &Variable{
			Import: "", // error doesn't have an import
			Name:   "err",
			Type:   "error",
		})
	}
	// Create the provider
	return &Provider{
		Target:    target,
		Imports:   gen.Imports.List(),
		Externals: gen.Externals,
		Code:      gen.Code.String(),
		Results:   outputs,
	}
}

// Helper to mark a dependency as external returning a variable to that external
// value
func (g *Generator) External(n *Node) *External {
	importPath := n.Original.Import
	dataType := n.Original.Type
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
func (g *Generator) Identifier(importPath, name string) string {
	if g.Target != name {
		pkg := g.Imports.Add(importPath)
		return toDataType(pkg, name)
	}
	return name
}

// Helper to create an data type (e.g. *web.Web) based on the import path and
// data type. This function will also add an import automatically if the
// importPath doesn't match our target path.
func (g *Generator) DataType(importPath, dataType string) string {
	if g.Target != importPath {
		pkg := g.Imports.Add(importPath)
		return toDataType(pkg, dataType)
	}
	return dataType
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
func (g *Generator) Variable(importPath, dataType string) string {
	if dataType == "error" {
		return "err"
	}
	name := strings.TrimLeft(dataType, "*[]")
	pkg := g.Imports.Reserve(importPath)
	return pkg + name
}

// Generate recursively. This function starts at the top with nodes that don't
// have any dependencies. It generates those, then uses the results to as inputs
// to downstream nodes. It does this recursively until we've generated the full
// graph.
func (n *Node) generate(gen *Generator, params ...*Variable) []*Variable {
	if outputs, ok := gen.Seen[n.ID()]; ok {
		return outputs
	}
	var results []*Variable
	if n.External {
		external := gen.External(n)
		results = append(results, external.Variable)
		gen.Seen[n.ID()] = results
		return results
	}
	for _, dep := range n.Dependencies {
		outputs := dep.generate(gen, params...)
		results = append(results, outputs[0])
	}
	outputs := n.Declaration.Generate(gen, results)
	gen.Seen[n.ID()] = outputs
	return outputs
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
	return strings.TrimPrefix(last, "*")
}
