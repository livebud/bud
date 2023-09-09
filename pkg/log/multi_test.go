package log_test

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/livebud/bud/pkg/color"
	"github.com/livebud/bud/pkg/log"
	"github.com/matryer/is"
)

func TestMulti(t *testing.T) {
	is := is.New(t)
	buf1 := new(bytes.Buffer)
	buf2 := new(bytes.Buffer)
	date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	log := log.Multi(
		slog.NewJSONHandler(buf1, &slog.HandlerOptions{
			Level: log.LevelDebug,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					a.Value = slog.TimeValue(date)
				}
				return a
			},
		}),
		log.Filter(log.LevelInfo, &log.Console{Color: color.Ignore(), Writer: buf2}),
	)
	log.Debug("hello", "args", 10)
	log.Info("hello", "planet", "world", "args", 10)
	log.Warn("hello", "planet", "world", "args", 10)
	log.Error("hello world", "planet", "world", "args", 10)
	lines1 := strings.Split(strings.TrimRight(buf1.String(), "\n"), "\n")
	is.Equal(len(lines1), 4)
	is.Equal(string(lines1[0]), `{"time":"2023-01-01T00:00:00Z","level":"DEBUG","msg":"hello","args":10}`)
	is.Equal(string(lines1[1]), `{"time":"2023-01-01T00:00:00Z","level":"INFO","msg":"hello","planet":"world","args":10}`)
	is.Equal(string(lines1[2]), `{"time":"2023-01-01T00:00:00Z","level":"WARN","msg":"hello","planet":"world","args":10}`)
	is.Equal(string(lines1[3]), `{"time":"2023-01-01T00:00:00Z","level":"ERROR","msg":"hello world","planet":"world","args":10}`)
	lines2 := strings.Split(strings.TrimRight(buf2.String(), "\n"), "\n")
	is.Equal(len(lines2), 3)
	is.Equal(string(lines2[0]), "info: hello planet=world args=10")
	is.Equal(string(lines2[1]), "warn: hello planet=world args=10")
	is.Equal(string(lines2[2]), "error: hello world planet=world args=10")
}

func ExampleMulti() {
	log := log.Multi(
		log.Filter(log.LevelInfo, &log.Console{Writer: os.Stderr}),
		slog.NewJSONHandler(os.Stderr, nil),
	)
	log.WithGroup("hello").Debug("world", "args", 10)
	log.Info("hello", "planet", "world", "args", 10)
	log.Warn("hello", "planet", "world", "args", 10)
	log.Error("hello world", "planet", "world", "args", 10)
	// Output:
}
