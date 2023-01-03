package dom_test

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/livebud/bud/package/genfs"

	"github.com/livebud/bud/framework/view/nodemodules"

	"github.com/livebud/bud/package/log/testlog"

	"github.com/livebud/bud/framework/transform/transformrt"
	"github.com/livebud/bud/package/gomod"

	"github.com/livebud/bud/framework/view/dom"
	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
	v8 "github.com/livebud/bud/package/js/v8"
	"github.com/livebud/bud/package/svelte"
)

func TestServeFile(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx := context.Background()
	dir := t.TempDir()
	vm, err := v8.Load()
	is.NoErr(err)
	svelteCompiler, err := svelte.Load(vm)
	is.NoErr(err)
	transformer, err := transformrt.Default(log, svelteCompiler)
	is.NoErr(err)
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>index</h1>`
	td.Files["view/about/index.svelte"] = `<h2>about</h2>`
	td.NodeModules["svelte"] = versions.Svelte
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	gfs := genfs.New(dag.Discard, module, log)
	gfs.FileServer("bud/view", dom.New(module, transformer))
	// Read the wrapped version of index.svelte with node_modules rewritten
	code, err := fs.ReadFile(gfs, "bud/view/_index.svelte.js")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h1");`))
	is.True(strings.Contains(string(code), `text("index")`))
	is.True(strings.Contains(string(code), `"/bud/view/index.svelte": view_default`))
	is.True(strings.Contains(string(code), `page: "/bud/view/index.svelte",`))
	is.True(strings.Contains(string(code), `hot: new Hot("http://127.0.0.1:35729/bud/hot/view/index.svelte", components)`))

	// Unwrapped version with node_modules rewritten
	code, err = fs.ReadFile(gfs, "bud/view/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h1");`))
	is.True(strings.Contains(string(code), `text("index")`))
	// Unwrapped version doesn't contain wrapping
	is.True(!strings.Contains(string(code), `"/bud/view/index.svelte": view_default`))
	is.True(!strings.Contains(string(code), `page: "/bud/view/index.svelte",`))
	is.True(!strings.Contains(string(code), `hot: new Hot("http://127.0.0.1:35729/bud/hot/view/index.svelte", components)`))

	// Read the wrapped version of about/index.svelte with node_modules rewritten
	code, err = fs.ReadFile(gfs, "bud/view/about/_index.svelte.js")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h2");`))
	is.True(strings.Contains(string(code), `text("about")`))
	is.True(strings.Contains(string(code), `"/bud/view/about/index.svelte": about_default`))
	is.True(strings.Contains(string(code), `page: "/bud/view/about/index.svelte",`))
	is.True(strings.Contains(string(code), `hot: new Hot("http://127.0.0.1:35729/bud/hot/view/about/index.svelte", components)`))

	// Unwrapped version with node_modules rewritten
	code, err = fs.ReadFile(gfs, "bud/view/about/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h2");`))
	is.True(strings.Contains(string(code), `text("about")`))
	// Unwrapped version doesn't contain wrapping
	is.True(!strings.Contains(string(code), `"/bud/view/about/index.svelte": about_default`))
	is.True(!strings.Contains(string(code), `page: "/bud/view/about/index.svelte",`))
	is.True(!strings.Contains(string(code), `hot: new Hot("http://127.0.0.1:35729/bud/hot/view/about/index.svelte", components)`))
}

func TestNodeModules(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>hi world</h1>`
	td.NodeModules["svelte"] = versions.Svelte
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	gfs := genfs.New(dag.Discard, module, log)
	gfs.FileServer("bud/node_modules", nodemodules.New(module))
	// Read the re-written node_modules
	code, err := fs.ReadFile(gfs, "bud/node_modules/svelte/internal")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `function element(`))
	is.True(strings.Contains(string(code), `function text(`))
}

