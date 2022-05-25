package bud

import (
	"github.com/livebud/bud/runtime/bud"
)

type Flag = bud.Flag

// type Flag struct {
// 	Hot    bool
// 	Minify bool
// 	Embed  bool
// }

// // Map flags into a map to be generated
// func (f *Flag) Map() map[string]string {
// 	return map[string]string{
// 		"Embed":  strconv.FormatBool(f.Embed),
// 		"Hot":    strconv.FormatBool(f.Hot),
// 		"Minify": strconv.FormatBool(f.Minify),
// 	}
// }
