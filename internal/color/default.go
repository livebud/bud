package color

import "os"

func Default() Writer {
	if os.Getenv("NO_COLOR") != "" {
		return Ignore()
	}
	return New()
}

var defaultColor = Default()

// Dim a message
func Dim(msg ...interface{}) string {
	return defaultColor.Dim(msg...)
}

// Black formats a black message
func Black(msg ...interface{}) string {
	return defaultColor.Black(msg...)
}

// Red formats a red message
func Red(msg ...interface{}) string {
	return defaultColor.Red(msg...)
}

// Green formats a green message
func Green(msg ...interface{}) string {
	return defaultColor.Green(msg...)
}

// Yellow formats a yellow message
func Yellow(msg ...interface{}) string {
	return defaultColor.Yellow(msg...)
}

// Blue formats a blue message
func Blue(msg ...interface{}) string {
	return defaultColor.Blue(msg...)
}

// Magenta formats a magenta message
func Magenta(msg ...interface{}) string {
	return defaultColor.Magenta(msg...)
}

// Cyan formats a cyan message
func Cyan(msg ...interface{}) string {
	return defaultColor.Cyan(msg...)
}

// White formats a white message
func White(msg ...interface{}) string {
	return defaultColor.White(msg...)
}
