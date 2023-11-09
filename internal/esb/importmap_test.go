package esb_test

import (
	"net/http"
	"testing"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/esb"
	"github.com/livebud/bud/pkg/logs"
	"github.com/matryer/is"
)

func TestImportMap(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"view/index.js": `
			import { SvelteComponent } from 'svelte'
			import { onMount } from 'svelte/internal'
			export function createElement() { console.log(SvelteComponent, onMount) }
		`,
	})
	file, err := serve(esbuild.BuildOptions{
		AbsWorkingDir: dir,
		EntryPoints:   []string{"./view/index.js"},
		Plugins: []esbuild.Plugin{
			esb.HTTP(http.DefaultClient),
			esb.ImportMap(logs.New(logs.Buffer()), map[string]string{
				"svelte":  "https://esm.run/svelte@4",
				"svelte/": "https://esm.run/svelte@4/",
			}),
		},
	})
	is.NoErr(err)
	code := string(file.Contents)
	isIn(t, code, `function createElement() {`)
	isIn(t, code, `"Component was already destroyed"`)
	isIn(t, code, `on_mount`)
}
