package mux

import (
	"github.com/livebud/bud/pkg/mux/ast"
	"github.com/livebud/bud/pkg/mux/internal/parser"
)

// Parse a route
func Parse(route string) (*ast.Route, error) {
	return parser.Parse(route)
}
