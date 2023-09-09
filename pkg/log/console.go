package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"runtime"
	"strings"
	"sync"

	"github.com/go-logfmt/logfmt"
	"github.com/livebud/bud/pkg/color"
)

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

// Console handler for printing logs to the terminal
type Console struct {
	Color     color.Writer
	Writer    io.Writer
	AddSource bool
	prefixes  map[Level]string
	mu        sync.Mutex // mu protects the writer
	once      sync.Once  // Called once to setup defaults
	attrs     []slog.Attr
	groups    []string
}

var _ Handler = (*Console)(nil)

// Setup defaults
func (c *Console) setup() {
	if c.Writer == nil {
		c.Writer = os.Stderr
	}
	if c.Color == nil {
		c.Color = color.Default()
	}
	c.prefixes = prefixes(c.Color)
}

// Enabled is always set to true. Use log.Filter to filter out log levels
func (c *Console) Enabled(context.Context, Level) bool {
	return true
}

func (c *Console) Handle(ctx context.Context, record slog.Record) error {
	c.once.Do(c.setup)
	// Format the message
	msg := new(strings.Builder)
	msg.WriteString(c.format(record.Level, c.prefixes[record.Level]) + " " + record.Message)
	// Format and log the fields
	fields := new(strings.Builder)
	enc := logfmt.NewEncoder(fields)
	if record.NumAttrs() > 0 {
		prefix := strings.Join(c.groups, ".")
		record.Attrs(func(attr slog.Attr) bool {
			if prefix != "" {
				attr.Key = prefix + "." + attr.Key
			}
			enc.EncodeKeyval(attr.Key, attr.Value.String())
			return true
		})
	}
	if c.AddSource {
		enc.EncodeKeyval("source", c.source(record.PC))
	}
	enc.Reset()
	msg.WriteString(" " + c.Color.Dim(fields.String()))
	msg.WriteString("\n")

	// Write out
	c.mu.Lock()
	fmt.Fprint(c.Writer, msg.String())
	c.mu.Unlock()

	return nil
}

func (c *Console) source(pc uintptr) string {
	fs := runtime.CallersFrames([]uintptr{pc})
	f, _ := fs.Next()
	return fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
}

func (c *Console) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Console{
		Color:     c.Color,
		Writer:    c.Writer,
		AddSource: c.AddSource,
		groups:    c.groups,
		attrs:     append(append([]slog.Attr{}, c.attrs...), attrs...),
	}
}

func (c *Console) WithGroup(group string) slog.Handler {
	return &Console{
		Color:     c.Color,
		Writer:    c.Writer,
		AddSource: c.AddSource,
		attrs:     c.attrs,
		groups:    append(append([]string{}, c.groups...), group),
	}
}

func (c *Console) format(level Level, msg string) string {
	switch level {
	case LevelDebug:
		return c.Color.Dim(msg)
	case LevelInfo:
		return c.Color.Blue(msg)
	case LevelWarn:
		return c.Color.Yellow(msg)
	case LevelError:
		return c.Color.Red(msg)
	default:
		return ""
	}
}
