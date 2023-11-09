package logs_test

import (
	"errors"
	"testing"

	"github.com/livebud/bud/pkg/logs"
	"github.com/matryer/is"
)

func TestSampleLog(t *testing.T) {
	is := is.New(t)
	log := logs.Buffer()
	log.Info("hello")
	log.Info("world")
	log.Warnf("hello %s!", "mars")
	log.Fields(map[string]interface{}{
		"file":   "memory_test.go",
		"detail": "file not exist",
	}).Field("one", "two").Errorf("%d. oh noz", 1)
	entries := log.Entries()
	is.Equal(len(entries), 4)
	is.Equal(entries[0].Level.String(), "info")
	is.Equal(entries[0].Message, "hello")
	is.Equal(entries[1].Level.String(), "info")
	is.Equal(entries[1].Message, "world")
	is.Equal(entries[2].Level.String(), "warn")
	is.Equal(entries[2].Message, "hello mars!")
	is.Equal(entries[3].Level.String(), "error")
	is.Equal(entries[3].Message, "1. oh noz")
	fields := entries[3].Fields
	is.Equal(len(fields), 3)
	is.Equal(fields.Keys()[0], "detail")
	is.Equal(fields.Get(fields.Keys()[0]), "file not exist")
	is.Equal(fields.Keys()[1], "file")
	is.Equal(fields.Get(fields.Keys()[1]), "memory_test.go")
	is.Equal(fields.Keys()[2], "one")
	is.Equal(fields.Get(fields.Keys()[2]), "two")
}

func TestErrLog(t *testing.T) {
	is := is.New(t)
	log := logs.Buffer()
	log.Error(errors.New("one"), "two", "three")
	entries := log.Entries()
	is.Equal(len(entries), 1)
	is.Equal(entries[0].Level.String(), "error")
	is.Equal(entries[0].Message, "one two three")
	fields := entries[0].Fields
	is.Equal(len(fields), 0)
}

func TestSampleHandler(t *testing.T) {
	is := is.New(t)
	handler := logs.Buffer()
	log := logs.New(handler)
	log.Info("hello")
	log.Info("world")
	log.Warnf("hello %s!", "mars")
	log.Fields(map[string]interface{}{
		"file":   "memory_test.go",
		"detail": "file not exist",
	}).Field("one", "two").Errorf("%d. oh noz", 1)
	entries := handler.Entries()
	is.Equal(len(entries), 4)
	is.Equal(entries[0].Level.String(), "info")
	is.Equal(entries[0].Message, "hello")
	is.Equal(entries[1].Level.String(), "info")
	is.Equal(entries[1].Message, "world")
	is.Equal(entries[2].Level.String(), "warn")
	is.Equal(entries[2].Message, "hello mars!")
	is.Equal(entries[3].Level.String(), "error")
	is.Equal(entries[3].Message, "1. oh noz")
	fields := entries[3].Fields
	is.Equal(len(fields), 3)
	is.Equal(fields.Keys()[0], "detail")
	is.Equal(fields.Get(fields.Keys()[0]), "file not exist")
	is.Equal(fields.Keys()[1], "file")
	is.Equal(fields.Get(fields.Keys()[1]), "memory_test.go")
	is.Equal(fields.Keys()[2], "one")
	is.Equal(fields.Get(fields.Keys()[2]), "two")
}

func TestErrHandler(t *testing.T) {
	is := is.New(t)
	handler := logs.Buffer()
	log := logs.New(handler)
	log.Error(errors.New("one"), "two", "three")
	entries := handler.Entries()
	is.Equal(len(entries), 1)
	is.Equal(entries[0].Level.String(), "error")
	is.Equal(entries[0].Message, "one two three")
	fields := entries[0].Fields
	is.Equal(len(fields), 0)
}
