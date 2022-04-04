package ssr_test

import (
	"encoding/json"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/package/gomod"
	"gitlab.com/mnm/bud/package/js"
	v8 "gitlab.com/mnm/bud/package/js/v8"
	"gitlab.com/mnm/bud/package/overlay"
	"gitlab.com/mnm/bud/package/svelte"
	"gitlab.com/mnm/bud/runtime/transform"
	"gitlab.com/mnm/bud/runtime/view"
	"gitlab.com/mnm/bud/runtime/view/ssr"
)

func TestSvelteHello(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	td.Files["view/index.svelte"] = `<h1>hi world</h1>`
	td.NodeModules["svelte"] = "3.46.4"
	is.NoErr(td.Write(dir))
	vm, err := v8.Load()
	is.NoErr(err)
	svelteCompiler, err := svelte.Load(vm)
	is.NoErr(err)
	transformer := transform.MustLoad(svelte.NewTransformable(svelteCompiler))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	overlay, err := overlay.Load(module)
	is.NoErr(err)
	overlay.FileGenerator("bud/view/_ssr.js", ssr.New(module, transformer.SSR))
	// Read the wrapped version of index.svelte with node_modules rewritten
	code, err := fs.ReadFile(overlay, "bud/view/_ssr.js")
	is.NoErr(err)
	is.True(strings.Contains(string(code), `create_ssr_component(`))
	is.True(strings.Contains(string(code), `<h1>hi world</h1>`))
	is.True(strings.Contains(string(code), `views["/"] = `))
	result, err := vm.Eval("render.js", string(code)+`; bud.render("/", {})`)
	is.NoErr(err)
	var res view.Response
	err = json.Unmarshal([]byte(result), &res)
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.True(strings.Contains(res.Body, `<script id="bud_props" type="text/template" defer>{"props":{}}</script>`))
	is.True(strings.Contains(res.Body, `<script type="module" src="/bud/view/_index.svelte" defer></script>`))
	is.True(strings.Contains(res.Body, `<div id="bud_target">`))
	is.True(strings.Contains(res.Body, `<h1>hi world</h1>`))
}

func TestSvelteAwait(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("all good"))
	}))
	defer server.Close()
	td := testdir.New()
	td.Files["view/index.svelte"] = `
		<script>
			let promise = fetch("` + server.URL + `").then(res => res.text())
		</script>

		<div>
			{#await promise}
				Loading...
			{:then value}
				response: {value}
			{/await}
		</div>
	`
	td.NodeModules["svelte"] = "3.46.4"
	is.NoErr(td.Write(dir))
	vm, err := v8.Load()
	is.NoErr(err)
	svelteCompiler, err := svelte.Load(vm)
	is.NoErr(err)
	transformer := transform.MustLoad(svelte.NewTransformable(svelteCompiler))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	overlay, err := overlay.Load(module)
	is.NoErr(err)
	overlay.FileGenerator("bud/view/_ssr.js", ssr.New(module, transformer.SSR))
	// Read the wrapped version of index.svelte with node_modules rewritten
	code, err := fs.ReadFile(overlay, "bud/view/_ssr.js")
	is.NoErr(err)
	result, err := vm.Eval("render.js", string(code)+`; bud.render("/", {})`)
	is.NoErr(err)
	var res view.Response
	err = json.Unmarshal([]byte(result), &res)
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.True(strings.Contains(res.Body, `<script id="bud_props" type="text/template" defer>{"props":{}}</script>`))
	is.True(strings.Contains(res.Body, `<script type="module" src="/bud/view/_index.svelte" defer></script>`))
	is.True(strings.Contains(res.Body, `<div id="bud_target">`))
	is.True(strings.Contains(res.Body, `Loading...`))
}

