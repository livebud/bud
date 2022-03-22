package bud

import (
	"strconv"
)

type Flag struct {
	Embed  bool
	Hot    bool
	Minify bool
	Cache  bool
}

func (f *Flag) List(cachePath string) []string {
	args := []string{
		"--embed=" + strconv.FormatBool(f.Embed),
		"--hot=" + strconv.FormatBool(f.Hot),
		"--minify=" + strconv.FormatBool(f.Minify),
	}
	// Add the cache path if the cache is enabled
	if f.Cache {
		args = append(args,
			"--cache="+cachePath,
		)
	}
	return args
}
