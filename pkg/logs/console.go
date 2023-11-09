package logs

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/go-logfmt/logfmt"
	"github.com/livebud/bud/internal/color"
)

// Debugger log
func Debugger() *Logger {
	return New(Console(color.Default(), os.Stderr))
}

// Default log
func Default() *Logger {
	return New(Filter(LevelInfo, Console(color.Default(), os.Stderr)))
}

// Parse the logger with a given filter
func Parse(filter string) (*Logger, error) {
	level, err := ParseLevel(filter)
	if err != nil {
		return nil, err
	}
	return New(Filter(level, Console(color.Default(), os.Stderr))), nil
}

// Console handler
func Console(color color.Writer, writer io.Writer) Handler {
	return &console{color, sync.Mutex{}, writer, prefixes(color)}
}

// console logger
type console struct {
	color    color.Writer
	mu       sync.Mutex
	writer   io.Writer
	prefixes map[Level]string
}

// Prefixes
func prefixes(color color.Writer) map[Level]string {
	if color.Enabled() {
		return map[Level]string{
			LevelDebug: "|",
			LevelInfo:  "|",
			LevelWarn:  "|",
			LevelError: "|",
		}
	}
	return map[Level]string{
		LevelDebug: "debug:",
		LevelInfo:  "info:",
		LevelWarn:  "warn:",
		LevelError: "error:",
	}
}

func (c *console) format(level Level, msg string) string {
	switch level {
	case LevelDebug:
		return c.color.Dim(msg)
	case LevelInfo:
		return c.color.Blue(msg)
	case LevelNotice:
		return c.color.Magenta(msg)
	case LevelWarn:
		return c.color.Yellow(msg)
	case LevelError:
		return c.color.Red(msg)
	default:
		return ""
	}
}

// Log implements Logger
func (c *console) Log(entry *Entry) error {
	// Format the message
	msg := new(strings.Builder)
	msg.WriteString(c.format(entry.Level, c.prefixes[entry.Level]) + " " + entry.Message)

	// Format and log the fields
	if len(entry.Fields) > 0 {
		keys := entry.Fields.Keys()
		fields := new(strings.Builder)
		enc := logfmt.NewEncoder(fields)
		for _, key := range keys {
			enc.EncodeKeyval(key, entry.Fields.Get(key))
		}
		enc.Reset()
		msg.WriteString(" " + c.color.Dim(fields.String()))
	}
	msg.WriteString("\n")

	// Write out
	c.mu.Lock()
	fmt.Fprint(c.writer, msg.String())
	c.mu.Unlock()

	return nil
}

// Stderr is a console log singleton that writes to stderr
var stderr = Default()

var (
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
