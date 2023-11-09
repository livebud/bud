package esb_test

import (
	"net/http"
	"testing"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/esb"
	"github.com/matryer/is"
)

func TestHTTP(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"view/index.js": `
			import { uid } from 'https://esm.run/uid'
			export function createElement() { return uid() }
		`,
	})
	file, err := serve(esbuild.BuildOptions{
		AbsWorkingDir: dir,
		EntryPoints:   []string{"./view/index.js"},
		Plugins: []esbuild.Plugin{
			esb.HTTP(http.DefaultClient),
		},
	})
	is.NoErr(err)
	code := string(file.Contents)
	isIn(t, code, "function createElement() {")
	isIn(t, code, "(t + 256).toString(16).substring(1)")
	isIn(t, code, "Math.random()")
}

func TestHTTPDepOfDep(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"view/index.js": `
			import { uid } from 'https://esm.run/uid/secure'
			export function createElement() { return uid() }
		`,
	})
	file, err := serve(esbuild.BuildOptions{
		AbsWorkingDir: dir,
		EntryPoints:   []string{"./view/index.js"},
		Plugins: []esbuild.Plugin{
			esb.HTTP(http.DefaultClient),
		},
	})
	is.NoErr(err)
	code := string(file.Contents)
	isIn(t, code, "function createElement() {")
	isIn(t, code, "(n + 256).toString(16).substring(1)")
	isIn(t, code, "crypto.getRandomValues")
}
