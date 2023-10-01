package mux

import (
	"github.com/livebud/bud/mux/ast"
	"github.com/livebud/bud/mux/internal/parser"
)

// Parse a route
func Parse(route string) (*ast.Route, error) {
	return parser.Parse(route)
}
