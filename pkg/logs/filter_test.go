package logs_test

import (
	"bytes"
	"testing"

	"github.com/livebud/bud/internal/color"
	"github.com/livebud/bud/pkg/logs"
	"github.com/matryer/is"
)

func TestFilterDebug(t *testing.T) {
	is := is.New(t)
	buf := new(bytes.Buffer)
	log := logs.New(logs.Filter(logs.LevelInfo, logs.Console(color.Ignore(), buf)))
	log.Debug("hello", "args", 10)
	log.Field("planet", "world").Field("args", 10).Info("hello")
	log.Field("planet", "world").Field("args", 10).Warn("hello")
	log.Field("planet", "world").Field("args", 10).Error("hello world")
	lines := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n"), []byte("\n"))
	is.Equal(len(lines), 3)
	is.Equal(string(lines[0]), "info: hello args=10 planet=world")
	is.Equal(string(lines[1]), "warn: hello args=10 planet=world")
	is.Equal(string(lines[2]), "error: hello world args=10 planet=world")
}

func TestFilterError(t *testing.T) {
	is := is.New(t)
	buf := new(bytes.Buffer)
	log := logs.New(logs.Filter(logs.LevelError, logs.Console(color.Ignore(), buf)))
	log.Debug("hello", "args", 10)
	log.Field("planet", "world").Field("args", 10).Info("hello")
	log.Field("planet", "world").Field("args", 10).Warn("hello")
	log.Field("planet", "world").Field("args", 10).Error("hello world")
	lines := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n"), []byte("\n"))
	is.Equal(len(lines), 1)
	is.Equal(string(lines[0]), "error: hello world args=10 planet=world")
}
