package svelte_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/viewer/svelte"
	"github.com/livebud/js"
	"github.com/livebud/js/goja"
)

func TestSSR(t *testing.T) {
	is := is.New(t)
	js := goja.New(&js.Console{
		Log:   os.Stdout,
		Error: os.Stderr,
	})
	flag := &framework.Flag{}
	ctx := context.Background()
	compiler, err := svelte.Load(flag, js)
	is.NoErr(err)
	ssr, err := compiler.SSR(ctx, "test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(ssr.JS, `import { create_ssr_component } from "svelte/internal";`))
	is.True(strings.Contains(ssr.JS, `<h1>hi world!</h1>`))
}

func TestSSRRecovery(t *testing.T) {
	is := is.New(t)
	js := goja.New(&js.Console{
		Log:   os.Stdout,
		Error: os.Stderr,
	})
	flag := &framework.Flag{}
	ctx := context.Background()
	compiler, err := svelte.Load(flag, js)
	is.NoErr(err)
	ssr, err := compiler.SSR(ctx, "test.svelte", []byte(`<h1>hi world!</h1></h1>`))
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), `</h1> attempted to close an element that was not open`))
	is.True(strings.Contains(err.Error(), `<h1>hi world!</h1></h1`))
	is.Equal(ssr, nil)
	ssr, err = compiler.SSR(ctx, "test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(ssr.JS, `import { create_ssr_component } from "svelte/internal";`))
	is.True(strings.Contains(ssr.JS, `<h1>hi world!</h1>`))
}

func TestDOM(t *testing.T) {
	is := is.New(t)
	js := goja.New(&js.Console{
		Log:   os.Stdout,
		Error: os.Stderr,
	})
	flag := &framework.Flag{}
	ctx := context.Background()
	compiler, err := svelte.Load(flag, js)
	is.NoErr(err)
	dom, err := compiler.DOM(ctx, "test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(dom.JS, `from "svelte/internal"`))
	is.True(strings.Contains(dom.JS, `function create_fragment`))
	is.True(strings.Contains(dom.JS, `element("h1")`))
	is.True(strings.Contains(dom.JS, `text("hi world!")`))
}

func TestDOMRecovery(t *testing.T) {
	is := is.New(t)
	js := goja.New(&js.Console{
		Log:   os.Stdout,
		Error: os.Stderr,
	})
	flag := &framework.Flag{}
	ctx := context.Background()
	compiler, err := svelte.Load(flag, js)
	is.NoErr(err)
	dom, err := compiler.DOM(ctx, "test.svelte", []byte(`<h1>hi world!</h1></h1>`))
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), `</h1> attempted to close an element that was not open`))
	is.True(strings.Contains(err.Error(), `<h1>hi world!</h1></h1`))
	is.Equal(dom, nil)
	dom, err = compiler.DOM(ctx, "test.svelte", []byte(`<h1>hi world!</h1>`))
	is.NoErr(err)
	is.True(strings.Contains(dom.JS, `from "svelte/internal"`))
	is.True(strings.Contains(dom.JS, `function create_fragment`))
	is.True(strings.Contains(dom.JS, `element("h1")`))
	is.True(strings.Contains(dom.JS, `text("hi world!")`))
}

// TODO: test compiler.Dev = false
