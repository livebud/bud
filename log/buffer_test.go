package log_test

import (
	"errors"
	"testing"

	"github.com/livebud/bud/log"
	"github.com/matryer/is"
)

func TestLog(t *testing.T) {
	is := is.New(t)
	handler := log.Buffer()
	log := log.New(handler)
	log.Info("hello")
	log.Info("world")
	log.Warnf("hello %s!", "mars")
	log.Fields(map[string]interface{}{
		"file":   "memory_test.go",
		"detail": "file not exist",
	}).Field("one", "two").Errorf("%d. oh noz", 1)
	is.Equal(len(handler.Entries), 4)
	is.Equal(handler.Entries[0].Level.String(), "info")
	is.Equal(handler.Entries[0].Message, "hello")
	is.Equal(handler.Entries[1].Level.String(), "info")
	is.Equal(handler.Entries[1].Message, "world")
	is.Equal(handler.Entries[2].Level.String(), "warn")
	is.Equal(handler.Entries[2].Message, "hello mars!")
	is.Equal(handler.Entries[3].Level.String(), "error")
	is.Equal(handler.Entries[3].Message, "1. oh noz")
	fields := handler.Entries[3].Fields
	is.Equal(len(fields), 3)
	is.Equal(fields.Keys()[0], "detail")
	is.Equal(fields.Get(fields.Keys()[0]), "file not exist")
	is.Equal(fields.Keys()[1], "file")
	is.Equal(fields.Get(fields.Keys()[1]), "memory_test.go")
	is.Equal(fields.Keys()[2], "one")
	is.Equal(fields.Get(fields.Keys()[2]), "two")
}

func TestErr(t *testing.T) {
	is := is.New(t)
	handler := log.Buffer()
	log := log.New(handler)
	log.Error(errors.New("one"), "two", "three")
	is.Equal(len(handler.Entries), 1)
	is.Equal(handler.Entries[0].Level.String(), "error")
	is.Equal(handler.Entries[0].Message, "one two three")
	fields := handler.Entries[0].Fields
	is.Equal(len(fields), 0)
}
