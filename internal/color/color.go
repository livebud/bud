package color

import (
	"fmt"
)

type Writer interface {
	Enabled() bool
	Blue(v ...interface{}) string
	Red(v ...interface{}) string
	Yellow(v ...interface{}) string
	Green(v ...interface{}) string
	Pink(v ...interface{}) string
	Dim(v ...interface{}) string
}

const (
	escape = "\x1b"
	reset  = 0
	dim    = 2
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	pink   = 35
)

func New() Writer {
	return &writer{}
}

// Format a message with a color
func format(color int, msg string) string {
	return fmt.Sprintf("%s[%dm%s%s[%dm", escape, color, msg, escape, reset)
}

type writer struct{}

var _ Writer = (*writer)(nil)

func (w *writer) Enabled() bool {
	return true
}

func (w *writer) Blue(v ...interface{}) string {
	return format(blue, fmt.Sprint(v...))
}

func (w *writer) Red(v ...interface{}) string {
	return format(red, fmt.Sprint(v...))
}

func (w *writer) Yellow(v ...interface{}) string {
	return format(yellow, fmt.Sprint(v...))
}

func (w *writer) Green(v ...interface{}) string {
	return format(green, fmt.Sprint(v...))
}

func (w *writer) Pink(v ...interface{}) string {
	return format(pink, fmt.Sprint(v...))
}

func (w *writer) Dim(v ...interface{}) string {
	return format(dim, fmt.Sprint(v...))
}
