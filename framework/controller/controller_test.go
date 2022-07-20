package controller_test

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/lithammer/dedent"
	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
	"github.com/matthewmueller/diff"
)

func TestNoActions(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// HTML response
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 404 Not Found
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff
	`))
	is.NoErr(app.Close())
}

func TestIndexString(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string {
			return "Root"
		}
	`
	td.Files["controller/about/controller.go"] = `
		package about
		type Controller struct {}
		func (c *Controller) Index() string {
			return "About"
		}
	`
	td.Files["controller/posts/comments/controller.go"] = `
		package comments
		type Controller struct {}
		func (c *Controller) Index() string {
			return "Comments"
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// HTML response
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: text/html
	`))
	is.In(res.Body().String(), "Root")
	// JSON response
	res, err = app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"Root"
	`))
	// HTML response
	res, err = app.Get("/about")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: text/html
	`))
	is.In(res.Body().String(), "About")
	// JSON response
	res, err = app.GetJSON("/about")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"About"
	`))
	// HTML response
	res, err = app.Get("/posts/some-slug/comments")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: text/html
	`))
	is.In(res.Body().String(), "Comments")
	// JSON response
	res, err = app.GetJSON("/posts/some-slug/comments")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"Comments"
	`))
	is.NoErr(app.Close())
}

func TestCreateRedirect(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Create() {}
	`
	td.Files["controller/users/controller.go"] = `
		package users
		type Controller struct {}
		func (c *Controller) Create() {}
	`
	td.Files["controller/posts/comments/controller.go"] = `
		package comments
		type Controller struct {}
		func (c *Controller) Create() {}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Redirect /
	res, err := app.Post("/", nil)
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 302 Found
		Location: /
	`))
	is.Equal(res.Body().Len(), 0)
	// Redirect /users
	res, err = app.Post("/users", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /users
	`))
	is.Equal(res.Body().Len(), 0)
	// Redirect /posts/10/comments
	res, err = app.Post("/posts/10/comments", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /posts/10/comments
	`))
	is.Equal(res.Body().Len(), 0)
	is.NoErr(app.Close())
}

func TestNoContent(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() {}
		func (c *Controller) Show(id string) {}
		func (c *Controller) New() {}
		func (c *Controller) Edit(id int) {}
	`
	td.Files["controller/posts/comments/controller.go"] = `
		package comments
		type Controller struct {}
		func (c *Controller) Index(postId int) error { return nil }
		func (c *Controller) Show(postId int, id string) error { return nil }
		func (c *Controller) New(postId int) error { return nil }
		func (c *Controller) Edit(postId int, id int) error { return nil }
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Root
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.Get("/10")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.Get("/new")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.Get("/10/edit")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.GetJSON("/10")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.GetJSON("/new")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.GetJSON("/10/edit")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	// Comments
	res, err = app.Get("/posts/1/comments")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.Get("/posts/1/comments/5")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.Get("/posts/1/comments/new")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.Get("/posts/1/comments/10/edit")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.GetJSON("/posts/1/comments")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.GetJSON("/posts/1/comments/10")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.GetJSON("/posts/1/comments/new")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.GetJSON("/posts/1/comments/10/edit")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	is.NoErr(app.Close())
}

func TestIndex500(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		import "errors"
		type Controller struct {}
		type Post struct {}
		func (c *Controller) Index() ([]*Post, error) {
			return nil, errors.New("unable to list posts")
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json

		{"error":"unable to list posts"}
	`))
	is.NoErr(app.Close())
}

func TestIndexList500(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		import "errors"
		type Controller struct {}
		type Post struct {}
		func (c *Controller) Index() (int, string, error) {
			return 0, "", errors.New("unable to list posts")
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /
	`))
	is.NoErr(app.Close())
}

func TestIndexList200(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		type Post struct {}
		func (c *Controller) Index() (int, string, error) {
			return 0, "a", nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /
	`))
	is.NoErr(app.Close())
}

func TestIndexListObject500(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		import "errors"
		type Controller struct {}
		type Post struct {}
		func (c *Controller) Index() (a int, b string, err error) {
			return 0, "a", errors.New("unable to list posts")
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /
	`))
	is.NoErr(app.Close())
}

func TestIndexListObject200(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		type Post struct {}
		func (c *Controller) Index() (a int, b string, err error) {
			return 0, "a", nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /
	`))
	is.NoErr(app.Close())
}

func TestIndexStructs200(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		type Post struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
		func (c *Controller) Index() (list []*Post, err error) {
			return []*Post{{0, "a"}, {1, "b"}}, nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":0,"title":"a"},{"id":1,"title":"b"}]
	`))
	is.NoErr(app.Close())
}

