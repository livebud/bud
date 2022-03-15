package dom_test

import (
	"fmt"
	"io/fs"
	"strings"
	"testing"

	"gitlab.com/mnm/bud/package/overlay"

	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/runtime/transform"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
	v8 "gitlab.com/mnm/bud/pkg/js/v8"
	"gitlab.com/mnm/bud/pkg/svelte"
	"gitlab.com/mnm/bud/runtime/view/dom"
)

func TestRunner(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	svelteCompiler := svelte.New(v8.New())
	transformer := transform.MustLoad(
		svelte.NewTransformable(svelteCompiler),
	)
	td := testdir.New()
	td.Files["view/index.svelte"] = `<h1>index</h1>`
	td.Files["view/about/index.svelte"] = `<h2>about</h2>`
	is.NoErr(td.Write(dir))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	overlay, err := overlay.Load(module)
	is.NoErr(err)
	overlay.FileServer("bud/view", dom.Runner(overlay, module, transformer))
	// Read the wrapped version of index.svelte with node_modules rewritten
	code, err := fs.ReadFile(overlay, "bud/view/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h1");`))
	is.True(strings.Contains(string(code), `text("index")`))
	is.True(strings.Contains(string(code), `"/bud/view/index.svelte": view_default`))
	is.True(strings.Contains(string(code), `page: "/bud/view/index.svelte",`))
	is.True(strings.Contains(string(code), `hot: new Hot("http://0.0.0.0:35729/?page=%2Fbud%2Fview%2Findex.svelte", components)`))

	// Unwrapped version with node_modules rewritten
	code, err = fs.ReadFile(overlay, "bud/view/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h1");`))
	is.True(strings.Contains(string(code), `text("index")`))
	// Unwrapped version doesn't contain wrapping
	is.True(!strings.Contains(string(code), `"/bud/view/index.svelte": view_default`))
	is.True(!strings.Contains(string(code), `page: "/bud/view/index.svelte",`))
	is.True(!strings.Contains(string(code), `hot: new Hot("http://0.0.0.0:35729/?page=%2Fbud%2Fview%2Findex.svelte", components)`))

	// Read the wrapped version of about/index.svelte with node_modules rewritten
	code, err = fs.ReadFile(overlay, "bud/view/about/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h2");`))
	is.True(strings.Contains(string(code), `text("about")`))
	is.True(strings.Contains(string(code), `"/bud/view/about/index.svelte": about_default`))
	is.True(strings.Contains(string(code), `page: "/bud/view/about/index.svelte",`))
	is.True(strings.Contains(string(code), `hot: new Hot("http://0.0.0.0:35729/?page=%2Fbud%2Fview%2Fabout%2Findex.svelte", components)`))

	// Unwrapped version with node_modules rewritten
	code, err = fs.ReadFile(overlay, "bud/view/about/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h2");`))
	is.True(strings.Contains(string(code), `text("about")`))
	// Unwrapped version doesn't contain wrapping
	is.True(!strings.Contains(string(code), `"/bud/view/about/index.svelte": about_default`))
	is.True(!strings.Contains(string(code), `page: "/bud/view/about/index.svelte",`))
	is.True(!strings.Contains(string(code), `hot: new Hot("http://0.0.0.0:35729/?page=%2Fbud%2Fview%2Fabout%2Findex.svelte", components)`))
}

func TestNodeModules(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	td.Files["view/index.svelte"] = `<h1>hi world</h1>`
	td.NodeModules["svelte"] = "3.42.3"
	is.NoErr(td.Write(dir))
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

func TestBuilder(t *testing.T) {
	chunkPath := "chunk-H7BRTJPS.js"
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	td.Files["view/index.svelte"] = `<h1>index</h1>`
	td.Files["view/about/index.svelte"] = `<h2>about</h2>`
	td.NodeModules["livebud"] = "*"
	td.NodeModules["svelte"] = "3.46.4"
	is.NoErr(td.Write(dir))
	svelteCompiler := svelte.New(v8.New())
	transformer := transform.MustLoad(svelte.NewTransformable(svelteCompiler))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	overlay, err := overlay.Load(module)
	is.NoErr(err)
	overlay.DirGenerator("bud/view", dom.Builder(overlay, module, transformer))
	des, err := fs.ReadDir(overlay, "bud/view")
	is.NoErr(err)
	is.Equal(len(des), 3)
	is.Equal(des[0].Name(), "_index.svelte")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[1].Name(), "about")
	is.Equal(des[1].IsDir(), true)
	is.Equal(des[2].Name(), chunkPath)
	is.Equal(des[2].IsDir(), false)
	des, err = fs.ReadDir(overlay, "bud/view/about")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "_index.svelte")
	is.Equal(des[0].IsDir(), false)

	code, err := fs.ReadFile(overlay, "bud/view/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"H1"`))
	is.True(strings.Contains(string(code), `"index"`))
	is.True(strings.Contains(string(code), fmt.Sprintf(`from"./%s"`, chunkPath)))
	is.True(strings.Contains(string(code), `page:"/bud/view/index.svelte"`))
	is.True(strings.Contains(string(code), `document.getElementById("bud_target")`))
	// TODO: remove hot
	// is.True(!strings.Contains(string(code), `hot:`))

	code, err = fs.ReadFile(overlay, "bud/view/about/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"H2"`))
	is.True(strings.Contains(string(code), `"about"`))
	is.True(strings.Contains(string(code), fmt.Sprintf(`from"../%s"`, chunkPath)))
	is.True(strings.Contains(string(code), `page:"/bud/view/about/index.svelte"`))
	is.True(strings.Contains(string(code), `document.getElementById("bud_target")`))
	// TODO: remove hot
	// is.True(!strings.Contains(string(code), `hot:`))

	code, err = fs.ReadFile(overlay, fmt.Sprintf("bud/view/%s", chunkPath))
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"SvelteDOMInsert"`))
	is.True(strings.Contains(string(code), `"bud_props"`))
}

func TestImportLocal(t *testing.T) {
	t.SkipNow()
}

func TestImportNodeModule(t *testing.T) {
	t.SkipNow()
}
