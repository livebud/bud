package cli_test

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testcli"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/testdir"
)

func TestNewControllerNoActions(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	err = td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(td.Directory())
	result, err := cli.Run(ctx, "new", "controller", "hello")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("controller/hello/controller.go"))
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 404 Not Found
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`))
	is.NoErr(app.Close())
}

func TestNewControllerNoActionsRoute(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	err = td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(td.Directory())
	result, err := cli.Run(ctx, "new", "controller", "hello:/")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("controller/controller.go"))
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 404 Not Found
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`))
	is.NoErr(app.Close())
}

func TestNewControllerAll(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	td.NodeModules["svelte"] = versions.Svelte
	err = td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(td.Directory())
	result, err := cli.Run(ctx, "new", "controller", "posts", "index", "show", "create", "update", "delete", "edit", "new")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("controller/posts/controller.go"))
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()

	// Post index
	res, err := app.Get("/posts")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	sel, err := res.Query("#bud_target")
	is.NoErr(err)
	html, err := sel.Html()
	is.NoErr(err)
	is.In(html, `<h1>Post Index</h1>`)
	is.In(html, `<table `)
	is.In(html, `</table>`)

	// New post
	res, err = app.Get("/posts/new")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.In(html, `<h1>New Post</h1>`)
	is.In(html, `<form method="post" action="/posts">`)
	is.In(html, `<input type="submit" value="Create Post"/>`)
	is.In(html, `</form>`)
	is.In(html, `<a href="/posts">Back</a>`)

	// Create post
	res, err = app.Post("/posts", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /posts/0
	`))

	// Edit post
	res, err = app.Get("/posts/10/edit")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.In(html, `<h1>Edit Post</h1>`)
	is.In(html, `<form method="post" action="/posts/10">`)
	is.In(html, `<input type="hidden" name="_method" value="patch"/>`)
	is.In(html, `<input type="submit" value="Update Post"/>`)
	is.In(html, `</form>`)
	is.In(html, `<a href="/posts">Back</a>`)

	// Update post with patch
	res, err = app.Patch("/posts/10", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /posts/10
	`))
	// Update post using method override
	values := url.Values{}
	values.Set("_method", http.MethodPatch)
	req, err := app.PostRequest("/posts/10", bytes.NewBufferString(values.Encode()))
	is.NoErr(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err = app.Do(req)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /posts/10
	`))

	// Delete post
	res, err = app.Delete("/posts/10", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /posts
	`))

	is.NoErr(app.Close())
}

func TestNewControllerAllRoot(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	td.NodeModules["svelte"] = versions.Svelte
	err = td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(td.Directory())
	result, err := cli.Run(ctx, "new", "controller", "posts:/", "index", "show", "create", "update", "delete", "edit", "new")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("controller/controller.go"))
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()

	// Post index
	res, err := app.Get("/")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	sel, err := res.Query("#bud_target")
	is.NoErr(err)
	html, err := sel.Html()
	is.NoErr(err)
	is.In(html, `<h1>Post Index</h1>`)
	is.In(html, `<table `)
	is.In(html, `</table>`)

	// New post
	res, err = app.Get("/new")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.In(html, `<h1>New Post</h1>`)
	is.In(html, `<form method="post" action="/">`)
	is.In(html, `<input type="submit" value="Create Post"/>`)
	is.In(html, `</form>`)
	is.In(html, `<a href="/">Back</a>`)

	// Create post
	res, err = app.Post("/", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /0
	`))

	// Edit post
	res, err = app.Get("/10/edit")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.In(html, `<h1>Edit Post</h1>`)
	is.In(html, `<form method="post" action="/10">`)
	is.In(html, `<input type="hidden" name="_method" value="patch"/>`)
	is.In(html, `<input type="submit" value="Update Post"/>`)
	is.In(html, `</form>`)
	is.In(html, `<a href="/">Back</a>`)

	// Update post with patch
	res, err = app.Patch("/10", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /10
	`))
	// Update post using method override
	values := url.Values{}
	values.Set("_method", http.MethodPatch)
	req, err := app.PostRequest("/10", bytes.NewBufferString(values.Encode()))
	is.NoErr(err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res, err = app.Do(req)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /10
	`))

	// Delete post
	res, err = app.Delete("/10", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /
	`))

	is.NoErr(app.Close())
}

