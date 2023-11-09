package esb_test

import (
	"strings"
	"testing"

	esbuild "github.com/evanw/esbuild/pkg/api"
	"github.com/livebud/bud/internal/esb"
	"github.com/matryer/is"
)

func TestDisableCSSImportsInJS(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	writeFiles(t, dir, map[string]string{
		"view/index.js": `
			import { uid } from 'https://esm.run/uid'
			import "./index.css"
			export function createElement() { return uid() }
		`,
		"view/index.css": `
			body { background: red; }
		`,
	})
	file, err := serve(esbuild.BuildOptions{
		AbsWorkingDir: dir,
		EntryPoints:   []string{"./view/index.js"},
		Plugins: []esbuild.Plugin{
			esb.DisableCSSImportsInJS(),
		},
	})
	is.True(err != nil)
	is.True(strings.Contains(err.Error(), "not allowed"))
	is.Equal(file, nil)
}
