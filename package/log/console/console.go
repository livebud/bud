package console

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/go-logfmt/logfmt"

	"github.com/livebud/bud/internal/ansi"
	log "github.com/livebud/bud/package/log"
)

func paint(level log.Level, msg string) string {
	switch level {
	case log.DebugLevel:
		return ansi.White(msg)
	case log.InfoLevel:
		return ansi.Blue(msg)
	case log.NoticeLevel:
		return ansi.Pink(msg)
	case log.WarnLevel:
		return ansi.Yellow(msg)
	case log.ErrorLevel:
		return ansi.Red(msg)
	default:
		return ""
	}
}

// prefixes
var prefixes = func() [6]string {
	if os.Getenv("NO_COLOR") != "" {
		return [6]string{
			log.DebugLevel:  "debug:",
			log.InfoLevel:   "info:",
			log.NoticeLevel: "notice:",
			log.WarnLevel:   "warn:",
			log.ErrorLevel:  "error:",
		}
	}
	return [6]string{
		log.DebugLevel:  "|",
		log.InfoLevel:   "|",
		log.NoticeLevel: "|",
		log.WarnLevel:   "|",
		log.ErrorLevel:  "|",
	}
}()

// New console handler
func New(w io.Writer) log.Handler {
	return &console{w: w}
}

// console logger
type console struct {
	mu sync.Mutex
	w  io.Writer
}

// Log implements Logger
func (c *console) Log(log *log.Entry) error {
	// Format the message
	msg := new(strings.Builder)
	msg.WriteString(paint(log.Level, prefixes[log.Level]) + " " + log.Message)

	// Format and log the fields
	if len(log.Fields) > 0 {
		keys := log.Fields.Keys()
		fields := new(strings.Builder)
		enc := logfmt.NewEncoder(fields)
		for _, key := range keys {
			enc.EncodeKeyval(key, log.Fields.Get(key))
		}
		enc.Reset()
		msg.WriteString(" " + ansi.Dim(fields.String()))
	}
	msg.WriteString("\n")

	// Write out
	c.mu.Lock()
	fmt.Fprint(c.w, msg.String())
	c.mu.Unlock()

	return nil
}

// Stderr is a console log singleton that writes to stderr
var stderr = log.New(New(os.Stderr))

var (
	// Return a logger with a field
	Field = stderr.Field
	// Debug message is written to the console
	Debug = stderr.Debug
	// Debugf message is written to the console
	Debugf = stderr.Debugf
	// Info message is written to the console
	Info = stderr.Info
	// Infof message is written to the console
	Infof = stderr.Infof
	// Notice message is written to the console
	Notice = stderr.Notice
	// Noticef message is written to the console
	Noticef = stderr.Noticef
	// Warn message is written to the console
	Warn = stderr.Warn
	// Warnf message is written to the console
	Warnf = stderr.Warnf
	// Error message is written to the console
	Error = stderr.Error
	// Errorf message is written to the console
	Errorf = stderr.Errorf
)
