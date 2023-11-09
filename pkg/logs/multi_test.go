package logs_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/livebud/bud/internal/color"
	"github.com/livebud/bud/pkg/logs"
	"github.com/matryer/is"
)

func TestMulti(t *testing.T) {
	is := is.New(t)
	buf1 := new(bytes.Buffer)
	buf2 := new(bytes.Buffer)
	date := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	now := time.Now
	logs.Now = func() time.Time { return date }
	defer func() { logs.Now = now }()
	log := logs.Multi(
		logs.Json(buf1),
		logs.Filter(logs.LevelInfo, logs.Console(color.Ignore(), buf2)),
	)
	log.Field("args", 10).Debug("hello")
	log.Field("args", 10).Field("planet", "world").Info("hello")
	log.Field("args", 10).Field("planet", "world").Warn("hello")
	log.Field("args", 10).Field("planet", "world").Error("hello world")
	lines1 := strings.Split(strings.TrimRight(buf1.String(), "\n"), "\n")
	is.Equal(len(lines1), 4)
	is.Equal(string(lines1[0]), `{"ts":"2023-01-01T00:00:00Z","lvl":"debug","msg":"hello","fields":{"args":10}}`)
	is.Equal(string(lines1[1]), `{"ts":"2023-01-01T00:00:00Z","lvl":"info","msg":"hello","fields":{"args":10,"planet":"world"}}`)
	is.Equal(string(lines1[2]), `{"ts":"2023-01-01T00:00:00Z","lvl":"warn","msg":"hello","fields":{"args":10,"planet":"world"}}`)
	is.Equal(string(lines1[3]), `{"ts":"2023-01-01T00:00:00Z","lvl":"error","msg":"hello world","fields":{"args":10,"planet":"world"}}`)
	lines2 := strings.Split(strings.TrimRight(buf2.String(), "\n"), "\n")
	is.Equal(len(lines2), 3)
	is.Equal(string(lines2[0]), "info: hello args=10 planet=world")
	is.Equal(string(lines2[1]), "warn: hello args=10 planet=world")
	is.Equal(string(lines2[2]), "error: hello world args=10 planet=world")
}

func ExampleMulti() {
	log := logs.Multi(
		logs.Filter(logs.LevelDebug, logs.Console(color.Default(), os.Stderr)),
		logs.Json(os.Stderr),
	)
	log.Debug("world", "args", 10)
	log.Field("planet", "world").Field("args", 10).Info("hello")
	log.Warnf("hello %s", "world")
	log.Error("hello world", "planet", "world", "args", 10)
	// Output:
}
