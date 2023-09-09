package log_test

import (
	"bytes"
	"testing"

	"github.com/livebud/bud/pkg/color"
	"github.com/livebud/bud/pkg/log"
	"github.com/matryer/is"
)

func TestFilterDebug(t *testing.T) {
	is := is.New(t)
	buf := new(bytes.Buffer)
	log := log.New(log.Filter(log.LevelInfo, &log.Console{
		Color:  color.Ignore(),
		Writer: buf,
	}))
	log.Debug("hello", "args", 10)
	log.Info("hello", "planet", "world", "args", 10)
	log.Warn("hello", "planet", "world", "args", 10)
	log.Error("hello world", "planet", "world", "args", 10)
	lines := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n"), []byte("\n"))
	is.Equal(len(lines), 3)
	is.Equal(string(lines[0]), "info: hello planet=world args=10")
	is.Equal(string(lines[1]), "warn: hello planet=world args=10")
	is.Equal(string(lines[2]), "error: hello world planet=world args=10")
}

func TestFilterError(t *testing.T) {
	is := is.New(t)
	buf := new(bytes.Buffer)
	log := log.New(log.Filter(log.LevelError, &log.Console{
		Color:  color.Ignore(),
		Writer: buf,
	}))
	log.Debug("hello", "args", 10)
	log.Info("hello", "planet", "world", "args", 10)
	log.Warn("hello", "planet", "world", "args", 10)
	log.Error("hello world", "planet", "world", "args", 10)
	lines := bytes.Split(bytes.TrimRight(buf.Bytes(), "\n"), []byte("\n"))
	is.Equal(len(lines), 1)
	is.Equal(string(lines[0]), "error: hello world planet=world args=10")
}
