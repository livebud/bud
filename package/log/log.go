// Package log is inspired by apex/log.
package log

import (
	"fmt"
	"sort"
	"time"

	"github.com/livebud/bud/internal/stacktrace"
)

type Fields map[string]interface{}

func (f Fields) Get(key string) interface{} {
	return f[key]
}

func (f Fields) Keys() []string {
	keys := make([]string, len(f))
	i := 0
	for key := range f {
		keys[i] = key
		i++
	}
	sort.Strings(keys)
	return keys
}

// Now returns the current time.
var Now = time.Now

type Log interface {
	Field(key string, value interface{}) Log
	Fields(fields map[string]interface{}) Log
	Debug(msg string, args ...interface{}) error
	Info(msg string, args ...interface{}) error
	Notice(msg string, args ...interface{}) error
	Warn(msg string, args ...interface{}) error
	Error(msg string, args ...interface{}) error
	Err(err error, msg string, args ...interface{}) error
}

func Error(log Log, err error) Log {
	return log.Field("error", err)
}

type Entry struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Fields    Fields
}

type Handler interface {
	Log(log *Entry) error
}

// New logger
func New(level Level, handler Handler) *Logger {
	return &Logger{level, handler}
}

type Logger struct {
	Level   Level
	Handler Handler
}

var _ Log = (*Logger)(nil)

func (l *Logger) Fields(fields map[string]interface{}) Log {
	return &sublogger{l, fields}
}

func (l *Logger) Field(key string, value interface{}) Log {
	return &sublogger{l, Fields{key: value}}
}

func (l *Logger) Debug(msg string, args ...interface{}) error {
	return l.log(DebugLevel, msg, args, nil)
}

func (l *Logger) Info(msg string, args ...interface{}) error {
	return l.log(InfoLevel, msg, args, nil)
}

func (l *Logger) Warn(msg string, args ...interface{}) error {
	return l.log(WarnLevel, msg, args, nil)
}

func (l *Logger) Notice(msg string, args ...interface{}) error {
	return l.log(NoticeLevel, msg, args, nil)
}

func (l *Logger) Error(msg string, args ...interface{}) error {
	return l.log(ErrorLevel, msg, args, nil)
}

func (l *Logger) Err(err error, msg string, args ...interface{}) error {
	return l.log(ErrorLevel, msg, args, Fields{
		"error":  err.Error(),
		"source": stacktrace.Source(1),
	})
}

func (l *Logger) log(level Level, msg string, args []interface{}, fields map[string]interface{}) error {
	if level < l.Level {
		return nil
	}
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return l.Handler.Log(&Entry{
		Timestamp: Now(),
		Level:     level,
		Message:   msg,
		Fields:    fields,
	})
}

type Field struct {
	Key   string
	Value interface{}
}

type logger interface {
	Log
	log(level Level, msg string, args []interface{}, fields map[string]interface{}) error
}

type sublogger struct {
	logger logger
	fields map[string]interface{}
}

func (l *sublogger) Fields(fields map[string]interface{}) Log {
	for k, v := range l.fields {
		if _, ok := fields[k]; !ok {
			fields[k] = v
		}
	}
	return &sublogger{l.logger, fields}
}

func (l *sublogger) Field(key string, value interface{}) Log {
	return l.Fields(Fields{key: value})
}

func (l *sublogger) Debug(msg string, args ...interface{}) error {
	return l.log(DebugLevel, msg, args, l.fields)
}

func (l *sublogger) Info(msg string, args ...interface{}) error {
	return l.log(InfoLevel, msg, args, l.fields)
}

func (l *sublogger) Warn(msg string, args ...interface{}) error {
	return l.log(WarnLevel, msg, args, l.fields)
}

func (l *sublogger) Notice(msg string, args ...interface{}) error {
	return l.log(NoticeLevel, msg, args, l.fields)
}

func (l *sublogger) Error(msg string, args ...interface{}) error {
	return l.log(ErrorLevel, msg, args, l.fields)
}

func (l *sublogger) Err(err error, msg string, args ...interface{}) error {
	return l.Fields(Fields{
		"error":  err.Error(),
		"source": stacktrace.Source(1),
	}).Error(msg, args...)
}

func (l *sublogger) log(level Level, msg string, args []interface{}, fields map[string]interface{}) error {
	return l.logger.log(level, msg, args, fields)
}
