package bud

import (
	"strconv"
	"strings"
)

type Flag struct {
	Embed  bool
	Hot    bool
	Minify bool
}

func (f *Flag) Args() []string {
	return []string{
		"--embed=" + strconv.FormatBool(f.Embed),
		"--hot=" + strconv.FormatBool(f.Hot),
		"--minify=" + strconv.FormatBool(f.Minify),
	}
}

func (f *Flag) String() string {
	return strings.Join(f.Args(), " ")
}
