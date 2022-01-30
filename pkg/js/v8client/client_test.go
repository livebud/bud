package v8client_test

import (
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/pkg/js/v8client"
)

func TestLaunch(t *testing.T) {
	is := is.New(t)
	client := v8client.New("bud", "tool", "v8", "client")
	result, err := client.Eval("stdin.js", "10 + 3")
	is.NoErr(err)
	is.Equal(result, "13")
	result, err = client.Eval("stdin.js", "10 + 5")
	is.NoErr(err)
	is.Equal(result, "15")
}
