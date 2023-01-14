package framework

import (
	"strconv"

	"github.com/livebud/bud"
	"github.com/livebud/bud/package/config"
)

func From(c *config.Config) *Flag {
	return &Flag{
		Embed:  c.Embed,
		Minify: c.Minify,
		Hot:    c.Hot,
	}
}

// Flag is used by many of the framework generators
type Flag struct {
	Embed  bool
	Minify bool
	Hot    bool
}

func (f *Flag) Flags() []string {
	return []string{
		"--embed=" + strconv.FormatBool(f.Embed),
		"--minify=" + strconv.FormatBool(f.Minify),
		"--hot=" + strconv.FormatBool(f.Hot),
	}
}

func (f *Flag) Config() *bud.Config {
	return &bud.Config{
		Embed:  f.Embed,
		Minify: f.Minify,
		Hot:    f.Hot,
	}
}
