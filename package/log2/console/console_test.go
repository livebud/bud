package console_test

import (
	"testing"

	"github.com/livebud/bud/package/log2/console"
)

func TestConsole(t *testing.T) {
	console.Field("file", "console_test.go").Field("another", "cool story").Debug("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Info("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Notice("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Warn("hello %s", "mars")
	console.Field("file", "console_test.go").Field("another", "cool story").Error("hello %s", "mars")
}
