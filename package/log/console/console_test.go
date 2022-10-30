package console_test

import (
	"errors"
	"os"
	"testing"

	log "github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
)

func TestConsole(t *testing.T) {
	console.Field("file", "console_test.go").Field("another", "cool story").Debug("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Info("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Notice("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Warn("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Error("hello %s", "mars")
	console.Err(errors.New("one"), "two %s", "three")
	logger := log.New(console.New(os.Stdout))
	logger.Err(errors.New("one"), "two %s", "three")
}
