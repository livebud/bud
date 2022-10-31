package memory_test

import (
	"errors"
	"testing"

	"github.com/livebud/bud/internal/is"
	log "github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/memory"
)

func TestLog(t *testing.T) {
	is := is.New(t)
	handler := memory.New()
	log := log.New(handler)
	log.Info("hello")
	log.Info("world")
	log.Warn("hello %s!", "mars")
	log.Fields(map[string]interface{}{
		"file":   "memory_test.go",
		"detail": "file not exist",
	}).Field("one", "two").Error("%d. oh noz", 1)
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
	handler := memory.New()
	log := log.New(handler)
	log.Err(errors.New("one"), "two %s", "three")
	is.Equal(len(handler.Entries), 1)
	is.Equal(handler.Entries[0].Level.String(), "error")
	is.Equal(handler.Entries[0].Message, "two three")
	fields := handler.Entries[0].Fields
	is.Equal(len(fields), 2)
	is.Equal(fields.Keys()[0], "error")
	is.Equal(fields.Get(fields.Keys()[0]), "one")
	is.Equal(fields.Keys()[1], "source")
	source := fields.Get(fields.Keys()[1])
	is.In(source, "memory_test.TestErr:46") // line number may change
}
