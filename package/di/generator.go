package di

import (
	"strings"

	"github.com/livebud/bud/internal/imports"
)

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
	if node.External || node.Hoist {
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
		Key:      toTypeName(dataType),
		Hoisted:  n.Hoist,
		FullType: g.DataType(importPath, dataType),
		Variable: &Variable{
			Import: importPath,
			Name:   g.Variable(importPath, dataType),
			Type:   dataType,
			// Kind:   0, // Unknown kind
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
//	package web
//	func Load(log log.Log) *Web {}
//	type Web struct {}
//
//	package log
//	type Log interface {
//	  Log(msg string)
//	}
//
//	package console
//	type Console struct {}
//	func (c *Console) Log(msg string)
//
// If we generate into the `genweb` package for `*Web`, the `log` package isn't
// actually referenced.
//
//	package genweb
//	import "web"
//	import "console"
//	func Load() *web.Web {
//	  consoleConsole := &console.Console{}
//	  webWeb := Load(consoleConsole)
//	  return webWeb
//	}
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
