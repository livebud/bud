package bud

import (
	"context"
	"io"
	"strconv"

	"github.com/livebud/bud/internal/pubsub"
	"github.com/livebud/bud/package/socket"
)

type Config struct {
	Dir string
	Log string

	Embed  bool
	Hot    bool
	Minify bool

	ListenWeb string
	ListenDev string
	ListenAFS string

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string

	// Used for testing
	Bus         pubsub.Client
	WebListener socket.Listener
	DevListener socket.Listener
	AFSListener socket.Listener
}

func (c *Config) Flags() []string {
	return []string{
		"--log=" + strconv.Quote(c.Log),
		"--embed=" + strconv.FormatBool(c.Embed),
		"--minify=" + strconv.FormatBool(c.Minify),
		"--hot=" + strconv.FormatBool(c.Hot),
	}
}

type Create struct {
	Dir string
}

type Creator interface {
	Create(context.Context, *Create) error
}

type Generate struct {
	Config
	DevAddress string
	Packages   []string
}

type Generator interface {
	Generate(context.Context, *Generate) error
}

type Build struct {
	Config
}

type Builder interface {
	Build(context.Context, *Build) error
}

type Run struct {
	Config
	Watch      bool
	WebAddress string
	DevAddress string
}

type Run2 struct {
	Watch bool
}

type Runner interface {
	Run(context.Context, *Run) error
}
