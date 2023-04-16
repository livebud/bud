package es_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/es"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/testdir"
)

func TestBundleSSR(t *testing.T) {
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
	td.Files["view/show.jsx"] = `
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
	esb := es.New(flag, log)
	files, err := esb.Bundle(&es.Bundle{
		AbsDir:   td.Directory(),
		Entries:  []string{"./view/index.jsx", "./view/show.jsx"},
		Platform: es.SSR,
	})
	is.NoErr(err)
	is.Equal(len(files), 2)
	// TODO: figure out a better test
	// SSR would typically be just one big file anyway
	// DOM is where you chunk the files
	// fmt.Println(string(files[0].Contents))
	// fmt.Println(string(files[1].Contents))

	// vm, err := v8.Load()
	// is.NoErr(err)
	// defer vm.Close()
	// result, err := vm.Eval("view/index.jsx", fmt.Sprintf(`%s; bud.render({ title: "hello" })`, string(file.Contents)))
	// is.NoErr(err)
	// is.Equal(result, `<h1>hello</h1>`)
}