func TestJSONCreate204(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Create() {}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PostJSON("/", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	is.NoErr(app.Close())
}

func TestJSONCreate500(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		import "errors"
		type Controller struct {}
		func (c *Controller) Create() (string, error) {
			return "", errors.New("Not implemented yet")
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PostJSON("/", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json

		{"error":"Not implemented yet"}
	`))
	is.NoErr(app.Close())
}

func TestDependencyHoist(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["postgres/pool.go"] = `
		package postgres
		func New() *Pool { return &Pool{1} }
		type Pool struct { ID int }
	`
	td.Files["controller/controller.go"] = `
		package controller
		import "app.com/postgres"
		type Controller struct {
			Pool *postgres.Pool
		}
		func (c *Controller) Index() int {
			return c.Pool.ID
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		1
	`))
	is.NoErr(app.Close())
}

func TestDependencyRequest(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["postgres/pool.go"] = `
		package postgres
		import "net/http"
		func New(r *http.Request) *Pool { return &Pool{r.URL.Path} }
		type Pool struct { Path string }
	`
	td.Files["controller/controller.go"] = `
		package controller
		import "app.com/postgres"
		type Controller struct {
			Pool *postgres.Pool
		}
		func (c *Controller) Index() string {
			return c.Pool.Path
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: text/html
	`))
	is.In(res.Body().String(), `/`)
	is.NoErr(app.Close())
}

func TestShareStruct(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["article/article.go"] = `
		package article
		type Article struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
	`
	td.Files["controller/controller.go"] = `
		package controller
		import "app.com/article"
		type Controller struct {}
		func (c *Controller) Update(a *article.Article) (*article.Article, error) {
			return a, nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PatchJSON("/10", bytes.NewBufferString(`{"id": 1, "title": "a"}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"title":"a"}
	`))
	is.NoErr(app.Close())
}

func TestJSONCreateNested(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["postgres/pool.go"] = `
		package postgres
		func New(r *http.Request) *Pool { return &Pool{r.URL.Path} }
		type Pool struct { Path string }
	`
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		type Post struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
		func (c *Controller) Create(p *Post) *Post {
			return p
		}
	`
	td.Files["controller/users/users.go"] = `
		package users
		type Controller struct {}
		type Post struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
		func (c *Controller) Create(p *Post) *Post {
			return p
		}
	`
	td.Files["controller/users/admin/admin.go"] = `
		package admin
		type Controller struct {}
		type Post struct {
			UserID int ` + "`" + `json:"user_id"` + "`" + `
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
		func (c *Controller) Create(p *Post) *Post {
			return p
		}
	`
	td.Files["controller/articles/articles.go"] = `
		package articles
		type Controller struct {}
		type Post struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
		func (c *Controller) Create(p *Post) *Post {
			return p
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PostJSON("/", bytes.NewBufferString(`{"id": 1, "title": "a"}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"title":"a"}
	`))
	res, err = app.PostJSON("/users", bytes.NewBufferString(`{"id": 2, "title": "b"}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"title":"b"}
	`))
	res, err = app.PostJSON("/users/1/admin", bytes.NewBufferString(`{"id": 3, "title": "c"}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"user_id":1,"id":3,"title":"c"}
	`))
	res, err = app.PostJSON("/articles", bytes.NewBufferString(`{"id": 4, "title": "d"}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":4,"title":"d"}
	`))
	is.NoErr(app.Close())
}

func TestJSONDelete500(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		import "errors"
		import "context"
		type Controller struct {}
		func (c *Controller) Delete(ctx context.Context, id int) (err error) {
			return errors.New("Not implemented yet")
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.DeleteJSON("/1", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json

		{"error":"Not implemented yet"}
	`))
	is.NoErr(app.Close())
}

func TestJSONDelete200(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		type Post struct {
			ID int
			Title string
		}
		func (c *Controller) Delete(id int) (*Post, error) {
			return &Post{id, "a"}, nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.DeleteJSON("/1", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"ID":1,"Title":"a"}
	`))
	is.NoErr(app.Close())
}

func TestJSONMultipleActions(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string {
			return "hello world"
		}
		// Show route
		func (c *Controller) Show(id int) int {
			return id
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"hello world"
	`))
	res, err = app.GetJSON("/10")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		10
	`))
	is.NoErr(app.Close())
}

func TestJSONUpdate500(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		import "errors"
		type Controller struct {}
		func (c *Controller) Update() (string, error) {
			return "", errors.New("Not implemented yet")
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PatchJSON("/1", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json

		{"error":"Not implemented yet"}
	`))
	is.NoErr(app.Close())
}

func TestJSONUpdate200(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		type Post struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
		func (c *Controller) Update(id int) (*Post, error) {
			return &Post{id, "a"}, nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PatchJSON("/1", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"title":"a"}
	`))
	is.NoErr(app.Close())
}

func TestNestedResource(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/users/users.go"] = `
		package users
		type DB struct {}
		type Controller struct {
			DB *DB
		}
		type User struct {
			ID   int ` + "`" + `json:"id"` + "`" + `
			Name string ` + "`" + `json:"name"` + "`" + `
			Age  int ` + "`" + `json:"age"` + "`" + `
		}
		func (c *Controller) Index() ([]*User, error) {
			return []*User{{1, "a", 2}, {2, "b", 3}}, nil
		}
		func (c *Controller) New() {}
		func (c *Controller) Create(name string, age int) (*User, error) {
			return &User{3, name, age}, nil
		}
		func (c *Controller) Show(id int) (*User, error) {
			return &User{id, "d", 5}, nil
		}
		func (c *Controller) Edit(id int) {}
		func (c *Controller) Update(id int, name *string, age *int) error {
			return nil
		}
		func (c *Controller) Delete(id int) error {
			return nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/users")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":1,"name":"a","age":2},{"id":2,"name":"b","age":3}]
	`))
	res, err = app.GetJSON("/users/new")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.PostJSON("/users?name=matt&age=10", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":3,"name":"matt","age":10}
	`))
	res, err = app.GetJSON("/users/10")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"d","age":5}
	`))
	res, err = app.GetJSON("/users/10/edit")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.PatchJSON("/users/10", bytes.NewBufferString(`{"name": "matt", "age": 10}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.DeleteJSON("/users/10", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	is.NoErr(app.Close())
}

func TestDeepNestedResource(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/posts/comments/comments.go"] = `
		package comments
		type DB struct {}
		type Controller struct {
			DB *DB
		}
		type Comment struct {
			ID     int ` + "`" + `json:"id,omitempty"` + "`" + `
			PostID int ` + "`" + `json:"post_id,omitempty"` + "`" + `
			Title  string ` + "`" + `json:"title,omitempty"` + "`" + `
		}
		func (c *Controller) Index(postID int) ([]*Comment, error) {
			return []*Comment{{2, postID, "a"}, {3, postID, "b"}}, nil
		}
		func (c *Controller) New(postID int) {}
		func (c *Controller) Create(postID int, title string) (*Comment, error) {
			return &Comment{1, postID, title}, nil
		}
		func (c *Controller) Show(postID, id int) (*Comment, error) {
			return &Comment{id, postID, "a"}, nil
		}
		func (c *Controller) Edit(postID, id int) {}
		func (c *Controller) Update(postID, id int, title *string) (*Comment, error) {
			if title == nil {
				return &Comment{postID, id, ""}, nil
			}
			return &Comment{postID, id, *title}, nil
		}
		func (c *Controller) Delete(postID, id int) (*Comment, error) {
			return &Comment{postID, id, ""}, nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/posts/1/comments")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":2,"post_id":1,"title":"a"},{"id":3,"post_id":1,"title":"b"}]
	`))
	res, err = app.GetJSON("/posts/1/comments/new")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.PostJSON("/posts/1/comments", bytes.NewBufferString(`{"title":"1st"}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"post_id":1,"title":"1st"}
	`))
	res, err = app.GetJSON("/posts/1/comments/2")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"post_id":1,"title":"a"}
	`))
	res, err = app.GetJSON("/posts/1/comments/2/edit")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
	`))
	res, err = app.PatchJSON("/posts/1/comments/2", bytes.NewBufferString(`{"title":"1st"}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"post_id":2,"title":"1st"}
	`))
	res, err = app.PatchJSON("/posts/1/comments/2", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"post_id":2}
	`))
	res, err = app.DeleteJSON("/posts/1/comments/2", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"post_id":2}
	`))
	is.NoErr(app.Close())
}

func TestRedirectResource(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		type Post struct {
			ID     int ` + "`" + `json:"id,omitempty"` + "`" + `
			Title  string ` + "`" + `json:"title,omitempty"` + "`" + `
		}
		func (c *Controller) Create(title string) (*Post, error) {
			return &Post{2, title}, nil
		}
		func (c *Controller) Update(id int, title *string) (*Post, error) {
			if title == nil {
				return &Post{id, ""}, nil
			}
			return &Post{id, *title}, nil
		}
		func (c *Controller) Delete(id int) (*Post, error) {
			return &Post{id, ""}, nil
		}
	`
	td.Files["controller/posts/posts.go"] = `
		package posts
		type Controller struct {}
		type Post struct {
			ID     int ` + "`" + `json:"id,omitempty"` + "`" + `
			Title  string ` + "`" + `json:"title,omitempty"` + "`" + `
		}
		func (c *Controller) Create(title string) (*Post, error) {
			return &Post{2, title}, nil
		}
		func (c *Controller) Update(id int, title *string) (*Post, error) {
			if title == nil {
				return &Post{id, ""}, nil
			}
			return &Post{id, *title}, nil
		}
		func (c *Controller) Delete(id int) (*Post, error) {
			return &Post{id, ""}, nil
		}
	`
	td.Files["controller/posts/comments/comments.go"] = `
		package comments
		type DB struct {}
		type Controller struct {
			DB *DB
		}
		type Comment struct {
			ID     int ` + "`" + `json:"id,omitempty"` + "`" + `
			PostID int ` + "`" + `json:"post_id,omitempty"` + "`" + `
			Title  string ` + "`" + `json:"title,omitempty"` + "`" + `
		}
		func (c *Controller) Index(postID int) ([]*Comment, error) {
			return []*Comment{{2, postID, "a"}, {3, postID, "b"}}, nil
		}
		func (c *Controller) New(postID int) {}
		func (c *Controller) Create(postID int, title string) (*Comment, error) {
			return &Comment{2, postID, title}, nil
		}
		func (c *Controller) Show(postID, id int) (*Comment, error) {
			return &Comment{id, postID, "a"}, nil
		}
		func (c *Controller) Edit(postID, id int) {}
		func (c *Controller) Update(postID, id int, title *string) (*Comment, error) {
			if title == nil {
				return &Comment{postID, id, ""}, nil
			}
			return &Comment{postID, id, *title}, nil
		}
		func (c *Controller) Delete(postID, id int) (*Comment, error) {
			return &Comment{postID, id, ""}, nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Root
	res, err := app.Post("/", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /2
	`))
	res, err = app.Patch("/1", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /1
	`))
	res, err = app.Delete("/1", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /
	`))
	// Posts
	res, err = app.Post("/posts", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /posts/2
	`))
	res, err = app.Patch("/posts/1", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /posts/1
	`))
	res, err = app.Delete("/posts/1", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /posts
	`))
	// Comments
	res, err = app.Post("/posts/1/comments", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /posts/1/comments/2
	`))
	res, err = app.Patch("/posts/1/comments/2", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /posts/1/comments/2
	`))
	res, err = app.Delete("/posts/1/comments/2", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 302 Found
		Location: /posts/1/comments
	`))
	is.NoErr(app.Close())
}

