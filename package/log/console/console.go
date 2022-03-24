package console

import (
	"fmt"
	"io"
	"os"
	"sync"

	"gitlab.com/mnm/bud/internal/ansi"
	"gitlab.com/mnm/bud/package/log"
)

func color(level log.Level) string {
	switch level {
	case log.DebugLevel:
		return ansi.Color.White
	case log.InfoLevel:
		return ansi.Color.Blue
	case log.WarnLevel:
		return ansi.Color.Yellow
	case log.ErrorLevel:
		return ansi.Color.Red
	default:
		return ""
	}
}

// prefixes
var prefixes = func() [6]string {
	if os.Getenv("NO_COLOR") != "" {
		return [6]string{
			log.DebugLevel: "debug:",
			log.InfoLevel:  "info:",
			log.WarnLevel:  "warn:",
			log.ErrorLevel: "error:",
		}
	}
	return [6]string{
		log.DebugLevel: "|",
		log.InfoLevel:  "|",
		log.WarnLevel:  "|",
		log.ErrorLevel: "|",
	}
}()

// New console handler
func New(w io.Writer) log.Handler {
	return &console{Writer: w}
}

// console logger
type console struct {
	mu     sync.Mutex
	Writer io.Writer
}

// Log implements Logger
func (c *console) Log(log log.Entry) {
	// Load the level prefix and color
	level := prefixes[log.Level]
	color := color(log.Level)

	// Format the message
	msg := color + level + ansi.Color.Reset + " " + log.Message
	for _, field := range log.Fields {
		msg += ansi.Color.Dim + " " + field.Key + "=" + field.Value + ansi.Color.Reset
	}
	msg += "\n"

	// Write out
	c.mu.Lock()
	fmt.Fprint(c.Writer, msg)
	c.mu.Unlock()
}

// Stderr is a console log singleton that writes to stderr
var Stderr = log.New(New(os.Stderr))

// Debug message is written to the console
func Debug(message string, args ...interface{}) {
	Stderr.Debug(message, args...)
}

// Info message is written to the console
func Info(message string, args ...interface{}) {
	Stderr.Info(message, args...)
}

// Notice message is written to the console
func Notice(message string, args ...interface{}) {
	Stderr.Notice(message, args...)
}

// Warn message is written to the console
func Warn(message string, args ...interface{}) {
	Stderr.Warn(message, args...)
}

// Error message is written to the console
func Error(message string, args ...interface{}) {
	Stderr.Error(message, args...)
}