func render(vm js.VM, code, path string, props interface{}) (*view.Response, error) {
	input, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}
	result, err := vm.Eval("render.js", string(code)+`; bud.render("`+path+`", `+string(input)+`)`)
	if err != nil {
		return nil, err
	}
	var res view.Response
	if err = json.Unmarshal([]byte(result), &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func TestSvelteProps(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	td.Files["view/index.svelte"] = `
		<script>
			export let props
		</script>
		<h1>{@html JSON.stringify(props)}</h1>
	`
	td.Files["view/show.svelte"] = `
		<script>
			export let props
		</script>
		<h2>{@html JSON.stringify(props)}</h2>
	`
	td.Files["view/users/index.svelte"] = `
		<script>
			export let props
		</script>
		<h3>{@html JSON.stringify(props)}</h3>
	`
	td.Files["view/users/show.svelte"] = `
		<script>
			export let props
		</script>
		<h4>{@html JSON.stringify(props)}</h4>
	`
	td.Files["view/posts/comments/index.svelte"] = `
		<script>
			export let props
		</script>
		<h5>{@html JSON.stringify(props)}</h5>
	`
	td.Files["view/posts/comments/show.svelte"] = `
		<script>
			export let props
		</script>
		<h6>{@html JSON.stringify(props)}</h6>
	`
	td.NodeModules["svelte"] = "3.46.4"
	is.NoErr(td.Write(dir))
	vm, err := v8.Load()
	is.NoErr(err)
	svelteCompiler, err := svelte.Load(vm)
	is.NoErr(err)
	transformer := transform.MustLoad(svelte.NewTransformable(svelteCompiler))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	overlay, err := overlay.Load(module)
	is.NoErr(err)
	overlay.FileGenerator("bud/view/_ssr.js", ssr.New(module, transformer.SSR))
	// Read the wrapped version of index.svelte with node_modules rewritten
	code, err := fs.ReadFile(overlay, "bud/view/_ssr.js")
	is.NoErr(err)
	// index
	type User struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	res, err := render(vm, string(code), "/", []*User{
		{"Alice", "alice@livebud.com"},
		{"Tom", "tom@livebud.com"},
	})
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.True(strings.Contains(res.Body, `<h1><!-- HTML_TAG_START -->[{"name":"Alice","email":"alice@livebud.com"},{"name":"Tom","email":"tom@livebud.com"}]<!-- HTML_TAG_END --></h1>`))
	// show
	res, err = render(vm, string(code), "/:id", &User{"Alice", "alice@livebud.com"})
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.True(strings.Contains(res.Body, `<h2><!-- HTML_TAG_START -->{"name":"Alice","email":"alice@livebud.com"}<!-- HTML_TAG_END --></h2>`))
	// users/index
	res, err = render(vm, string(code), "/users", []*User{
		{"Alice", "alice@livebud.com"},
		{"Tom", "tom@livebud.com"},
	})
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.True(strings.Contains(res.Body, `<h3><!-- HTML_TAG_START -->[{"name":"Alice","email":"alice@livebud.com"},{"name":"Tom","email":"tom@livebud.com"}]<!-- HTML_TAG_END --></h3>`))
	// users/show
	res, err = render(vm, string(code), "/users/:id", &User{"Alice", "alice@livebud.com"})
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.True(strings.Contains(res.Body, `<h4><!-- HTML_TAG_START -->{"name":"Alice","email":"alice@livebud.com"}<!-- HTML_TAG_END --></h4>`))
	// posts/comments/index
	type Comment struct {
		Name  string `json:"name"`
		Title string `json:"title"`
	}
	res, err = render(vm, string(code), "/posts/:post_id/comments", []*Comment{
		{"Alice", "first"},
		{"Tom", "second"},
	})
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.True(strings.Contains(res.Body, `<h5><!-- HTML_TAG_START -->[{"name":"Alice","title":"first"},{"name":"Tom","title":"second"}]<!-- HTML_TAG_END --></h5>`))
	// posts/comments/:id
	res, err = render(vm, string(code), "/posts/:post_id/comments/:id", &Comment{"Alice", "first"})
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.True(strings.Contains(res.Body, `<h6><!-- HTML_TAG_START -->{"name":"Alice","title":"first"}<!-- HTML_TAG_END --></h6>`))
}

func TestSvelteLocalImports(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	td := testdir.New()
	td.Files["view/Comment.svelte"] = `
		<script>
			export let comment = {}
		</script>
		<h2>{comment.message}</h2>
	`
	td.Files["view/Story.svelte"] = `
		<script>
			export let story = {}
		</script>
		<h1>{story.title}</h1>
	`
	td.Files["view/show.svelte"] = `
		<script>
			import Story from "./Story.svelte"
			import Comment from "./Comment.svelte"
			export let props = {
					comments: []
			}
		</script>
		<Story story={props} />
		{#each props.comments as comment}
			<Comment {comment} />
		{/each}
	`
	td.NodeModules["svelte"] = "3.46.4"
	is.NoErr(td.Write(dir))
	vm, err := v8.Load()
	is.NoErr(err)
	svelteCompiler, err := svelte.Load(vm)
	is.NoErr(err)
	transformer := transform.MustLoad(svelte.NewTransformable(svelteCompiler))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	overlay, err := overlay.Load(module)
	is.NoErr(err)
	overlay.FileGenerator("bud/view/_ssr.js", ssr.New(module, transformer.SSR))
	// Read the wrapped version of index.svelte with node_modules rewritten
	code, err := fs.ReadFile(overlay, "bud/view/_ssr.js")
	is.NoErr(err)
	type Comment struct {
		Message string `json:"message"`
	}
	type Story struct {
		Title    string     `json:"title"`
		Comments []*Comment `json:"comments"`
	}
	res, err := render(vm, string(code), "/:id", &Story{
		Title: "first story",
		Comments: []*Comment{
			{Message: "first comment"},
			{Message: "second comment"},
		},
	})
	is.NoErr(err)
	is.Equal(res.Status, 200)
	is.Equal(len(res.Headers), 1)
	is.Equal(res.Headers["Content-Type"], "text/html")
	is.True(strings.Contains(res.Body, `<h1>first story</h1>`))
	is.True(strings.Contains(res.Body, `<h2>first comment</h2><h2>second comment</h2>`))
}
