package tester_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/tester"
)

func TestFiles(t *testing.T) {
	is := is.New(t)
	tr := tester.New(t)
	tr.WriteFiles(map[string]string{
		"ok.txt": "whatever",
	})
	is.True(tr.Exists("ok.txt"))
	is.True(!tr.Exists("okz.txt"))
}

func TestRun(t *testing.T) {
	tr := tester.New(t)
	result := tr.Run("echo", "hi").NoErr()
	result.Stdout.Equal("hi\n")
}
