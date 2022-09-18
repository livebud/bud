package framework

import "io"

// Flag is used by many of the framework generators
type Flag struct {
	Embed  bool
	Minify bool
	Hot    bool

	// Comes from *bud.Input
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string
}
