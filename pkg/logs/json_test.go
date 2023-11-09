package logs_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/livebud/bud/pkg/logs"
	"github.com/matryer/is"
)

func TestJson(t *testing.T) {
	is := is.New(t)
	buf := new(bytes.Buffer)
	log := logs.New(logs.Json(buf))
	log.Field("args", 10).Debug("hello")
	log.Field("args", 10).Field("planet", "world").Info("hello")
	log.Warn("hello")
	log.Field("planet", "world").Field("args", 10).Error("hello", "world")
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	is.Equal(len(lines), 4)
	fmt.Println(lines)
}
