package es_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/livebud/bud/framework"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/es"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/testdir"
)

func TestHTTP(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["view/index.js"] = `
		import { uid } from 'https://esm.run/uid'
		export function createElement() { return uid() }
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
		},
	})
	is.NoErr(err)
	code := string(file.Contents)
	is.In(code, "function createElement() {")
	is.In(code, "(t + 256).toString(16).substring(1)")
	is.In(code, "Math.random()")
}

func TestHTTPDepOfDep(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	log := testlog.New()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Files["view/index.js"] = `
		import { uid } from 'https://esm.run/uid/secure'
		export function createElement() { return uid() }
	`
	is.NoErr(td.Write(ctx))
	flag := &framework.Flag{}
	esb := es.New(flag, log)
	file, err := esb.Serve(&es.Serve{
		AbsDir:   td.Directory(),
		Entry:    "./view/index.js",
		Platform: es.SSR,
		Plugins: []es.Plugin{
			es.HTTP(http.DefaultClient),
		},
	})
	is.NoErr(err)
	code := string(file.Contents)
	is.In(code, "function createElement() {")
	is.In(code, "(n + 256).toString(16).substring(1)")
	is.In(code, "crypto.getRandomValues")
}
