package v8client_test

import (
	"testing"

	"github.com/go-duo/bud/js/v8client"
	"github.com/matryer/is"
)

func TestLaunch(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	client := v8client.Launch("duo", "tool", "v8")
	result, err := client.Eval("stdin", "10 + 3")
	is.NoErr(err)
	is.Equal(result, "13")
	result, err = client.Eval("stdin", "10 + 5")
	is.NoErr(err)
	is.Equal(result, "15")
}