func TestViewUnnamed(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.Files["view/index.svelte"] = `
		<script>
			export let users = []
		</script>
		{#each users as user}
		<h1>index: {user.id} {user.name}</h1>
		{/each}
	`
	td.Files["view/new.svelte"] = `
		<script>
			export let user = {}
		</script>
		<h1>new: {user.id} {user.name}</h1>
	`
	td.Files["view/show.svelte"] = `
		<script>
			export let user = {}
		</script>
		<h1>show: {user.id} {user.name}</h1>
	`
	td.Files["view/edit.svelte"] = `
		<script>
			export let user = {}
		</script>
		<h1>edit: {user.id} {user.name}</h1>
	`
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		type User struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Name string ` + "`" + `json:"name"` + "`" + `
		}
		func (c *Controller) Index() []*User {
			return []*User{{1, "a"}, {2, "b"}}
		}
		func (c *Controller) New() *User {
			return &User{3, "c"}
		}
		func (c *Controller) Show(id int) *User {
			return &User{id, "s"}
		}
		func (c *Controller) Edit(id int) *User {
			return &User{id, "e"}
		}
	`
	// Generate the app
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()

	// /
	res, err := app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":1,"name":"a"},{"id":2,"name":"b"}]
	`))
	res, err = app.Get("/")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err := res.Query("#bud_target")
	is.NoErr(err)
	html, err := el.Html()
	is.NoErr(err)
	is.Equal(`<h1>index: 1 a</h1><h1>index: 2 b</h1>`, html)

	// /new
	res, err = app.GetJSON("/new")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":3,"name":"c"}
	`))
	res, err = app.Get("/new")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>new: 3 c</h1>`, html)

	// /:id
	res, err = app.GetJSON("/10")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"s"}
	`))
	res, err = app.Get("/10")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>show: 10 s</h1>`, html)

	// /:id/edit
	res, err = app.GetJSON("/10/edit")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"e"}
	`))
	res, err = app.Get("/10/edit")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>edit: 10 e</h1>`, html)
	is.NoErr(app.Close())
}

