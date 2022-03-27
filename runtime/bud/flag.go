package bud

import "strconv"

type Flag struct {
	Embed  bool
	Hot    bool
	Minify bool
}

// Map flags into a map to be generated
func (f *Flag) Map() map[string]string {
	return map[string]string{
		"Embed":  strconv.FormatBool(f.Embed),
		"Hot":    strconv.FormatBool(f.Hot),
		"Minify": strconv.FormatBool(f.Minify),
	}
}
