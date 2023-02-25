package esb_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/esb"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/virtual"
)

func TestServeSSR(t *testing.T) {
	is := is.New(t)
	fsys := virtual.Tree{
		"node_modules/react-dom/server.js": &virtual.File{
			Data: []byte(`export function renderToString() { return "<h1>hello</h1>" }`),
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
	ssr := esb.SSR(&framework.Flag{}, "./view/index.jsx")
	file, err := esb.Serve(http.DefaultTransport, fsys, ssr)
	is.NoErr(err)
	vm, err := v8.Load()
	is.NoErr(err)
	defer vm.Close()
	result, err := vm.Eval("view/index.jsx", fmt.Sprintf(`%s; bud.render({ title: "hello" })`, string(file.Contents)))
	is.NoErr(err)
	is.Equal(result, `<h1>hello</h1>`)
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
