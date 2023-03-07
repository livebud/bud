package es_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/es"
	"github.com/livebud/bud/package/gomod"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/testdir"
)

func TestServeSSR(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["node_modules/react-dom/server.js"] = `
		export function renderToString() { return "<h1>hello</h1>" }
	`
	td.Files["node_modules/react/index.js"] = `
		export function createElement() { return {} }
	`
	td.Files["node_modules/@pkg/slugify/index.js"] = `
		export default function slugify(title) { return title }
	`
	td.Files["view/H1.jsx"] = `
		export default (props) => <h1>{props.children}</h1>
	`
	td.Files["view/Header.jsx"] = `
		import H1 from './H1.jsx'
		export default (props) => <H1>{props.title}</H1>
	`
	td.Files["view/index.jsx"] = `
		import { renderToString } from 'react-dom/server'
		import slugify from '@pkg/slugify'
		import * as React from 'react'
		import Header from './Header.jsx'
		export function render (props) {
			return renderToString(<Header title={slugify(props.title)} />)
		}
	`
	is.NoErr(td.Write(ctx))
	flag := &framework.Flag{}
	module, err := gomod.Find(td.Directory())
	is.NoErr(err)
	esb := es.New(flag, log, module)
	file, err := esb.Serve(&es.Serve{
		Entry:    "./view/index.jsx",
		Platform: es.SSR,
	})
	is.NoErr(err)
	vm, err := v8.Load()
	is.NoErr(err)
	defer vm.Close()
	result, err := vm.Eval("view/index.jsx", fmt.Sprintf(`%s; bud.render({ title: "hello" })`, string(file.Contents)))
	is.NoErr(err)
	is.Equal(result, `<h1>hello</h1>`)
}

func TestServeDOM(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["node_modules/react-dom/client.js"] = `
		export function renderToString() { return "<h1>hello</h1>" }
	`
	td.Files["node_modules/react/index.js"] = `
		export function createElement() { return {} }
	`
	td.Files["node_modules/@pkg/slugify/index.js"] = `
		export default function slugify(title) { return title }
	`
	td.Files["view/H1.jsx"] = `
		export default (props) => <h1>{props.children}</h1>
	`
	td.Files["view/Header.jsx"] = `
		import H1 from './H1.jsx'
		export default (props) => <H1>{props.title}</H1>
	`
	td.Files["view/index.jsx"] = `
		import { renderToString } from 'react-dom/client'
		import slugify from '@pkg/slugify'
		import * as React from 'react'
		import Header from './Header.jsx'
		export function render (props) {
			return renderToString(<Header title={slugify(props.title)} />)
		}
	`
	is.NoErr(td.Write(ctx))
	flag := &framework.Flag{}
	module, err := gomod.Find(td.Directory())
	is.NoErr(err)
	esb := es.New(flag, log, module)
	file, err := esb.Serve(&es.Serve{
		Entry:    "./view/index.jsx",
		Platform: es.DOM,
	})
	is.NoErr(err)
	code := string(file.Contents)
	is.In(code, "/node_modules/react-dom/client")
	is.In(code, "/node_modules/@pkg/slugify")
	is.In(code, "/node_modules/react")
	is.NotIn(code, "./Header.jsx")
	is.In(code, "export {\n  render\n};")
}

func TestServeModuleDOM(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["node_modules/react-dep/index.js"] = `
		export function dep() { return {} }
	`
	td.Files["node_modules/react/index.js"] = `
		import { dep } from 'react-dep'
		export function createElement() { return dep() }
	`
	is.NoErr(td.Write(ctx))
	flag := &framework.Flag{}
	module, err := gomod.Find(td.Directory())
	is.NoErr(err)
	esb := es.New(flag, log, module)
	file, err := esb.Serve(&es.Serve{
		Entry:    "react",
		Platform: es.DOM,
	})
	is.True(err != nil)
	is.True(errors.Is(err, es.ErrNotRelative))
	is.Equal(file, nil)
}

