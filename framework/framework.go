package framework

import (
	"strconv"
)

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
