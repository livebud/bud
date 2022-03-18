package command

import (
	"strconv"
	"strings"
)

type Option func(c *Config)

type Config struct {
	Embed  bool
	Hot    bool
	Minify bool
	Trace  bool
}

func (c *Config) Flags() []string {
	return []string{
		"--embed=" + strconv.FormatBool(c.Embed),
		"--hot=" + strconv.FormatBool(c.Hot),
		"--minify=" + strconv.FormatBool(c.Minify),
		"--trace=" + strconv.FormatBool(c.Trace),
	}
}

func (c *Config) String() string {
	return strings.Join(c.Flags(), " ")
}
