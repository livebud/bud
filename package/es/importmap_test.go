package es_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/es"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/testdir"
)

func TestImportMap(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["view/index.js"] = `
		import { SvelteComponent } from 'svelte'
		import { create_component } from 'svelte/internal'
		export function createElement() { console.log(SvelteComponent, create_component) }
	`
	is.NoErr(td.Write(ctx))
	flag := &framework.Flag{}
	esb := es.New(flag, log)
	file, err := esb.Serve(&es.Serve{
		AbsDir:   td.Directory(),
		Entry:    "./view/index.js",
		Platform: es.DOM,
		Plugins: []es.Plugin{
			es.HTTP(http.DefaultClient),
			es.ImportMap(log, map[string]string{
				"svelte":  "https://esm.run/svelte@" + versions.Svelte,
				"svelte/": "https://esm.run/svelte@" + versions.Svelte + "/",
			}),
		},
	})
	is.NoErr(err)
	code := string(file.Contents)
	is.In(code, `function createElement() {`)
	is.In(code, `"Component was already destroyed"`)
	is.In(code, `on_mount`)
}
