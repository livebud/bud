package bud

import (
	"strconv"
)

type Flag struct {
	Embed  bool
	Hot    bool
	Minify bool
	// Cache  bool
}

func (f *Flag) List() []string {
	args := []string{
		"--embed=" + strconv.FormatBool(f.Embed),
		"--hot=" + strconv.FormatBool(f.Hot),
		"--minify=" + strconv.FormatBool(f.Minify),
		// "--cache=" + strconv.FormatBool(f.Cache),
	}
	return args
}
