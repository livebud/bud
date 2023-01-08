package bud

import (
	"context"
	"strconv"
)

type Config struct {
	Embed  bool
	Hot    bool
	Minify bool
}

func (c *Config) Flags() []string {
	return []string{
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
	Packages []string
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

type Runner interface {
	Run(context.Context, *Run) error
}
