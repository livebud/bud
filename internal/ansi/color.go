package ansi

import "os"

type Colors struct {
	Reset     string
	Bold      string
	Dim       string
	Underline string

	White  string
	Teal   string
	Blue   string
	Yellow string
	Red    string
	Green  string
}

// Color set
var Color = func() Colors {
	if os.Getenv("NO_COLOR") != "" {
		return Colors{}
	}
	return Colors{
		Reset:     "\033[0m",
		Bold:      "\033[1m",
		Dim:       "\033[37m",
		Underline: "\033[4m",

		White:  "\033[37m",
		Teal:   "\033[36m",
		Blue:   "\033[34m",
		Yellow: "\033[33m",
		Red:    "\033[31m",
		Green:  "\033[32m",
	}
}()
