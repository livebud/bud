package framework

import (
	"io"

	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/socket"
)

// Flag is used by many of the framework generators
type Flag struct {
	Embed  bool
	Minify bool
	Hot    bool

	// Comes from *bud.Input
	// TODO: remove
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string
	// Currently passed in only for testing
	BudLn socket.Listener // Can be nil
	WebLn socket.Listener // Can be nil
	Bus   pubsub.Client   // Can be nil
}
