package svelte_test

import (
	"strings"
	"testing"

	v8 "github.com/go-duo/bud/js/v8"
	"github.com/go-duo/bud/svelte"
	"github.com/matryer/is"
)

func TestSSR(t *testing.T) {
	is := is.New(t)
	vm := v8.New()
	compiler := svelte.New(&svelte.Input{
		VM:  vm,
		Dev: true,
	})
	ssr, err := compiler.SSR("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(ssr.JS, `import { create_ssr_component } from "svelte/internal";`))
	is.True(strings.Contains(ssr.JS, `<h1>hi world!</h1>`))
}

func TestDOM(t *testing.T) {
	is := is.New(t)
	vm := v8.New()
	compiler := svelte.New(&svelte.Input{
		VM:  vm,
		Dev: true,
	})
	dom, err := compiler.DOM("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(dom.JS, `from "svelte/internal"`))
	is.True(strings.Contains(dom.JS, `function create_fragment`))
	is.True(strings.Contains(dom.JS, `element("h1")`))
	is.True(strings.Contains(dom.JS, `text("hi world!")`))
}
