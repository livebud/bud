// Package color provides color formatting for messages.
package color

import (
	"fmt"
)

type Writer interface {
	Enabled() bool
	Dim(msg ...interface{}) string
	Black(msg ...interface{}) string
	Red(msg ...interface{}) string
	Green(msg ...interface{}) string
	Yellow(msg ...interface{}) string
	Blue(msg ...interface{}) string
	Magenta(msg ...interface{}) string
	Cyan(msg ...interface{}) string
	White(msg ...interface{}) string
}

const (
	escape  = "\x1b"
	reset   = 0
	dim     = 2
	black   = 30
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 34
	magenta = 35
	cyan    = 36
	white   = 37
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

// Dim a message
func (w *writer) Dim(msg ...interface{}) string {
	return format(dim, fmt.Sprint(msg...))
}

// Black formats a black message
func (w *writer) Black(msg ...interface{}) string {
	return format(black, fmt.Sprint(msg...))
}

// Red formats a red message
func (w *writer) Red(msg ...interface{}) string {
	return format(red, fmt.Sprint(msg...))
}

// Green formats a green message
func (w *writer) Green(msg ...interface{}) string {
	return format(green, fmt.Sprint(msg...))
}

// Yellow formats a yellow message
func (w *writer) Yellow(msg ...interface{}) string {
	return format(yellow, fmt.Sprint(msg...))
}

// Blue formats a blue message
func (w *writer) Blue(msg ...interface{}) string {
	return format(blue, fmt.Sprint(msg...))
}

// Magenta formats a magenta message
func (w *writer) Magenta(msg ...interface{}) string {
	return format(magenta, fmt.Sprint(msg...))
}

// Cyan formats a cyan message
func (w *writer) Cyan(msg ...interface{}) string {
	return format(cyan, fmt.Sprint(msg...))
}

// White formats a white message
func (w *writer) White(msg ...interface{}) string {
	return format(white, fmt.Sprint(msg...))
}