func TestViewNestedUnnamed(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.Files["view/users/index.svelte"] = `
		<script>
			export let users = []
		</script>
		{#each users as user}
		<h1>index: {user.id} {user.name}</h1>
		{/each}
	`
	td.Files["view/users/new.svelte"] = `
		<script>
			export let user = {}
		</script>
		<h1>new: {user.id} {user.name}</h1>
	`
	td.Files["view/users/show.svelte"] = `
		<script>
			export let user = {}
		</script>
		<h1>show: {user.id} {user.name}</h1>
	`
	td.Files["view/users/edit.svelte"] = `
		<script>
			export let user = {}
		</script>
		<h1>edit: {user.id} {user.name}</h1>
	`
	td.Files["controller/users/users.go"] = `
		package users
		type Controller struct {}
		type User struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Name string ` + "`" + `json:"name"` + "`" + `
		}
		func (c *Controller) Index() []*User {
			return []*User{{1, "a"}, {2, "b"}}
		}
		func (c *Controller) New() *User {
			return &User{3, "c"}
		}
		func (c *Controller) Show(id int) *User {
			return &User{id, "s"}
		}
		func (c *Controller) Edit(id int) *User {
			return &User{id, "e"}
		}
	`
	// Generate the app
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()

	// /users
	res, err := app.GetJSON("/users")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":1,"name":"a"},{"id":2,"name":"b"}]
	`))
	res, err = app.Get("/users")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err := res.Query("#bud_target")
	is.NoErr(err)
	html, err := el.Html()
	is.NoErr(err)
	is.Equal(`<h1>index: 1 a</h1><h1>index: 2 b</h1>`, html)

	// /users/new
	res, err = app.GetJSON("/users/new")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":3,"name":"c"}
	`))
	res, err = app.Get("/users/new")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>new: 3 c</h1>`, html)

	// /users/:id
	res, err = app.GetJSON("/users/10")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"s"}
	`))
	res, err = app.Get("/users/10")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>show: 10 s</h1>`, html)

	// /users/:id/edit
	res, err = app.GetJSON("/users/10/edit")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"e"}
	`))
	res, err = app.Get("/users/10/edit")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>edit: 10 e</h1>`, html)
	is.NoErr(app.Close())
}

func TestViewDeepUnnamed(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.Files["view/teams/users/index.svelte"] = `
		<script>
			export let onlineUsers = []
		</script>
		{#each onlineUsers as user}
		<h1>index: {user.id} {user.name} {user.createdAt}</h1>
		{/each}
	`
	td.Files["view/teams/users/new.svelte"] = `
		<script>
			export let onlineUser = {}
		</script>
		<h1>new: {onlineUser.id} {onlineUser.name} {onlineUser.createdAt}</h1>
	`
	td.Files["view/teams/users/show.svelte"] = `
		<script>
			export let onlineUser = {}
		</script>
		<h1>show: {onlineUser.id} {onlineUser.name} {onlineUser.createdAt}</h1>
	`
	td.Files["view/teams/users/edit.svelte"] = `
		<script>
			export let onlineUser = {}
		</script>
		<h1>edit: {onlineUser.id} {onlineUser.name} {onlineUser.createdAt}</h1>
	`
	td.Files["controller/teams/users/users.go"] = `
		package users
		import "time"
		type Controller struct {}
		type OnlineUser struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Name string ` + "`" + `json:"name"` + "`" + `
			CreatedAt time.Time ` + "`" + `json:"createdAt"` + "`" + `
		}
		var now = time.Date(2021, 8, 4, 14, 56, 0, 0, time.UTC)
		func (c *Controller) Index() []*OnlineUser {
			return []*OnlineUser{{1, "a", now}, {2, "b", now}}
		}
		func (c *Controller) New() *OnlineUser {
			return &OnlineUser{3, "c", now}
		}
		func (c *Controller) Show(id int) *OnlineUser {
			return &OnlineUser{id, "s", now}
		}
		func (c *Controller) Edit(id int) *OnlineUser {
			return &OnlineUser{id, "e", now}
		}
	`
	// Generate the app
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()

	// /teams/:team_id/users
	res, err := app.GetJSON("/teams/5/users")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":1,"name":"a","createdAt":"2021-08-04T14:56:00Z"},{"id":2,"name":"b","createdAt":"2021-08-04T14:56:00Z"}]
	`))
	res, err = app.Get("/teams/5/users")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err := res.Query("#bud_target")
	is.NoErr(err)
	html, err := el.Html()
	is.NoErr(err)
	is.Equal(`<h1>index: 1 a 2021-08-04T14:56:00Z</h1><h1>index: 2 b 2021-08-04T14:56:00Z</h1>`, html)

	// /teams/:team_id/users/new
	res, err = app.GetJSON("/teams/5/users/new")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":3,"name":"c","createdAt":"2021-08-04T14:56:00Z"}
	`))
	res, err = app.Get("/teams/5/users/new")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>new: 3 c 2021-08-04T14:56:00Z</h1>`, html)

	// /teams/:team_id/users/:id
	res, err = app.GetJSON("/teams/5/users/10")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"s","createdAt":"2021-08-04T14:56:00Z"}
	`))
	res, err = app.Get("/teams/5/users/10")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>show: 10 s 2021-08-04T14:56:00Z</h1>`, html)

	// /teams/:team_id/users/:id/edit
	res, err = app.GetJSON("/teams/5/users/10/edit")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"e","createdAt":"2021-08-04T14:56:00Z"}
	`))
	res, err = app.Get("/teams/5/users/10/edit")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>edit: 10 e 2021-08-04T14:56:00Z</h1>`, html)
	is.NoErr(app.Close())
}

