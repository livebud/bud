package logs_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/livebud/bud/internal/color"
	"github.com/livebud/bud/pkg/logs"
	"github.com/matryer/is"
)

func TestConsole(t *testing.T) {
	is := is.New(t)
	buf := new(bytes.Buffer)
	log := logs.New(logs.Console(color.New(), buf))
	log.Debug("world", "args", 10)
	log.Info("hello", "planet", "world", "args", 10)
	log.Warn("hello", "planet", "world", "args", 10)
	log.Error("hello world", "planet", "world", "args", 10)
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	is.Equal(len(lines), 4)
}

func ExampleConsole() {
	log := logs.New(logs.Console(color.Ignore(), os.Stdout))
	log.Debug("world", "args", 10)
	log.Info("hello", "planet", "world", "args", 10)
	log.Warn("hello", "planet", "world", "args", 10)
	// log.Error("hello world", slog.String("planet", "world"), "args", 10)
	// Output:
	// debug: world hello.args=10
	// info: hello planet=world args=10
	// warn: hello planet=world args=10
	// error: hello world planet=world args=10
}
