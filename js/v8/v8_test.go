package v8_test

import (
	"testing"

	v8 "github.com/go-duo/bud/js/v8"
	"github.com/matryer/is"
)

func TestEval(t *testing.T) {
	is := is.New(t)
	result, err := v8.Eval("TestEval.js", "2*5")
	is.NoErr(err)
	is.Equal("10", result)
}

func TestScript(t *testing.T) {
	is := is.New(t)
	v8 := v8.New()
	v8.Script("bootstrap.js", `
		function multiply(x, y) {
			return x * y
		}
	`)
	result, err := v8.Eval("TestScript.js", "multiply(2, 10)")
	is.NoErr(err)
	is.Equal("20", result)
	result, err = v8.Eval("TestScript.js", "multiply(2, 5)")
	is.NoErr(err)
	is.Equal("10", result)
}

func TestConsole(t *testing.T) {
	is := is.New(t)
	_, err := v8.Eval("TestConsole.js", `console.log("a", 3, { hi: "world" })`)
	is.NoErr(err)
}
func TestFetch(t *testing.T) {
	is := is.New(t)
	res, err := v8.Eval("TestFetch.js", `fetch("http://google.com").then(res => res.status)`)
	is.NoErr(err)
	is.Equal(res, "200")
}