func TestResourceContext(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/users/users.go"] = `
		package users
		import contexts "context"
		type User struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Name string ` + "`" + `json:"name"` + "`" + `
		}
		type Controller struct {}
		func (c *Controller) Index(ctx contexts.Context) []*User {
			return []*User{{1, "a"}, {2, "b"}}
		}
		func (c *Controller) Create(ctx contexts.Context, name string) *User {
			return &User{1, name}
		}
		func (c *Controller) Show(ctx contexts.Context, id int) *User {
			return &User{id, "a"}
		}
		func (c *Controller) Edit(ctx contexts.Context, id int) *User {
			return &User{id, "a"}
		}
		func (c *Controller) Update(ctx contexts.Context, id int, name string) (*User, error) {
			return &User{id, name}, nil
		}
		func (c *Controller) Delete(ctx contexts.Context, id int) (*User, error) {
			return &User{id, "a"}, nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/users")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":1,"name":"a"},{"id":2,"name":"b"}]
	`))
	res, err = app.PostJSON("/users", bytes.NewBufferString(`{"name":"b"}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"name":"b"}
	`))
	res, err = app.GetJSON("/users/2")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"name":"a"}
	`))
	res, err = app.GetJSON("/users/2/edit")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"name":"a"}
	`))
	res, err = app.PatchJSON("/users/2", bytes.NewBufferString(`{"name":"b"}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"name":"b"}
	`))
	res, err = app.DeleteJSON("/users/2", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"name":"a"}
	`))
	is.NoErr(app.Close())
}

func TestOkChangeOk(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string {
			return "Hello Users!"
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: text/html
	`))
	is.In(res.Body().String(), `Hello Users!`)
	// Update controller
	controllerFile := filepath.Join(dir, "controller", "controller.go")
	is.NoErr(os.MkdirAll(filepath.Dir(controllerFile), 0755))
	is.NoErr(os.WriteFile(controllerFile, []byte(dedent.Dedent(`
		package controller
		type Controller struct {}
		func (c *Controller) Index() string {
			return "Hello Humans!"
		}
	`)), 0644))
	is.NoErr(app.Ready(ctx))
	// Try again with the new file
	res, err = app.Get("/")
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: text/html
	`))
	is.In(res.Body().String(), `Hello Humans!`)
	is.NoErr(app.Close())
}

func TestEmptyActionWithView(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() {}
	`
	td.Files["view/index.svelte"] = `<h1>hello</h1>`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	// HTML response
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Transfer-Encoding: chunked
		Content-Type: text/html
	`))
	is.In(res.Body().String(), `<h1>hello</h1>`)
	is.NoErr(app.Close())
}

