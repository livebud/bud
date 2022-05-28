package dom_test

import (
	"context"
	"fmt"
	"io/fs"
	"strings"
	"testing"

	"github.com/livebud/bud/package/overlay"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/runtime/transform"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/version"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/svelte"
	"github.com/livebud/bud/runtime/view/dom"
)

func TestServeFile(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	vm, err := v8.Load()
	is.NoErr(err)
	svelteCompiler, err := svelte.Load(vm)
	is.NoErr(err)
	transformer := transform.MustLoad(
		svelte.NewTransformable(svelteCompiler),
	)
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>index</h1>`
	td.Files["view/about/index.svelte"] = `<h2>about</h2>`
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	overlay, err := overlay.Load(module)
	is.NoErr(err)
	overlay.FileServer("bud/view", dom.New(module, transformer.DOM))
	// Read the wrapped version of index.svelte with node_modules rewritten
	code, err := fs.ReadFile(overlay, "bud/view/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h1");`))
	is.True(strings.Contains(string(code), `text("index")`))
	is.True(strings.Contains(string(code), `"/bud/view/index.svelte": view_default`))
	is.True(strings.Contains(string(code), `page: "/bud/view/index.svelte",`))
	is.True(strings.Contains(string(code), `hot: new Hot("http://127.0.0.1:35729/?page=%2Fbud%2Fview%2Findex.svelte", components)`))

	// Unwrapped version with node_modules rewritten
	code, err = fs.ReadFile(overlay, "bud/view/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h1");`))
	is.True(strings.Contains(string(code), `text("index")`))
	// Unwrapped version doesn't contain wrapping
	is.True(!strings.Contains(string(code), `"/bud/view/index.svelte": view_default`))
	is.True(!strings.Contains(string(code), `page: "/bud/view/index.svelte",`))
	is.True(!strings.Contains(string(code), `hot: new Hot("http://127.0.0.1:35729/?page=%2Fbud%2Fview%2Findex.svelte", components)`))

	// Read the wrapped version of about/index.svelte with node_modules rewritten
	code, err = fs.ReadFile(overlay, "bud/view/about/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h2");`))
	is.True(strings.Contains(string(code), `text("about")`))
	is.True(strings.Contains(string(code), `"/bud/view/about/index.svelte": about_default`))
	is.True(strings.Contains(string(code), `page: "/bud/view/about/index.svelte",`))
	is.True(strings.Contains(string(code), `hot: new Hot("http://127.0.0.1:35729/?page=%2Fbud%2Fview%2Fabout%2Findex.svelte", components)`))

	// Unwrapped version with node_modules rewritten
	code, err = fs.ReadFile(overlay, "bud/view/about/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h2");`))
	is.True(strings.Contains(string(code), `text("about")`))
	// Unwrapped version doesn't contain wrapping
	is.True(!strings.Contains(string(code), `"/bud/view/about/index.svelte": about_default`))
	is.True(!strings.Contains(string(code), `page: "/bud/view/about/index.svelte",`))
	is.True(!strings.Contains(string(code), `hot: new Hot("http://127.0.0.1:35729/?page=%2Fbud%2Fview%2Fabout%2Findex.svelte", components)`))
}

func TestNodeModules(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>hi world</h1>`
	td.NodeModules["svelte"] = version.Svelte
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	overlay, err := overlay.Load(module)
	is.NoErr(err)
	overlay.FileServer("bud/node_modules", dom.NodeModules(module))
	// Read the re-written node_modules
	code, err := fs.ReadFile(overlay, "bud/node_modules/svelte/internal")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `function element(`))
	is.True(strings.Contains(string(code), `function text(`))
}

func TestGenerateDir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>index</h1>`
	td.Files["view/about/index.svelte"] = `<h2>about</h2>`
	td.NodeModules["livebud"] = "*"
	td.NodeModules["svelte"] = version.Svelte
	is.NoErr(td.Write(ctx))
	vm, err := v8.Load()
	is.NoErr(err)
	svelteCompiler, err := svelte.Load(vm)
	is.NoErr(err)
	transformer := transform.MustLoad(svelte.NewTransformable(svelteCompiler))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	overlay, err := overlay.Load(module)
	is.NoErr(err)
	overlay.DirGenerator("bud/view", dom.New(module, transformer.DOM))
	des, err := fs.ReadDir(overlay, "bud/view")
	is.NoErr(err)
	is.Equal(len(des), 3)
	is.Equal(des[0].Name(), "_index.svelte")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[1].Name(), "about")
	is.Equal(des[1].IsDir(), true)
	is.True(strings.HasPrefix(des[2].Name(), "chunk-"))
	is.Equal(des[2].IsDir(), false)
	chunkName := des[2].Name()
	des, err = fs.ReadDir(overlay, "bud/view/about")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "_index.svelte")
	is.Equal(des[0].IsDir(), false)

	code, err := fs.ReadFile(overlay, "bud/view/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"H1"`))
	is.True(strings.Contains(string(code), `"index"`))
	is.True(strings.Contains(string(code), fmt.Sprintf(`from"./%s"`, chunkName)))
	is.True(strings.Contains(string(code), `page:"/bud/view/index.svelte"`))
	is.True(strings.Contains(string(code), `document.getElementById("bud_target")`))
	// TODO: remove hot
	// is.True(!strings.Contains(string(code), `hot:`))

	code, err = fs.ReadFile(overlay, "bud/view/about/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"H2"`))
	is.True(strings.Contains(string(code), `"about"`))
	is.True(strings.Contains(string(code), fmt.Sprintf(`from"../%s"`, chunkName)))
	is.True(strings.Contains(string(code), `page:"/bud/view/about/index.svelte"`))
	is.True(strings.Contains(string(code), `document.getElementById("bud_target")`))
	// TODO: remove hot
	// is.True(!strings.Contains(string(code), `hot:`))

	code, err = fs.ReadFile(overlay, fmt.Sprintf("bud/view/%s", chunkName))
	is.NoErr(err)
	fmt.Println(string(code))
	is.True(strings.Contains(string(code), `"SvelteDOMInsert"`))
	is.True(strings.Contains(string(code), `"bud_props"`))
}

func TestImportLocal(t *testing.T) {
	t.SkipNow()
}

func TestImportNodeModule(t *testing.T) {
	t.SkipNow()
}
