package color

import "fmt"

func Ignore() Writer {
	return &ignore{}
}

type ignore struct{}

var _ Writer = (*ignore)(nil)

func (i ignore) Enabled() bool {
	return false
}

// Dim a message
func (ignore) Dim(msg ...interface{}) string {
	return fmt.Sprint(msg...)
}

// Black formats a black message
func (ignore) Black(msg ...interface{}) string {
	return fmt.Sprint(msg...)
}

// Red formats a red message
func (ignore) Red(msg ...interface{}) string {
	return fmt.Sprint(msg...)
}

// Green formats a green message
func (ignore) Green(msg ...interface{}) string {
	return fmt.Sprint(msg...)
}

// Yellow formats a yellow message
func (ignore) Yellow(msg ...interface{}) string {
	return fmt.Sprint(msg...)
}

// Blue formats a blue message
func (ignore) Blue(msg ...interface{}) string {
	return fmt.Sprint(msg...)
}

// Magenta formats a magenta message
func (ignore) Magenta(msg ...interface{}) string {
	return fmt.Sprint(msg...)
}

// Cyan formats a cyan message
func (ignore) Cyan(msg ...interface{}) string {
	return fmt.Sprint(msg...)
}

// White formats a white message
func (ignore) White(msg ...interface{}) string {
	return fmt.Sprint(msg...)
}
