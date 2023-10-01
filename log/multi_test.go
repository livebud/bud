package log_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/livebud/bud/internal/color"
	"github.com/livebud/bud/log"
	"github.com/matryer/is"
)

func TestMulti(t *testing.T) {
	is := is.New(t)
	buf1 := new(bytes.Buffer)
	buf2 := new(bytes.Buffer)
	date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now
	log.Now = func() time.Time { return date }
	defer func() { log.Now = now }()
	log := log.Multi(
		log.Json(buf1),
		log.Filter(log.LevelInfo, log.Console(color.Ignore(), buf2)),
	)
	log.Field("args", 10).Debug("hello")
	log.Field("args", 10).Field("planet", "world").Info("hello")
	log.Field("args", 10).Field("planet", "world").Warn("hello")
	log.Field("args", 10).Field("planet", "world").Error("hello world")
	lines1 := strings.Split(strings.TrimRight(buf1.String(), "\n"), "\n")
	is.Equal(len(lines1), 4)
	is.Equal(string(lines1[0]), `{"time":"2023-01-01T00:00:00Z","level":"debug","msg":"hello","fields":{"args":10}}`)
	is.Equal(string(lines1[1]), `{"time":"2023-01-01T00:00:00Z","level":"info","msg":"hello","fields":{"args":10,"planet":"world"}}`)
	is.Equal(string(lines1[2]), `{"time":"2023-01-01T00:00:00Z","level":"warn","msg":"hello","fields":{"args":10,"planet":"world"}}`)
	is.Equal(string(lines1[3]), `{"time":"2023-01-01T00:00:00Z","level":"error","msg":"hello world","fields":{"args":10,"planet":"world"}}`)
	lines2 := strings.Split(strings.TrimRight(buf2.String(), "\n"), "\n")
	is.Equal(len(lines2), 3)
	is.Equal(string(lines2[0]), "info: hello args=10 planet=world")
	is.Equal(string(lines2[1]), "warn: hello args=10 planet=world")
	is.Equal(string(lines2[2]), "error: hello world args=10 planet=world")
}

func ExampleMulti() {
	log := log.Multi(
		log.Filter(log.LevelInfo, log.Console(color.Ignore(), os.Stderr)),
		log.Json(os.Stderr),
	)
	log.Debug("world", "args", 10)
	log.Info("hello", "planet", "world", "args", 10)
	log.Warn("hello", "planet", "world", "args", 10)
	log.Error("hello world", "planet", "world", "args", 10)
	// Output:
}
