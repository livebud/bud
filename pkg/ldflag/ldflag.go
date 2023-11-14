package ldflag

import "strconv"

var (
	live   = "/live.js"
	embed  = ""
	minify = ""
)

func Live() string {
	return live
}

func Embed() string {
	return embed
}

func Minify() bool {
	if b, err := strconv.ParseBool(minify); err != nil {
		return false
	} else {
		return b
	}
}
