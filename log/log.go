// Package log is inspired by apex/log.
package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"
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

func (f Fields) clone() Fields {
	clone := make(Fields, len(f))
	for k, v := range f {
		clone[k] = v
	}
	return clone
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

type Entry struct {
	Time    time.Time `json:"time,omitempty"`
	Level   Level     `json:"level,omitempty"`
	Message string    `json:"msg,omitempty"`
	Fields  Fields    `json:"fields,omitempty"`
}

func (e *Entry) MarshalJSON() ([]byte, error) {
	w := new(bytes.Buffer)
	w.WriteByte('{')
	if !e.Time.IsZero() {
		w.WriteString(`"time":"`)
		w.WriteString(e.Time.Format(time.RFC3339))
		w.WriteString(`",`)
	}
	if e.Level != 0 {
		w.WriteString(`"level":"`)
		w.WriteString(e.Level.String())
		w.WriteString(`",`)
	}
	if e.Message != "" {
		w.WriteString(`"msg":`)
		w.WriteString(strconv.Quote(e.Message))
		w.WriteString(`,`)
	}
	if len(e.Fields) > 0 {
		w.WriteString(`"fields":`)
		fields, err := json.Marshal(e.Fields)
		if err != nil {
			return nil, err
		}
		w.Write(fields)
	}
	w.WriteByte('}')
	return w.Bytes(), nil
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
	return l.log(LevelDebug, args, nil)
}

func (l *Logger) Debugf(msg string, args ...interface{}) error {
	return l.logf(LevelDebug, msg, args, nil)
}

func (l *Logger) Info(args ...interface{}) error {
	return l.log(LevelInfo, args, nil)
}

func (l *Logger) Infof(msg string, args ...interface{}) error {
	return l.logf(LevelInfo, msg, args, nil)
}

func (l *Logger) Warn(args ...interface{}) error {
	return l.log(LevelWarn, args, nil)
}

func (l *Logger) Warnf(msg string, args ...interface{}) error {
	return l.logf(LevelWarn, msg, args, nil)
}

func (l *Logger) Notice(args ...interface{}) error {
	return l.log(LevelNotice, args, nil)
}

func (l *Logger) Noticef(msg string, args ...interface{}) error {
	return l.logf(LevelNotice, msg, args, nil)
}

func (l *Logger) Error(args ...interface{}) error {
	return l.log(LevelError, args, nil)
}

func (l *Logger) Errorf(msg string, args ...interface{}) error {
	return l.logf(LevelError, msg, args, nil)
}

func (l *Logger) log(level Level, args []interface{}, fields Fields) error {
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
		Time:    Now(),
		Level:   level,
		Message: msg.String(),
		Fields:  fields,
	})
}

func (l *Logger) logf(level Level, msg string, args []interface{}, fields Fields) error {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	return l.Handler.Log(&Entry{
		Time:    Now(),
		Level:   level,
		Message: msg,
		Fields:  fields.clone(),
	})
}

type logger interface {
	Log
	log(level Level, args []interface{}, fields Fields) error
	logf(level Level, msg string, args []interface{}, fields Fields) error
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
	return l.parent.log(LevelDebug, args, l.fields)
}

func (l *sublogger) Debugf(msg string, args ...interface{}) error {
	return l.parent.logf(LevelDebug, msg, args, l.fields)
}

func (l *sublogger) Info(args ...interface{}) error {
	return l.parent.log(LevelInfo, args, l.fields)
}

func (l *sublogger) Infof(msg string, args ...interface{}) error {
	return l.parent.logf(LevelInfo, msg, args, l.fields)
}

func (l *sublogger) Warn(args ...interface{}) error {
	return l.parent.log(LevelWarn, args, l.fields)
}

func (l *sublogger) Warnf(msg string, args ...interface{}) error {
	return l.parent.logf(LevelWarn, msg, args, l.fields)
}

func (l *sublogger) Notice(args ...interface{}) error {
	return l.parent.log(LevelNotice, args, l.fields)
}

func (l *sublogger) Noticef(msg string, args ...interface{}) error {
	return l.parent.logf(LevelNotice, msg, args, l.fields)
}

func (l *sublogger) Error(args ...interface{}) error {
	return l.parent.log(LevelError, args, l.fields)
}

func (l *sublogger) Errorf(msg string, args ...interface{}) error {
	return l.parent.logf(LevelError, msg, args, l.fields)
}
