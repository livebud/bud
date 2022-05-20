package bud

import "strconv"

type Flag struct {
	Embed  bool
	Minify bool
	Hot    string
}

// Map flags into a map to be generated
func (f *Flag) Map() map[string]string {
	return map[string]string{
		"Embed":  strconv.FormatBool(f.Embed),
		"Minify": strconv.FormatBool(f.Minify),
		"Hot":    strconv.Quote(f.Hot),
	}
}
