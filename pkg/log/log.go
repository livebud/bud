package log

import (
	"fmt"
)

type Fields map[string]string

type Entry struct {
	Level   Level
	Message string
	Fields  Fields
}

type Handler interface {
	Log(log Entry)
}

// Flusher is an optional interface
type Flusher interface {
	Flush()
}

type Log interface {
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Notice(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Error(message string, args ...interface{})
}

type dispatcher func(log Entry)

func (fn dispatcher) Log(log Entry) {
	fn(log)
}

var Discard = &logger{
	Handler: dispatcher(func(log Entry) {}),
}

// New logger
func New(handler Handler) Log {
	return &logger{handler}
}

type logger struct {
	Handler Handler
}

// Format the message
func (l *logger) format(message string, args ...interface{}) string {
	if len(args) == 0 {
		return message
	}
	return fmt.Sprintf(message, args...)
}

// Debug message is written to the console
func (l *logger) Debug(message string, args ...interface{}) {
	l.Handler.Log(Entry{
		Message: l.format(message, args...),
		Level:   DebugLevel,
	})
}

// Info message is written to the console
func (l *logger) Info(message string, args ...interface{}) {
	l.Handler.Log(Entry{
		Message: l.format(message, args...),
		Level:   InfoLevel,
	})
}

// Notice message is written to the console
func (l *logger) Notice(message string, args ...interface{}) {
	l.Handler.Log(Entry{
		Message: l.format(message, args...),
		Level:   NoticeLevel,
	})
}

// Warn message is written to the console
func (l *logger) Warn(message string, args ...interface{}) {
	l.Handler.Log(Entry{
		Message: l.format(message, args...),
		Level:   WarnLevel,
	})
}

// Error message is written to the console
func (l *logger) Error(message string, args ...interface{}) {
	l.Handler.Log(Entry{
		Message: l.format(message, args...),
		Level:   ErrorLevel,
	})
}