func TestCustomActions(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) About() string { return "about" }
	`
	td.Files["controller/users/users.go"] = `
		package users
		type Controller struct {}
		func (c *Controller) Deactivate() string { return "deactivate" }
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/about")
	is.NoErr(err)
	// HTML response
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: text/html
	`))
	is.In(res.Body().String(), `about`)
	res, err = app.Get("/users/deactivate")
	is.NoErr(err)
	// HTML response
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 200 OK
		Content-Type: text/html
	`))
	is.In(res.Body().String(), `deactivate`)
	is.NoErr(app.Close())
}

func TestHandlerFuncs(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/foos/bars/controller.go"] = `
		package controller
		import "io"
		import "net/http"
		type Controller struct {}
		func (c *Controller) Index() string {
			return "hello"
		}
		func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(r.URL.Query().Get("foo_id")))
			io.Copy(w, r.Body)
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Test POST
	res, err := app.Post("/foos/some/bars", bytes.NewBufferString("body"))
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 201 Created
		Content-Type: text/plain; charset=utf-8
	`)
	is.Equal(res.Body().String(), `somebody`)
	// Test that regular actions continue to work
	res, err = app.GetJSON("/foos/some/bars")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		"hello"
	`)
	is.NoErr(app.Close())
}

func TestEnvSupport(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		import "os"
		type Controller struct {}
		func (c *Controller) Index() string {
			return os.Getenv("BUDDY")
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	cli.Env["BUDDY"] = "buddy"
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Test GET /
	res, err := app.Get("/")
	is.NoErr(err)
	res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"buddy"
	`)
	is.NoErr(app.Close())
}

func TestStructInStruct(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["model/post/post.go"] = `
		package post
		type Model struct {}
	`
	td.Files["controller/controller.go"] = `
		package controller
		import "context"
		import "app.com/model/post"
		type Controller struct {}
		type Post struct {
			Model post.Model
		}
		func (c *Controller) Show(ctx context.Context, id int) (post *Post, err error) {
			return &Post{}, nil
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	result, err := cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
}

// https://github.com/livebud/bud/issues/101
func TestLoadController(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		func Load() (*Controller, error) {
			return &Controller{}, nil
		}
		type Controller struct {}
		func (c *Controller) Index() string { return "" }
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Test GET /
	res, err := app.Get("/")
	is.NoErr(err)
	res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		""
	`)
	is.NoErr(app.Close())
}

// https://github.com/livebud/bud/issues/135
func TestSameNestedName(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/users/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "/users" }
	`
	td.Files["controller/admins/users/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "/admins/:id/users" }
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/users")
	is.NoErr(err)
	res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"/users"
	`)
	is.NoErr(app.Close())
	res, err = app.Get("/admins/10/users")
	is.NoErr(err)
	res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"/admins/:id/users"
	`)
	is.NoErr(app.Close())
}

func TestControllerChange(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() (string, error) { return "/", nil }
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/")
	is.NoErr(err)
	res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"/"
	`)
	// Update controller
	controllerFile := filepath.Join(dir, "controller", "controller.go")
	is.NoErr(os.MkdirAll(filepath.Dir(controllerFile), 0755))
	is.NoErr(os.WriteFile(controllerFile, []byte(dedent.Dedent(`
		package controller
		type Controller struct {}
		func (c *Controller) Index() string { return "/" }
	`)), 0644))
	// Wait for the app to be ready again
	is.NoErr(app.Ready(ctx))
	// Try again with the new file
	res, err = app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"/"
	`))
	is.NoErr(app.Close())
}

func TestRequestMap(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/posts/comments/controller.go"] = `
		package comments
		type Controller struct {}
		type Input struct {
			PostID int ` + "`" + `json:"post_id"` + "`" + `
			Order string ` + "`" + `json:"order"` + "`" + `
			Author *string ` + "`" + `json:"author"` + "`" + `
		}
		func (c *Controller) Index(in *Input) *Input {
			return in
		}
		func (c *Controller) Create(in *Input) *Input {
			return in
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/posts/10/comments?order=asc&author=Alice")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"post_id":10,"order":"asc","author":"Alice"}
	`))
	res, err = app.PostJSON("/posts/10/comments?order=asc", bytes.NewBufferString(`{"author":"Alice"}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"post_id":10,"order":"asc","author":"Alice"}
	`))
	// Test optional
	res, err = app.PostJSON("/posts/10/comments?order=asc", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"post_id":10,"order":"asc","author":null}
	`))
	is.NoErr(app.Close())
}

