package console_test

import (
	"errors"
	"os"
	"testing"

	log "github.com/livebud/bud/package/log"
	"github.com/livebud/bud/package/log/console"
)

func TestConsole(t *testing.T) {
	console.Field("file", "console_test.go").Field("another", "cool story").Debugf("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Infof("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Noticef("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Warnf("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Errorf("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Error("hello", "mars")
	console.Error(errors.New("one"), "two", "three")
	logger := log.New(console.New(os.Stdout))
	logger.Error(errors.New("one"), 4, "three")
}
