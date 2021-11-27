package svelte_test

import (
	"strings"
	"testing"

	"github.com/matryer/is"
	v8 "gitlab.com/mnm/bud/js/v8"
	"gitlab.com/mnm/bud/svelte"
)

func TestSSR(t *testing.T) {
	is := is.New(t)
	vm := v8.New()
	compiler := svelte.New(vm)
	ssr, err := compiler.SSR("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(ssr.JS, `import { create_ssr_component } from "svelte/internal";`))
	is.True(strings.Contains(ssr.JS, `<h1>hi world!</h1>`))
}

func TestDOM(t *testing.T) {
	is := is.New(t)
	vm := v8.New()
	compiler := svelte.New(vm)
	dom, err := compiler.DOM("test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(dom.JS, `from "svelte/internal"`))
	is.True(strings.Contains(dom.JS, `function create_fragment`))
	is.True(strings.Contains(dom.JS, `element("h1")`))
	is.True(strings.Contains(dom.JS, `text("hi world!")`))
}

// TODO: test compiler.Dev = false