func TestComplexInput(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		type Email string
		type Op struct {
			Name string   ` + "`" + `json:"name"` + "`" + `
			Params []*Param ` + "`" + `json:"params"` + "`" + `
		}
		type Param struct {
			Version int
			Update bool ` + "`" + `json:"update"` + "`" + `
		}
		type Result struct {
			ID string
			Email string
			Op *Op
		}
		func (c *Controller) Update(id string, email string, op *Op) *Result {
			return &Result{id, email, op}
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PatchJSON("/123?email=alice@livebud.com", bytes.NewBufferString(`{
		"op": {
			"name": "update",
			"params": [
				{ "Version": 1, "update": true },
				{ "Version": 2 }
			]
		}
	}`))
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"ID":"123","Email":"alice@livebud.com","Op":{"name":"update","params":[{"Version":1,"update":true},{"Version":2,"update":false}]}}
	`))
}

func TestRedirectBack(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		import "errors"
		type Controller struct {}
		func (c *Controller) Index() string {
			return "index"
		}
		func (c *Controller) New() string {
			return "new"
		}
		func (c *Controller) Edit() string {
			return "new"
		}
		func (c *Controller) Create() error {
			return errors.New("create error")
		}
		func (c *Controller) Update() error {
			return errors.New("update error")
		}
		func (c *Controller) Delete() error {
			return errors.New("update error")
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Post request
	req, err := app.PostRequest("/", nil)
	is.NoErr(err)
	req.Header.Set("Referer", "/new")
	res, err := app.Do(req)
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 303 See Other
		Location: /new
	`))
	// Post request, no referer
	req, err = app.PostRequest("/", nil)
	is.NoErr(err)
	res, err = app.Do(req)
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 303 See Other
		Location: /
	`))
	// Patch request
	req, err = app.PatchRequest("/10", nil)
	is.NoErr(err)
	req.Header.Set("Referer", "/10/edit")
	res, err = app.Do(req)
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 303 See Other
		Location: /10/edit
	`))
	// Patch request, no referer
	req, err = app.PatchRequest("/10", nil)
	is.NoErr(err)
	res, err = app.Do(req)
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 303 See Other
		Location: /10
	`))
	// Delete request
	req, err = app.DeleteRequest("/10", nil)
	is.NoErr(err)
	req.Header.Set("Referer", "/10")
	res, err = app.Do(req)
	is.NoErr(err)
	is.NoErr(res.DiffHeaders(`
		HTTP/1.1 303 See Other
		Location: /10
	`))
}

func TestInject(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["log/log.go"] = `
		package log
		func New() *Logger { return &Logger{} }
		type Logger struct {}
		func (l *Logger) Info(msg string) {}
	`
	td.Files["session/session.go"] = `
		package session
		import "errors"
		import "net/http"
		import "app.com/log"
		func New(log *log.Logger, w http.ResponseWriter, r *http.Request) *Session {
			return &Session{log, w, r}
		}
		type Session struct {
			log *log.Logger
			w http.ResponseWriter
			r *http.Request
		}
		func (s *Session) Set(key, value string) {
			s.log.Info("setting session")
			http.SetCookie(s.w, &http.Cookie{Name: key, Value: value })
		}
		func (s *Session) Clear() error {
			return errors.New("session: unable to clear")
		}
	`
	td.Files["db/db.go"] = `
		package db
		import "context"
		var loaded = 0
		func Load(ctx context.Context) (*Client, error) {
			loaded++
			return &Client{}, nil
		}
		type Client struct{}
		func (c *Client) Loaded() int {
			return loaded
		}
	`
	td.Files["controller/controller.go"] = `
		package controller
		import "app.com/session"
		import "app.com/db"
		type Controller struct {
			Session *session.Session
			DB *db.Client
		}
		func (c *Controller) Create() error {
			c.Session.Set("sessionid", "some-key")
			return nil
		}
		func (c *Controller) Delete() error {
			return c.Session.Clear()
		}
		func (c *Controller) Index() int {
			return c.DB.Loaded()
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Test errors from dependencies
	res, err := app.DeleteJSON("/10", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
			HTTP/1.1 500 Internal Server Error
			Content-Type: application/json

			{"error":"session: unable to clear"}
		`))
	// Post request continue to work
	res, err = app.PostJSON("/", nil)
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 204 No Content
		Set-Cookie: sessionid=some-key
	`))
	// Post request continue to work
	res, err = app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		1
	`))
}

func TestEscapeProps(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		type Post struct {
			HTML string ` + "`json:\"html\"`" + `
		}
		func (c *Controller) Show(id string) *Post {
			return &Post{HTML: ` + "`" + `<b>hello ` + "`" + ` + id + ` + "`" + `<script type="text/javascript">alert('xss!')</script></b>` + "`" + `}
		}
	`
	td.Files["view/show.svelte"] = `
		<script>
			export let post = {}
		</script>
		<h1>{post.html}</h1>
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Post request
	res, err := app.Get("/alice")
	is.NoErr(err)
	props, err := res.Query("#bud_props")
	is.NoErr(err)
	is.Equal(props.Text(), `{"post":{"html":"<b>hello alice<script type=\"text/javascript\">alert('xss!')<\/script></b>"}}`)
	target, err := res.Query("#bud_target")
	is.NoErr(err)
	is.Equal(target.Text(), `<b>hello alice<script type="text/javascript">alert('xss!')</script></b>`)
}

func TestProtocol(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = versions.Svelte
	td.NodeModules["livebud"] = "*"
	td.Files["model/model.go"] = `
		package model
		type MyArticle struct {
			SlugID string
		}
	`
	td.Files["controller/controller.go"] = `
		package controller
		import "app.com/model"
		type Controller struct {}
		type MyPost struct {
			ID string ` + "`" + `json:"id"` + "`" + `
		}
		func (c *Controller) Index() string { return "index" }
		func (c *Controller) IndexErr() (string, error) { return "index_err", nil }
		func (c *Controller) Named() (name string) { return "named" }
		func (c *Controller) NamedErr() (name string, err error) { return "named_err", nil }
		func (c *Controller) Post() *MyPost { return &MyPost{ID:"post"} }
		func (c *Controller) PostErr() (*MyPost, error) { return &MyPost{ID:"post_err"}, nil }
		func (c *Controller) NamedPost() (post *MyPost) { return &MyPost{ID:"named_post"} }
		func (c *Controller) NamedPostErr() (post *MyPost, err error) { return &MyPost{ID:"named_post_err"}, nil }
		func (c *Controller) Article() *model.MyArticle { return &model.MyArticle{SlugID:"article"} }
		func (c *Controller) ArticleErr() (*model.MyArticle, error) { return &model.MyArticle{SlugID:"article_err"}, nil }
		func (c *Controller) NamedArticle() (article *model.MyArticle) { return &model.MyArticle{SlugID:"named_article"} }
		func (c *Controller) NamedArticleErr() (article *model.MyArticle, err error) { return &model.MyArticle{SlugID:"named_article_err"}, nil }
	`
	td.Files["view/index.svelte"] = `
		<script>
			export let _string = ""
		</script>
		<h1>{_string}</h1>
	`
	td.Files["view/index_err.svelte"] = `
		<script>
			export let _string = ""
		</script>
		<h1>{_string}</h1>
	`
	td.Files["view/named.svelte"] = `
		<script>
			export let name = ""
		</script>
		<h1>{name}</h1>
	`
	td.Files["view/named_err.svelte"] = `
		<script>
			export let name = ""
		</script>
		<h1>{name}</h1>
	`
	td.Files["view/post.svelte"] = `
		<script>
			export let myPost = {}
		</script>
		<h1>{myPost.id}</h1>
	`
	td.Files["view/post_err.svelte"] = `
		<script>
		export let myPost = {}
		</script>
		<h1>{myPost.id}</h1>
	`
	td.Files["view/named_post.svelte"] = `
		<script>
			export let post = {}
		</script>
		<h1>{post.id}</h1>
	`
	td.Files["view/named_post_err.svelte"] = `
		<script>
		export let post = {}
		</script>
		<h1>{post.id}</h1>
	`
	td.Files["view/article.svelte"] = `
		<script>
			export let myArticle = {}
		</script>
		<h1>{myArticle.SlugID}</h1>
	`
	td.Files["view/article_err.svelte"] = `
		<script>
		export let myArticle = {}
		</script>
		<h1>{myArticle.SlugID}</h1>
	`
	td.Files["view/named_article.svelte"] = `
		<script>
			export let article = {}
		</script>
		<h1>{article.SlugID}</h1>
	`
	td.Files["view/named_article_err.svelte"] = `
		<script>
		export let article = {}
		</script>
		<h1>{article.SlugID}</h1>
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	app, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// index
	res, err := app.GetJSON("/")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"index"
	`))
	res, err = app.Get("/")
	is.NoErr(err)
	sel, err := res.Query("#bud_target")
	is.NoErr(err)
	html, err := sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>index</h1>`)
	// index error
	res, err = app.GetJSON("/index_err")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"index_err"
	`))
	res, err = app.Get("/index_err")
	is.NoErr(err)
	is.Equal(res.Status(), 200)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>index_err</h1>`)
	// named
	res, err = app.GetJSON("/named")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"named"
	`))
	res, err = app.Get("/named")
	is.NoErr(err)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>named</h1>`)
	// named err
	res, err = app.GetJSON("/named_err")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		"named_err"
	`))
	res, err = app.Get("/named_err")
	is.NoErr(err)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>named_err</h1>`)
	// post
	res, err = app.GetJSON("/post")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":"post"}
	`))
	res, err = app.Get("/post")
	is.NoErr(err)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>post</h1>`)
	// post err
	res, err = app.GetJSON("/post_err")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":"post_err"}
	`))
	res, err = app.Get("/post_err")
	is.NoErr(err)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>post_err</h1>`)
	// named post
	res, err = app.GetJSON("/named_post")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":"named_post"}
	`))
	res, err = app.Get("/named_post")
	is.NoErr(err)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>named_post</h1>`)
	// named post err
	res, err = app.GetJSON("/named_post_err")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":"named_post_err"}
	`))
	res, err = app.Get("/named_post_err")
	is.NoErr(err)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>named_post_err</h1>`)
	// article
	res, err = app.GetJSON("/article")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"SlugID":"article"}
	`))
	res, err = app.Get("/article")
	is.NoErr(err)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>article</h1>`)
	// article err
	res, err = app.GetJSON("/article_err")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"SlugID":"article_err"}
	`))
	res, err = app.Get("/article_err")
	is.NoErr(err)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>article_err</h1>`)
	// named article
	res, err = app.GetJSON("/named_article")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"SlugID":"named_article"}
	`))
	res, err = app.Get("/named_article")
	is.NoErr(err)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>named_article</h1>`)
	// named article err
	res, err = app.GetJSON("/named_article_err")
	is.NoErr(err)
	is.NoErr(res.Diff(`
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"SlugID":"named_article_err"}
	`))
	res, err = app.Get("/named_article_err")
	is.NoErr(err)
	sel, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = sel.Html()
	is.NoErr(err)
	is.Equal(html, `<h1>named_article_err</h1>`)
	// close the app
	is.NoErr(app.Close())
}
