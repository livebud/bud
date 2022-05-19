package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
)

type Fields []Field

func (f Fields) Len() int {
	return len(f)
}

func (f Fields) Less(i, j int) bool {
	return f[i].Key < f[j].Key
}

func (f Fields) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

type Field struct {
	Key   string
	Value string // Can be empty, though usually a user error
}

type Entry struct {
	Level   Level
	Message string
	Fields  []Field
	Path    string // File path can be empty
}

type Handler interface {
	Log(log Entry)
}

// Flusher is an optional interface
type Flusher interface {
	Flush()
}

type Interface interface {
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

type Option func(logger *logger)

// WithPath determines whether or not to pass the file path to the handler
// This is typically turned off to improve performance, but is really handy
// for debugging.
func WithPath(includePath bool) Option {
	return func(l *logger) {
		l.includePath = includePath
	}
}

// New logger
func New(handler Handler, options ...Option) Interface {
	logger := &logger{
		Handler:     handler,
		includePath: false,
	}
	for _, option := range options {
		option(logger)
	}
	return logger
}

type logger struct {
	Handler     Handler
	fields      []Field
	includePath bool
}

func (l *logger) path() string {
	if !l.includePath {
		return ""
	}
	// Gets the filename. Uses 2 because we're two levels deep from the caller
	_, filename, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	return filepath.Dir(filename)
}

// Turns a list of key values into an array of fields
func (l *logger) keyValues(kvs ...interface{}) (list Fields) {
	size := len(kvs)
	// Special cases
	if size == 0 {
		return nil
	} else if size == 1 {
		return []Field{{Key: fmt.Sprintf("%s", kvs[0])}}
	}
	for i := 1; i < size; i += 2 {
		list = append(list, Field{
			Key:   fmt.Sprintf("%s", kvs[i-1]),
			Value: fmt.Sprintf("%v", kvs[i]),
		})
	}
	// Add in the fields
	list = append(list, l.fields...)
	// Sort the fields by key
	sort.Sort(list)
	return list
}

// New sub logger
func (l *logger) New(fields ...interface{}) Interface {
	return &logger{
		Handler:     l.Handler,
		includePath: l.includePath,
		fields:      append(l.keyValues(fields...), l.fields...),
	}
}

// Debug message is written to the console
func (l *logger) Debug(message string, fields ...interface{}) {
	l.Handler.Log(Entry{
		Message: message,
		Fields:  l.keyValues(fields...),
		Level:   DebugLevel,
		Path:    l.path(),
	})
}

// Info message is written to the console
func (l *logger) Info(message string, fields ...interface{}) {
	l.Handler.Log(Entry{
		Message: message,
		Fields:  l.keyValues(fields...),
		Level:   InfoLevel,
		Path:    l.path(),
	})
}

// Notice message is written to the console
func (l *logger) Notice(message string, fields ...interface{}) {
	l.Handler.Log(Entry{
		Message: message,
		Fields:  l.keyValues(fields...),
		Level:   NoticeLevel,
		Path:    l.path(),
	})
}

// Warn message is written to the console
func (l *logger) Warn(message string, fields ...interface{}) {
	l.Handler.Log(Entry{
		Message: message,
		Fields:  l.keyValues(fields...),
		Level:   WarnLevel,
		Path:    l.path(),
	})
}

// Error message is written to the console
func (l *logger) Error(message string, fields ...interface{}) {
	l.Handler.Log(Entry{
		Message: message,
		Fields:  l.keyValues(fields...),
		Level:   ErrorLevel,
		Path:    l.path(),
	})
}
