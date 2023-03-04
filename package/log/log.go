// Package log is inspired by apex/log.
package log

import (
	"bytes"
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
	Debug(args ...interface{}) error
	Debugf(msg string, args ...interface{}) error
	Info(args ...interface{}) error
	Infof(msg string, args ...interface{}) error
	Notice(args ...interface{}) error
	Noticef(msg string, args ...interface{}) error
	Warn(args ...interface{}) error
	Warnf(msg string, args ...interface{}) error
	Error(args ...interface{}) error
	Errorf(msg string, args ...interface{}) error
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
func New(handler Handler) *Logger {
	return &Logger{handler}
}

type Logger struct {
	Handler Handler
}

var _ Log = (*Logger)(nil)

func (l *Logger) Fields(fields map[string]interface{}) Log {
	return &sublogger{l, fields}
}

func (l *Logger) Field(key string, value interface{}) Log {
	return &sublogger{l, Fields{key: value}}
}

func (l *Logger) Debug(args ...interface{}) error {
	return l.log(DebugLevel, args, nil)
}

func (l *Logger) Debugf(msg string, args ...interface{}) error {
	return l.logf(DebugLevel, msg, args, nil)
}

func (l *Logger) Info(args ...interface{}) error {
	return l.log(InfoLevel, args, nil)
}

func (l *Logger) Infof(msg string, args ...interface{}) error {
	return l.logf(InfoLevel, msg, args, nil)
}

func (l *Logger) Warn(args ...interface{}) error {
	return l.log(WarnLevel, args, nil)
}

func (l *Logger) Warnf(msg string, args ...interface{}) error {
	return l.logf(WarnLevel, msg, args, nil)
}

func (l *Logger) Notice(args ...interface{}) error {
	return l.log(NoticeLevel, args, nil)
}

func (l *Logger) Noticef(msg string, args ...interface{}) error {
	return l.logf(NoticeLevel, msg, args, nil)
}

func (l *Logger) Error(args ...interface{}) error {
	return l.log(ErrorLevel, args, map[string]interface{}{
		"source": stacktrace.Source(1),
	})
}

func (l *Logger) Errorf(msg string, args ...interface{}) error {
	return l.logf(ErrorLevel, msg, args, map[string]interface{}{
		"source": stacktrace.Source(1),
	})
}

func (l *Logger) log(level Level, args []interface{}, fields map[string]interface{}) error {
	if len(args) == 0 {
		return nil
	}
	var msg bytes.Buffer
	// Add spaces between the arguments
	for argNum, arg := range args {
		if argNum > 0 {
			msg.WriteByte(' ')
		}
		msg.WriteString(fmt.Sprint(arg))
	}
	return l.Handler.Log(&Entry{
		Timestamp: Now(),
		Level:     level,
		Message:   msg.String(),
		Fields:    fields,
	})
}

func (l *Logger) logf(level Level, msg string, args []interface{}, fields map[string]interface{}) error {
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
	log(level Level, args []interface{}, fields map[string]interface{}) error
	logf(level Level, msg string, args []interface{}, fields map[string]interface{}) error
}

type sublogger struct {
	parent logger
	fields map[string]interface{}
}

func (l *sublogger) Fields(fields map[string]interface{}) Log {
	for k, v := range l.fields {
		if _, ok := fields[k]; !ok {
			fields[k] = v
		}
	}
	return &sublogger{l.parent, fields}
}

func (l *sublogger) Field(key string, value interface{}) Log {
	return l.Fields(Fields{key: value})
}

func (l *sublogger) Debug(args ...interface{}) error {
	return l.parent.log(DebugLevel, args, l.fields)
}

func (l *sublogger) Debugf(msg string, args ...interface{}) error {
	return l.parent.logf(DebugLevel, msg, args, l.fields)
}

func (l *sublogger) Info(args ...interface{}) error {
	return l.parent.log(InfoLevel, args, l.fields)
}

func (l *sublogger) Infof(msg string, args ...interface{}) error {
	return l.parent.logf(InfoLevel, msg, args, l.fields)
}

func (l *sublogger) Warn(args ...interface{}) error {
	return l.parent.log(WarnLevel, args, l.fields)
}

func (l *sublogger) Warnf(msg string, args ...interface{}) error {
	return l.parent.logf(WarnLevel, msg, args, l.fields)
}

func (l *sublogger) Notice(args ...interface{}) error {
	return l.parent.log(NoticeLevel, args, l.fields)
}

func (l *sublogger) Noticef(msg string, args ...interface{}) error {
	return l.parent.logf(NoticeLevel, msg, args, l.fields)
}

func (l *sublogger) Error(args ...interface{}) error {
	if _, ok := l.fields["source"]; !ok {
		l.fields["source"] = stacktrace.Source(1)
	}
	return l.parent.log(ErrorLevel, args, l.fields)
}

func (l *sublogger) Errorf(msg string, args ...interface{}) error {
	if _, ok := l.fields["source"]; !ok {
		l.fields["source"] = stacktrace.Source(1)
	}
	return l.parent.logf(ErrorLevel, msg, args, l.fields)
}
