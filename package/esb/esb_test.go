package esb_test

import (
	"fmt"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/esb"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/virtual"
)

func TestServeSSR(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.Tree{
		"node_modules/uid/index.js": &virtual.File{
			Data: []byte(`export default function uid() { return "uid" }`),
		},
		"node_modules/react-dom/server.js": &virtual.File{
			Data: []byte(`
				import uid from 'uid'
				export function renderToString() { return "<h1>hello</h1>" + uid() }
			`),
		},
		"node_modules/react/index.js": &virtual.File{
			Data: []byte(`export function createElement() { return {} }`),
		},
		"node_modules/@pkg/slugify/index.mjs": &virtual.File{
			Data: []byte(`export default function slugify(title) { return title }`),
		},
		"view/Header.jsx": &virtual.File{
			Data: []byte(`export default (props) => <h1>{props.title}</h1>`),
		},
		"view/index.jsx": &virtual.File{
			Data: []byte(`
				import { renderToString } from 'react-dom/server'
				import slugify from '@pkg/slugify'
				import * as React from 'react'
				import Header from './Header.jsx'
				export function render (props) {
					return renderToString(<Header title={slugify(props.title)} />)
				}
			`),
		},
	}
	builder := esb.New(fsys, log)
	flag := new(framework.Flag)
	ssr := esb.SSR(flag)
	ssr.Plugins = append(ssr.Plugins, esb.FS(fsys, "virtual"))
	file, err := builder.Serve(ssr, "./view/index.jsx")
	is.NoErr(err)
	vm, err := v8.Load()
	is.NoErr(err)
	defer vm.Close()
	result, err := vm.Eval("view/index.jsx", fmt.Sprintf(`%s; bud.render({ title: "hello" })`, string(file.Contents)))
	is.NoErr(err)
	is.Equal(result, `<h1>hello</h1>uid`)
}

func TestNoFSLeak(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.Tree{}
	builder := esb.New(fsys, log)
	ssr := esb.SSR(&framework.Flag{})
	file, err := builder.Serve(ssr, "./testdata/leak/index.js")
	is.True(err != nil)
	is.True(file == nil)
	is.In(err.Error(), `unable to find "./testdata/leak/index.js"`)
}

func TestFS(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.Tree{
		"view/index.jsx": &virtual.File{
			Data: []byte(`
				import * as React from 'react'
				export default function() {
					return <h1>hello</h1>
				}
			`),
		},
		"view/index.jsx.js": &virtual.File{
			Data: []byte(`
				import { renderToString } from 'react-dom/server'
				import { createElement } from 'react'
				import jsesc from 'jsesc'
				import Index from './index.jsx'
				export function render(props = {}) {
					return renderToString(createElement(Index, props, []))
				}
			`),
		},
	}
	builder := esb.New(fsys, log)
	flag := new(framework.Flag)
	ssr := esb.SSR(flag)
	viewerFS := virtual.Tree{
		"node_modules/react-dom/server.js": &virtual.File{
			Data: []byte(`export function renderToString() { return "<h1>hello</h1>" }`),
		},
		"node_modules/react/index.js": &virtual.File{
			Data: []byte(`export function createElement() { return {} }`),
		},
		"node_modules/jsesc/index.js": &virtual.File{
			Data: []byte(`export default function(input) { return input }`),
		},
	}
	ssr.Plugins = append(ssr.Plugins, esb.FS(fsys, "virtual"), esb.FS(viewerFS, "react"))
	file, err := builder.Serve(ssr, "./view/index.jsx.js")
	is.NoErr(err)
	vm, err := v8.Load()
	is.NoErr(err)
	defer vm.Close()
	result, err := vm.Eval("view/index.jsx.js", fmt.Sprintf(`%s; bud.render({ title: "hello" })`, string(file.Contents)))
	is.NoErr(err)
	is.Equal(result, `<h1>hello</h1>`)
}

func TestNodeModuleSSR(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := virtual.Tree{
		"node_modules/jsesc/jsesc.js": &virtual.File{
			Data: []byte(`export default function jsesc() { return "jsesc" }`),
		},
		"view/index.jsx": &virtual.File{
			Data: []byte(`
				import jsesc from 'jsesc'
				export function render () {
					return jsesc()
				}
			`),
		},
	}
	builder := esb.New(fsys, log)
	flag := new(framework.Flag)
	ssr := esb.SSR(flag)
	ssr.Plugins = append(ssr.Plugins, esb.FS(fsys, "virtual"))
	file, err := builder.Serve(ssr, "./view/index.jsx")
	is.NoErr(err)
	vm, err := v8.Load()
	is.NoErr(err)
	defer vm.Close()
	result, err := vm.Eval("view/index.jsx", fmt.Sprintf(`%s; bud.render()`, string(file.Contents)))
	is.NoErr(err)
	is.Equal(result, `jsesc`)
}

func TestNodeModuleDOM(t *testing.T) {
}

// TODO: Test serve DOM relative entry
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
