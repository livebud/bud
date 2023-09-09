package color

import "os"

func Default() Writer {
	if os.Getenv("NO_COLOR") != "" {
		return Ignore()
	}
	return New()
}

var defaultColor = Default()

func Blue(v ...interface{}) string {
	return defaultColor.Blue(v...)
}

func Red(v ...interface{}) string {
	return defaultColor.Red(v...)
}

func Yellow(v ...interface{}) string {
	return defaultColor.Yellow(v...)
}

func Green(v ...interface{}) string {
	return defaultColor.Green(v...)
}

func Pink(v ...interface{}) string {
	return defaultColor.Pink(v...)
}

func Dim(v ...interface{}) string {
	return defaultColor.Dim(v...)
}