func TestServeRelModuleDOM(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["node_modules/react-dep/index.js"] = `
		export function dep() { return {} }
	`
	td.Files["node_modules/react/index.js"] = `
		import { dep } from 'react-dep'
		export function createElement() { return dep() }
	`
	is.NoErr(td.Write(ctx))
	flag := &framework.Flag{}
	module, err := gomod.Find(td.Directory())
	is.NoErr(err)
	esb := es.New(flag, log, module)
	file, err := esb.Serve(&es.Serve{
		Entry:    "./node_modules/react",
		Platform: es.DOM,
	})
	is.NoErr(err)
	code := string(file.Contents)
	is.In(code, "function dep() {")
	is.In(code, "function createElement() {")
	is.In(code, "export {\n  createElement\n};")
}

func TestServeRelModuleSubpathDOM(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["node_modules/react-dep/index.js"] = `
		export function dep() { return {} }
	`
	td.Files["node_modules/react/client.js"] = `
		import { dep } from 'react-dep'
		export function createElement() { return dep() }
	`
	is.NoErr(td.Write(ctx))
	flag := &framework.Flag{}
	module, err := gomod.Find(td.Directory())
	is.NoErr(err)
	esb := es.New(flag, log, module)
	file, err := esb.Serve(&es.Serve{
		Entry:    "./node_modules/react/client",
		Platform: es.DOM,
	})
	is.NoErr(err)
	code := string(file.Contents)
	is.In(code, "function dep() {")
	is.In(code, "function createElement() {")
	is.In(code, "export {\n  createElement\n};")
}

func TestServeScopedModuleDOM(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["node_modules/slugify-dep/index.js"] = `
		export function dep(title) { return title }
	`
	td.Files["node_modules/@pkg/slugify/index.js"] = `
		import { dep } from 'slugify-dep'
		export default function slugify(title) { return dep(title) }
	`
	is.NoErr(td.Write(ctx))
	flag := &framework.Flag{}
	module, err := gomod.Find(td.Directory())
	is.NoErr(err)
	esb := es.New(flag, log, module)
	file, err := esb.Serve(&es.Serve{
		Entry:    "@pkg/slugify",
		Platform: es.DOM,
	})
	is.True(err != nil)
	is.True(errors.Is(err, es.ErrNotRelative))
	is.Equal(file, nil)
}

func TestServeRelScopedModuleSubpathDOM(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["node_modules/slugify-dep/index.js"] = `
		export function dep(title) { return title }
	`
	td.Files["node_modules/@pkg/slugify/titles.js"] = `
		export function titles(title) { return title }
	`
	td.Files["node_modules/@pkg/slugify/client.js"] = `
		import { dep } from 'slugify-dep'
		import { titles } from './titles'
		export default function slugify(title) { return dep(titles(title)) }
	`
	is.NoErr(td.Write(ctx))
	flag := &framework.Flag{}
	module, err := gomod.Find(td.Directory())
	is.NoErr(err)
	esb := es.New(flag, log, module)
	file, err := esb.Serve(&es.Serve{
		Entry:    "./node_modules/@pkg/slugify/client",
		Platform: es.DOM,
	})
	is.NoErr(err)
	code := string(file.Contents)
	is.In(code, `function titles(title)`)
	is.In(code, `function dep(title)`)
	is.In(code, "export {\n  slugify as default\n};")
}

// TODO: Test serve DOM node_module entry (e.g. "./node_modules/react")
// TODO: test resolving different relative path extensions (e.g. ./Header.svelte)
// TODO: test resolving different node_modules path extensions (e.g. "./node_modules/@ui/Grid.svelte")
// TODO: test dependencies of dependencies

// TODO: Test bundle SSR relative entries
// TODO: Test bundle DOM relative entries
// TODO: test minifying

// TODO: test tilde ~
// TODO: test http
// TODO: test resolving svelte from within an embedded virtual file system
// TODO: test injecting variables
