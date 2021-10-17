package dom_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	"gitlab.com/mnm/bud/transform"

	"gitlab.com/mnm/bud/bfs"
	"gitlab.com/mnm/bud/dom"
	"gitlab.com/mnm/bud/internal/npm"
	"gitlab.com/mnm/bud/internal/vfs"
	v8 "gitlab.com/mnm/bud/js/v8"
	"gitlab.com/mnm/bud/svelte"
	"github.com/matryer/is"
)

func TestRunner(t *testing.T) {
	is := is.New(t)
	cwd, err := os.Getwd()
	is.NoErr(err)
	dir := filepath.Join(cwd, "_tmp")
	is.NoErr(os.RemoveAll(dir))
	defer func() {
		if !t.Failed() {
			is.NoErr(os.RemoveAll(dir))
		}
	}()
	memfs := vfs.Memory{
		"view/index.svelte": &fstest.MapFile{
			Data: []byte(`<h1>index</h1>`),
		},
		"view/about/index.svelte": &fstest.MapFile{
			Data: []byte(`<h2>about</h2>`),
		},
	}
	is.NoErr(vfs.WriteAll(".", dir, memfs))
	dirfs := os.DirFS(dir)
	svelteCompiler := svelte.New(&svelte.Input{
		VM:  v8.New(),
		Dev: true,
	})
	transformer := transform.MustLoad(
		svelte.NewTransformable(svelteCompiler),
	)
	bf := bfs.New(dirfs)
	bf.Add(map[string]bfs.Generator{
		"bud/view": dom.Runner(dir, transformer),
	})
	// Read the wrapped version of index.svelte with node_modules rewritten
	code, err := fs.ReadFile(bf, "bud/view/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h1");`))
	is.True(strings.Contains(string(code), `text("index")`))
	is.True(strings.Contains(string(code), `"/bud/view/index.svelte": view_default`))
	is.True(strings.Contains(string(code), `page: "/bud/view/index.svelte",`))
	is.True(strings.Contains(string(code), `hot: new Hot("/bud/hot?page=%2Fbud%2Fview%2Findex.svelte", components)`))

	// Unwrapped version with node_modules rewritten
	code, err = fs.ReadFile(bf, "bud/view/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h1");`))
	is.True(strings.Contains(string(code), `text("index")`))
	// Unwrapped version doesn't contain wrapping
	is.True(!strings.Contains(string(code), `"/bud/view/index.svelte": view_default`))
	is.True(!strings.Contains(string(code), `page: "/bud/view/index.svelte",`))
	is.True(!strings.Contains(string(code), `hot: new Hot("/bud/hot?page=%2Fbud%2Fview%2Findex.svelte", components)`))

	// Read the wrapped version of about/index.svelte with node_modules rewritten
	code, err = fs.ReadFile(bf, "bud/view/about/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h2");`))
	is.True(strings.Contains(string(code), `text("about")`))
	is.True(strings.Contains(string(code), `"/bud/view/about/index.svelte": about_default`))
	is.True(strings.Contains(string(code), `page: "/bud/view/about/index.svelte",`))
	is.True(strings.Contains(string(code), `hot: new Hot("/bud/hot?page=%2Fbud%2Fview%2Fabout%2Findex.svelte", components)`))

	// Unwrapped version with node_modules rewritten
	code, err = fs.ReadFile(bf, "bud/view/about/index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `from "/bud/node_modules/svelte/internal"`))
	is.True(strings.Contains(string(code), `element("h2");`))
	is.True(strings.Contains(string(code), `text("about")`))
	// Unwrapped version doesn't contain wrapping
	is.True(!strings.Contains(string(code), `"/bud/view/about/index.svelte": about_default`))
	is.True(!strings.Contains(string(code), `page: "/bud/view/about/index.svelte",`))
	is.True(!strings.Contains(string(code), `hot: new Hot("/bud/hot?page=%2Fbud%2Fview%2Fabout%2Findex.svelte", components)`))
}

func TestImportLocal(t *testing.T) {

}

func TestImportNodeModule(t *testing.T) {

}

func TestNodeModules(t *testing.T) {
	is := is.New(t)
	cwd, err := os.Getwd()
	is.NoErr(err)
	dir := filepath.Join(cwd, "_tmp")
	is.NoErr(os.RemoveAll(dir))
	defer func() {
		if !t.Failed() {
			is.NoErr(os.RemoveAll(dir))
		}
	}()
	memfs := vfs.Memory{
		"view/index.svelte": &fstest.MapFile{
			Data: []byte(`<h1>hi world</h1>`),
		},
	}
	is.NoErr(vfs.WriteAll(".", dir, memfs))
	dirfs := os.DirFS(dir)
	err = npm.Install(dir, "svelte@3.42.3")
	is.NoErr(err)
	bf := bfs.New(dirfs)
	bf.Add(map[string]bfs.Generator{
		"bud/node_modules": dom.NodeModules(dir),
	})
	// Read the re-written node_modules
	code, err := fs.ReadFile(bf, "bud/node_modules/svelte/internal")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `function element(`))
	is.True(strings.Contains(string(code), `function text(`))
}

func TestBuilder(t *testing.T) {
	chunkPath := "chunk-ZBJUH6TH.js"
	is := is.New(t)
	cwd, err := os.Getwd()
	is.NoErr(err)
	dir := filepath.Join(cwd, "_tmp")
	is.NoErr(os.RemoveAll(dir))
	defer func() {
		if !t.Failed() {
			is.NoErr(os.RemoveAll(dir))
		}
	}()
	memfs := vfs.Memory{
		"view/index.svelte": &fstest.MapFile{
			Data: []byte(`<h1>index</h1>`),
		},
		"view/about/index.svelte": &fstest.MapFile{
			Data: []byte(`<h2>about</h2>`),
		},
	}
	is.NoErr(vfs.WriteAll(".", dir, memfs))
	dirfs := os.DirFS(dir)
	err = npm.Install(dir, "svelte@3.42.3")
	is.NoErr(err)
	err = npm.Link("../budjs", dir)
	is.NoErr(err)
	svelteCompiler := svelte.New(&svelte.Input{
		VM:  v8.New(),
		Dev: true,
	})
	transformer := transform.MustLoad(
		svelte.NewTransformable(svelteCompiler),
	)
	bf := bfs.New(dirfs)
	bf.Add(map[string]bfs.Generator{
		"bud/view": dom.Builder(dir, transformer),
	})
	des, err := fs.ReadDir(bf, "bud/view")
	is.NoErr(err)
	is.Equal(len(des), 3)
	is.Equal(des[0].Name(), "_index.svelte")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[1].Name(), "about")
	is.Equal(des[1].IsDir(), true)
	is.Equal(des[2].Name(), chunkPath)
	is.Equal(des[2].IsDir(), false)
	des, err = fs.ReadDir(bf, "bud/view/about")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "_index.svelte")
	is.Equal(des[0].IsDir(), false)

	code, err := fs.ReadFile(bf, "bud/view/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"H1"`))
	is.True(strings.Contains(string(code), `"index"`))
	is.True(strings.Contains(string(code), fmt.Sprintf(`from"./%s"`, chunkPath)))
	is.True(strings.Contains(string(code), `page:"/bud/view/index.svelte"`))
	is.True(strings.Contains(string(code), `document.getElementById("bud_target")`))
	// TODO: remove hot
	// is.True(!strings.Contains(string(code), `hot:`))

	code, err = fs.ReadFile(bf, "bud/view/about/_index.svelte")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"H2"`))
	is.True(strings.Contains(string(code), `"about"`))
	is.True(strings.Contains(string(code), fmt.Sprintf(`from"../%s"`, chunkPath)))
	is.True(strings.Contains(string(code), `page:"/bud/view/about/index.svelte"`))
	is.True(strings.Contains(string(code), `document.getElementById("bud_target")`))
	// TODO: remove hot
	// is.True(!strings.Contains(string(code), `hot:`))

	code, err = fs.ReadFile(bf, fmt.Sprintf("bud/view/%s", chunkPath))
	is.NoErr(err)
	is.True(strings.Contains(string(code), `"allowpaymentrequest"`))
	is.True(strings.Contains(string(code), `"readonly"`))
}
