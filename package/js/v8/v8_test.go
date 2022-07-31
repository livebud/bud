package v8_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	v8 "github.com/livebud/bud/package/js/v8"
)

func TestCompile(t *testing.T) {
	is := is.New(t)
	vm, err := v8.Compile("math.js", `const multiply = (a, b) => a * b`)
	is.NoErr(err)
	defer vm.Close()
	value, err := vm.Eval("run.js", "multiply(3, 2)")
	is.NoErr(err)
	is.Equal("6", value)
}

func TestEval(t *testing.T) {
	is := is.New(t)
	result, err := v8.Eval("TestEval.js", "2*5")
	is.NoErr(err)
	is.Equal("10", result)
}

func TestConsoleLog(t *testing.T) {
	is := is.New(t)
	_, err := v8.Eval("TestConsole.js", `console.log("a", 3, { hi: "world" })`)
	is.NoErr(err)
}

func TestConsoleWarn(t *testing.T) {
	is := is.New(t)
	_, err := v8.Eval("TestConsole.js", `console.warn("a", 3, { hi: "world" })`)
	is.NoErr(err)
}

func TestConsoleError(t *testing.T) {
	is := is.New(t)
	_, err := v8.Eval("TestConsole.js", `console.error("a", 3, { hi: "world" })`)
	is.NoErr(err)
}

func TestFetch(t *testing.T) {
	is := is.New(t)
	res, err := v8.Eval("TestFetch.js", `fetch("http://google.com").then(res => res.status)`)
	is.NoErr(err)
	is.Equal(res, "200")
}

func TestURL(t *testing.T) {
	is := is.New(t)
	res, err := v8.Eval("TestURL.js", `(new URL("http://google.com/hi")).host`)
	is.NoErr(err)
	is.Equal(res, "google.com")
}

func TestURLSearchParams(t *testing.T) {
	is := is.New(t)
	res, err := v8.Eval("TestURLSearchParams.js", `(new URLSearchParams('?a=1')).get('a')`)
	is.NoErr(err)
	is.Equal(res, "1")
}

func TestSetClearInterval(t *testing.T) {
	is := is.New(t)
	res, err := v8.Eval("TestSetClearInterval.js", `let id = setInterval(() => clearInterval(id), 500)`)
	is.NoErr(err)
	is.Equal(res, "undefined")
}

func TestSetClearTimeout(t *testing.T) {
	is := is.New(t)
	res, err := v8.Eval("TestSetClearTimeout.js", `let id = setTimeout(() => clearTimeout(id), 500)`)
	is.NoErr(err)
	is.Equal(res, "undefined")
}
