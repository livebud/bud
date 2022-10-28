package ansi

import (
	"os"

	"github.com/aybabtme/rgbterm"
)

var noColor = os.Getenv("NO_COLOR") != ""

func paint(msg string, r, g, b uint8) string {
	if noColor {
		return msg
	}
	return rgbterm.FgString(msg, r, g, b)
}

func Dim(msg string) string {
	if noColor {
		return msg
	}
	return "\033[37m" + msg + "\033[0m"
}

func Bold(msg string) string {
	if noColor {
		return msg
	}
	return "\033[1m" + msg + "\033[0m"
}

func White(msg string) string {
	return paint(msg, 226, 232, 240)
}

func Green(msg string) string {
	return paint(msg, 43, 255, 99)
}

func Blue(msg string) string {
	return paint(msg, 43, 199, 255)
}

func Yellow(msg string) string {
	return paint(msg, 255, 237, 43)
}

func Pink(msg string) string {
	return paint(msg, 192, 38, 211)
}

func Red(msg string) string {
	return paint(msg, 255, 43, 43)
}
