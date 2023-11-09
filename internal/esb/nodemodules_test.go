package esb_test

import (
	"testing"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/esb"
	"github.com/matryer/is"
)

func TestNodeModulesExternalDOM(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"view/H1.jsx": `
			import React from 'react'
			export default (props) => <h1>{props.children}</h1>
		`,
		"view/Header.jsx": `
			import React from 'react'
			import H1 from './H1.jsx'
			export default (props) => <H1>{props.title}</H1>
		`,
		"view/index.jsx": `
			import slugify from 'slugify'
			import React from 'react'
			import Header from './Header.jsx'
			export function render (props) {
				return <Header title={slugify("index-"+props.title)} />
			}
		`,
		"view/show.jsx": `
			import slugify from 'slugify'
			import React from 'react'
			import Header from './Header.jsx'
			export function render (props) {
				return <Header title={slugify("show-"+props.title)} />
			}
		`,
	})
	file, err := serve(esbuild.BuildOptions{
		AbsWorkingDir: dir,
		EntryPoints:   []string{"./view/index.jsx"},
		Plugins: []esbuild.Plugin{
			esb.ExternalNodeModules("/node_modules"),
		},
	})
	is.NoErr(err)
	code := string(file.Contents)
	isIn(t, code, "/node_modules/slugify")
	isIn(t, code, "/node_modules/react")
	isIn(t, code, "var Header_default")
	isIn(t, code, "export {\n  render\n};")
	// TODO: test against a real browser
}