func TestGenerateDir(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["view/index.svelte"] = `<h1>index</h1>`
	td.Files["view/about/index.svelte"] = `<h2>about</h2>`
	td.NodeModules["livebud"] = "*"
	td.NodeModules["svelte"] = versions.Svelte
	is.NoErr(td.Write(ctx))
	vm, err := v8.Load()
	is.NoErr(err)
	svelteCompiler, err := svelte.Load(vm)
	is.NoErr(err)
	transformer, err := transformrt.Default(log, svelteCompiler)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	gfs := genfs.New(dag.Discard, module, log)
	gfs.DirGenerator("bud/view", dom.New(module, transformer))
	des, err := fs.ReadDir(gfs, "bud/view")
	is.NoErr(err)
	is.Equal(len(des), 3)
	is.Equal(des[0].Name(), "_index.svelte.js")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[1].Name(), "about")
	is.Equal(des[1].IsDir(), true)
	is.True(strings.HasPrefix(des[2].Name(), "chunk-"))
	is.Equal(des[2].IsDir(), false)
	chunkName := des[2].Name()
	des, err = fs.ReadDir(gfs, "bud/view/about")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "_index.svelte.js")
	is.Equal(des[0].IsDir(), false)

	code, err := fs.ReadFile(gfs, "bud/view/_index.svelte.js")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"H1"`))
	is.True(strings.Contains(string(code), `"index"`))
	is.True(strings.Contains(string(code), fmt.Sprintf(`from"./%s"`, chunkName)))
	is.True(strings.Contains(string(code), `page:"/bud/view/index.svelte"`))
	is.True(strings.Contains(string(code), `document.getElementById("bud_target")`))
	// TODO: remove hot
	// is.True(!strings.Contains(string(code), `hot:`))

	code, err = fs.ReadFile(gfs, "bud/view/about/_index.svelte.js")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"H2"`))
	is.True(strings.Contains(string(code), `"about"`))
	is.True(strings.Contains(string(code), fmt.Sprintf(`from"../%s"`, chunkName)))
	is.True(strings.Contains(string(code), `page:"/bud/view/about/index.svelte"`))
	is.True(strings.Contains(string(code), `document.getElementById("bud_target")`))
	// TODO: remove hot
	// is.True(!strings.Contains(string(code), `hot:`))

	code, err = fs.ReadFile(gfs, fmt.Sprintf("bud/view/%s", chunkName))
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"SvelteDOMInsert"`))
	is.True(strings.Contains(string(code), `"bud_props"`))
}

func TestUpdateFile(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	ctx := context.Background()
	dir := t.TempDir()
	vm, err := v8.Load()
	is.NoErr(err)
	svelteCompiler, err := svelte.Load(vm)
	is.NoErr(err)
	transformer, err := transformrt.Default(log, svelteCompiler)
	is.NoErr(err)
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.Files["view/Story.svelte"] = `<h2>Story</h2>`
	td.Files["view/index.svelte"] = `
		<script>
			import Story from "./Story.svelte";
		</script>
		<h1>home</h1>
		<Story />
	`
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	cache, err := dag.Load(log, module.Directory("bud/bud.db"))
	is.NoErr(err)
	gfs := genfs.New(cache, module, log)
	gfs.FileServer("bud/view", dom.New(module, transformer))
	// check entry
	code, err := fs.ReadFile(gfs, "bud/view/_index.svelte.js")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"home"`), "missing home")
	is.True(strings.Contains(string(code), `"Story"`), "missing Story")
	// check component
	code, err = fs.ReadFile(gfs, "bud/view/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"home"`), "missing home")
	is.True(strings.Contains(string(code), `"Story"`), "missing Story")
	// Change view/Story.svelte and view/index.svelte
	os.WriteFile(filepath.Join(dir, "view/Story.svelte"), []byte(`<h2>Stories</h2>`), 0644)
	os.WriteFile(filepath.Join(dir, "view/index.svelte"), []byte(`
		<script>
			import Story from "./Story.svelte";
		</script>
		<h1>homies</h1>
		<Story />
	`), 0644)
	// check entry (cached)
	code, err = fs.ReadFile(gfs, "bud/view/_index.svelte.js")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"home"`), "missing home")
	is.True(strings.Contains(string(code), `"Story"`), "missing Story")
	// check component (cached)
	code, err = fs.ReadFile(gfs, "bud/view/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"home"`), "missing home")
	is.True(strings.Contains(string(code), `"Story"`), "missing Story")
	// Mark view/Story.svelte and view/index.svelte as changed
	is.NoErr(cache.Delete("view/index.svelte", "view/Story.svelte"))
	// check entry (uncached)
	code, err = fs.ReadFile(gfs, "bud/view/_index.svelte.js")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"homies"`), "missing homies")
	is.True(strings.Contains(string(code), `"Stories"`), "missing Stories")
	// check component (uncached)
	code, err = fs.ReadFile(gfs, "bud/view/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"homies"`), "missing homies")
	is.True(strings.Contains(string(code), `"Stories"`), "missing Stories")

	// Remove a file
	is.NoErr(os.Remove(filepath.Join(dir, "view/index.svelte")))
	// check page (cached)
	code, err = fs.ReadFile(gfs, "bud/view/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"homies"`), "missing homies")
	is.True(strings.Contains(string(code), `"Stories"`), "missing Stories")
	// Mark view/index.svelte  as changed
	is.NoErr(cache.Delete("view/index.svelte"))
	// check page (uncached)
	code, err = fs.ReadFile(gfs, "bud/view/index.svelte")
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(code, nil)
}
