package testlog

import (
	"flag"
	"fmt"
	"os"

	log "github.com/livebud/bud/package/log2"
	"github.com/livebud/bud/package/log2/console"
)

var logFlag = flag.String("log", "info", "choose a log level")

// Pattern returns the log level logFlag so we can pass it through arguments.
func Pattern() string {
	return *logFlag
}

// New logger for testing. You can set the log level for a given test by
// using the --log=<pattern> flag. For example, `go test ./... --log=debug` will
// run tests with debug logs on.
func New() log.Log {
	level, err := log.ParseLevel(*logFlag)
	if err != nil {
		panic(fmt.Sprintf("testlog: invalid --log=[level] %q" + *logFlag))
	}
	return log.New(level, console.New(os.Stderr))
}
