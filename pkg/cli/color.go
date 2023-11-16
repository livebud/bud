package cli

import (
	"os"
	"text/template"
)

var reset = color("\033[0m")
var dim = color("\033[37m")

var colors = template.FuncMap{
	"reset":     reset,
	"bold":      color("\033[1m"),
	"dim":       dim,
	"underline": color("\033[4m"),
	"teal":      color("\033[36m"),
	"blue":      color("\033[34m"),
	"yellow":    color("\033[33m"),
	"red":       color("\033[31m"),
	"green":     color("\033[32m"),
}

var nocolor = os.Getenv("NO_COLOR") != ""

func color(code string) func() string {
	return func() string {
		if nocolor {
			return ""
		}
		return code
	}
}