func TestNewControllerAllNested(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	td.NodeModules["svelte"] = versions.Svelte
	err = td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(td.Directory())
	result, err := cli.Run(ctx, "new", "controller", "posts/comments", "index", "show", "create", "update", "delete", "edit", "new")
	is.True(err != nil)
	is.Equal(err.Error(), `new controller: scaffolding the "index" or "new" action of a nested resource like "posts/comments" isn't supported yet, see https://github.com/livebud/bud/issues/209 for details`)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	// is.NoErr(err)
	// is.Equal(result.Stdout(), "")
	// is.Equal(result.Stderr(), "")
	// is.NoErr(td.Exists("controller/posts/comments/controller.go"))
	// app, err := cli.Start(ctx, "run")
	// is.NoErr(err)
	// defer app.Close()

	// // Post index
	// res, err := app.Get("/posts/10/comments")
	// is.NoErr(err)
	// is.Equal(res.Status(), 200)
	// sel, err := res.Query("#bud_target")
	// is.NoErr(err)
	// html, err := sel.Html()
	// is.NoErr(err)
	// is.In(html, `<h1>Comment Index</h1>`)
	// is.In(html, `<table `)
	// is.In(html, `</table>`)

	// // New post
	// res, err = app.Get("/posts/10/comments/new")
	// is.NoErr(err)
	// is.Equal(res.Status(), 200)
	// sel, err = res.Query("#bud_target")
	// is.NoErr(err)
	// html, err = sel.Html()
	// is.NoErr(err)
	// is.In(html, `<h1>New Comment</h1>`)
	// is.In(html, `<form method="post" action="/posts/10/comments">`)
	// is.In(html, `<input type="submit" value="Create Comment"/>`)
	// is.In(html, `</form>`)
	// is.In(html, `<a href="/posts/10/comments">Back</a>`)

	// // Create post
	// res, err = app.Post("/posts/10/comments", nil)
	// is.NoErr(err)
	// is.NoErr(res.Diff(`
	// 	HTTP/1.1 302 Found
	// 	Location: /posts/10/comments/0
	// `))

	// // Edit post
	// res, err = app.Get("/posts/10/comments/5/edit")
	// is.NoErr(err)
	// is.Equal(res.Status(), 200)
	// sel, err = res.Query("#bud_target")
	// is.NoErr(err)
	// html, err = sel.Html()
	// is.NoErr(err)
	// is.In(html, `<h1>Edit Post</h1>`)
	// is.In(html, `<form method="post" action="/posts/10/comments/5">`)
	// is.In(html, `<input type="hidden" name="_method" value="patch"/>`)
	// is.In(html, `<input type="submit" value="Update Post"/>`)
	// is.In(html, `</form>`)
	// is.In(html, `<a href="/posts/10/comments">Back</a>`)

	// // Update post with patch
	// res, err = app.Patch("/posts/10/comments/5", nil)
	// is.NoErr(err)
	// is.NoErr(res.Diff(`
	// 	HTTP/1.1 302 Found
	// 	Location: /posts/10/comments/5
	// `))
	// // Update post using method override
	// values := url.Values{}
	// values.Set("_method", http.MethodPatch)
	// req, err := app.PostRequest("/posts/10/comments/5", bytes.NewBufferString(values.Encode()))
	// is.NoErr(err)
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// res, err = app.Do(req)
	// is.NoErr(err)
	// is.NoErr(res.Diff(`
	// 	HTTP/1.1 302 Found
	// 	Location: /posts/10/comments/5
	// `))

	// // Delete post
	// res, err = app.Delete("/posts/10/comments/5", nil)
	// is.NoErr(err)
	// is.NoErr(res.Diff(`
	// 	HTTP/1.1 302 Found
	// 	Location: /posts/10/comments
	// `))

	// is.NoErr(app.Close())
}

func TestNewControllerCustom(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	td.NodeModules["svelte"] = versions.Svelte
	err = td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(td.Directory())
	result, err := cli.Run(ctx, "new", "controller", "posts:/", "custom")
	is.True(err != nil)
	is.Equal(err.Error(), `new controller: invalid action "custom", expected "index", "new", "create", "show", "edit", "update" or "delete"`)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
}

func TestNewControllerRemoveViewDir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	td.NodeModules["svelte"] = versions.Svelte
	err = td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(td.Directory())
	result, err := cli.Run(ctx, "new", "controller", "/", "index")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("controller/controller.go"))
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.In(res.Body().String(), `<h1>Resource Index</h1>`)

	// Remove the view directory
	is.NoErr(td.RemoveAll("view"))

	// Wait for the app to be ready again
	readyCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	is.NoErr(app.Ready(readyCtx))
	cancel()

	res, err = app.Get("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		[]
	`))

	is.NoErr(app.Close())
}
