package goja_test

import (
	"os"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/pkg/js/goja"
)

func TestAdd(t *testing.T) {
	is := is.New(t)
	vm := goja.New()
	result, err := vm.Eval("script.js", "1+1")
	is.NoErr(err)
	is.Equal(result, "2")
}

func TestSvelte(t *testing.T) {
	is := is.New(t)
	vm := goja.New()
	compiler, err := os.ReadFile("testdata/await.js")
	is.NoErr(err)
	err = vm.Script("compiler.js", string(compiler))
	is.NoErr(err)
}
