//go:build tools
// +build tools

// Tools we depend on. This file is here to prevent `go mod tidy` from cleaning
// up these dependencies
package bud

import (
	_ "github.com/evanw/esbuild/cmd/esbuild"
	_ "github.com/livebud/bud-test-plugin"
	_ "github.com/pointlander/peg"
	_ "honnef.co/go/tools/cmd/staticcheck"
	_ "src.techknowlogick.com/xgo"
)
