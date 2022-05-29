package controller_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/matthewmueller/diff"

	"github.com/livebud/bud/internal/cli"
	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/version"
)

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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// HTML response
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), "Root")
	// JSON response
	res, err = app.GetJSON("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		"Root"
	`)
	// HTML response
	res, err = app.Get("/about")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), "About")
	// JSON response
	res, err = app.GetJSON("/about")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		"About"
	`)
	// HTML response
	res, err = app.Get("/posts/some-slug/comments")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), "Comments")
	// JSON response
	res, err = app.GetJSON("/posts/some-slug/comments")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		"Comments"
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Redirect /
	res, err := app.Post("/", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 302 Found
		Location: /
	`)
	is.Equal(res.Body().Len(), 0)
	// Redirect /users
	res, err = app.Post("/users", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /users
	`)
	is.Equal(res.Body().Len(), 0)
	// Redirect /posts/10/comments
	res, err = app.Post("/posts/10/comments", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /posts/10/comments
	`)
	is.Equal(res.Body().Len(), 0)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Root
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.Get("/10")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.Get("/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.Get("/10/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.GetJSON("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.GetJSON("/10")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.GetJSON("/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.GetJSON("/10/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	// Comments
	res, err = app.Get("/posts/1/comments")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.Get("/posts/1/comments/5")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.Get("/posts/1/comments/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.Get("/posts/1/comments/10/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.GetJSON("/posts/1/comments")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.GetJSON("/posts/1/comments/10")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.GetJSON("/posts/1/comments/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.GetJSON("/posts/1/comments/10/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json

		{"error":"unable to list posts"}
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":0,"title":"a"},{"id":1,"title":"b"}]
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PostJSON("/", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PostJSON("/", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json

		{"error":"Not implemented yet"}
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		1
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), `/`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PatchJSON("/10", bytes.NewBufferString(`{"id": 1, "title": "a"}`))
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"title":"a"}
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PostJSON("/", bytes.NewBufferString(`{"id": 1, "title": "a"}`))
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"title":"a"}
	`)
	res, err = app.PostJSON("/users", bytes.NewBufferString(`{"id": 2, "title": "b"}`))
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"title":"b"}
	`)
	res, err = app.PostJSON("/users/1/admin", bytes.NewBufferString(`{"id": 3, "title": "c"}`))
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"user_id":1,"id":3,"title":"c"}
	`)
	res, err = app.PostJSON("/articles", bytes.NewBufferString(`{"id": 4, "title": "d"}`))
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":4,"title":"d"}
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}

func TestJSONDelete500(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		import "errors"
		type Controller struct {}
		func (c *Controller) Delete() (string, error) {
			return "", errors.New("Not implemented yet")
		}
	`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.DeleteJSON("/1", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json

		{"error":"Not implemented yet"}
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.DeleteJSON("/1", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"ID":1,"Title":"a"}
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		"hello world"
	`)
	res, err = app.GetJSON("/10")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		10
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PatchJSON("/1", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json

		{"error":"Not implemented yet"}
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.PatchJSON("/1", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"title":"a"}
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/users")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":1,"name":"a","age":2},{"id":2,"name":"b","age":3}]
	`)
	res, err = app.GetJSON("/users/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.PostJSON("/users?name=matt&age=10", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":3,"name":"matt","age":10}
	`)
	res, err = app.GetJSON("/users/10")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"d","age":5}
	`)
	res, err = app.GetJSON("/users/10/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.PatchJSON("/users/10", bytes.NewBufferString(`{"name": "matt", "age": 10}`))
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.DeleteJSON("/users/10", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/posts/1/comments")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":2,"post_id":1,"title":"a"},{"id":3,"post_id":1,"title":"b"}]
	`)
	res, err = app.GetJSON("/posts/1/comments/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.PostJSON("/posts/1/comments", bytes.NewBufferString(`{"title":"1st"}`))
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"post_id":1,"title":"1st"}
	`)
	res, err = app.GetJSON("/posts/1/comments/2")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"post_id":1,"title":"a"}
	`)
	res, err = app.GetJSON("/posts/1/comments/2/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 204 No Content
	`)
	res, err = app.PatchJSON("/posts/1/comments/2", bytes.NewBufferString(`{"title":"1st"}`))
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"post_id":2,"title":"1st"}
	`)
	res, err = app.PatchJSON("/posts/1/comments/2", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"post_id":2}
	`)
	res, err = app.DeleteJSON("/posts/1/comments/2", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"post_id":2}
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	// Root
	res, err := app.Post("/", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /2
	`)
	res, err = app.Patch("/1", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /1
	`)
	res, err = app.Delete("/1", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /
	`)
	// Posts
	res, err = app.Post("/posts", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /posts/2
	`)
	res, err = app.Patch("/posts/1", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /posts/1
	`)
	res, err = app.Delete("/posts/1", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /posts
	`)
	// Comments
	res, err = app.Post("/posts/1/comments", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /posts/1/comments/2
	`)
	res, err = app.Patch("/posts/1/comments/2", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /posts/1/comments/2
	`)
	res, err = app.Delete("/posts/1/comments/2", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 302 Found
		Location: /posts/1/comments
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}

func TestViewUnnamed(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = version.Svelte
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()

	// /
	res, err := app.GetJSON("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":1,"name":"a"},{"id":2,"name":"b"}]
	`)
	res, err = app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err := res.Query("#bud_target")
	is.NoErr(err)
	html, err := el.Html()
	is.NoErr(err)
	is.Equal(`<h1>index: 1 a</h1><h1>index: 2 b</h1>`, html)

	// /new
	res, err = app.GetJSON("/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":3,"name":"c"}
	`)
	res, err = app.Get("/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>new: 3 c</h1>`, html)

	// /:id
	res, err = app.GetJSON("/10")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"s"}
	`)
	res, err = app.Get("/10")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>show: 10 s</h1>`, html)

	// /:id/edit
	res, err = app.GetJSON("/10/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"e"}
	`)
	res, err = app.Get("/10/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>edit: 10 e</h1>`, html)

	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}

func TestViewNestedUnnamed(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = version.Svelte
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()

	// /users
	res, err := app.GetJSON("/users")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":1,"name":"a"},{"id":2,"name":"b"}]
	`)
	res, err = app.Get("/users")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err := res.Query("#bud_target")
	is.NoErr(err)
	html, err := el.Html()
	is.NoErr(err)
	is.Equal(`<h1>index: 1 a</h1><h1>index: 2 b</h1>`, html)

	// /users/new
	res, err = app.GetJSON("/users/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":3,"name":"c"}
	`)
	res, err = app.Get("/users/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>new: 3 c</h1>`, html)

	// /users/:id
	res, err = app.GetJSON("/users/10")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"s"}
	`)
	res, err = app.Get("/users/10")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>show: 10 s</h1>`, html)

	// /users/:id/edit
	res, err = app.GetJSON("/users/10/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"e"}
	`)
	res, err = app.Get("/users/10/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>edit: 10 e</h1>`, html)

	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}

func TestViewDeepUnnamed(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = version.Svelte
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()

	// /teams/:team_id/users
	res, err := app.GetJSON("/teams/5/users")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":1,"name":"a","createdAt":"2021-08-04T14:56:00Z"},{"id":2,"name":"b","createdAt":"2021-08-04T14:56:00Z"}]
	`)
	res, err = app.Get("/teams/5/users")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err := res.Query("#bud_target")
	is.NoErr(err)
	html, err := el.Html()
	is.NoErr(err)
	is.Equal(`<h1>index: 1 a 2021-08-04T14:56:00Z</h1><h1>index: 2 b 2021-08-04T14:56:00Z</h1>`, html)

	// /teams/:team_id/users/new
	res, err = app.GetJSON("/teams/5/users/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":3,"name":"c","createdAt":"2021-08-04T14:56:00Z"}
	`)
	res, err = app.Get("/teams/5/users/new")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>new: 3 c 2021-08-04T14:56:00Z</h1>`, html)

	// /teams/:team_id/users/:id
	res, err = app.GetJSON("/teams/5/users/10")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"s","createdAt":"2021-08-04T14:56:00Z"}
	`)
	res, err = app.Get("/teams/5/users/10")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>show: 10 s 2021-08-04T14:56:00Z</h1>`, html)

	// /teams/:team_id/users/:id/edit
	res, err = app.GetJSON("/teams/5/users/10/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":10,"name":"e","createdAt":"2021-08-04T14:56:00Z"}
	`)
	res, err = app.Get("/teams/5/users/10/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	el, err = res.Query("#bud_target")
	is.NoErr(err)
	html, err = el.Html()
	is.NoErr(err)
	is.Equal(`<h1>edit: 10 e 2021-08-04T14:56:00Z</h1>`, html)

	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.GetJSON("/users")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		[{"id":1,"name":"a"},{"id":2,"name":"b"}]
	`)
	res, err = app.PostJSON("/users", bytes.NewBufferString(`{"name":"b"}`))
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":1,"name":"b"}
	`)
	res, err = app.GetJSON("/users/2")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"name":"a"}
	`)
	res, err = app.GetJSON("/users/2/edit")
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"name":"a"}
	`)
	res, err = app.PatchJSON("/users/2", bytes.NewBufferString(`{"name":"b"}`))
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"name":"b"}
	`)
	res, err = app.DeleteJSON("/users/2", nil)
	is.NoErr(err)
	diff.TestHTTP(t, res.Dump().String(), `
		HTTP/1.1 200 OK
		Content-Type: application/json

		{"id":2,"name":"a"}
	`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	hot, err := app.Hot("/bud/view/index.svelte")
	is.NoErr(err)
	defer hot.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), `Hello Users!`)
	// Update file
	td = testdir.New(dir)
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() string {
			return "Hello Humans!"
		}
	`
	is.NoErr(td.Write(ctx))
	// Wait for the change event
	eventCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	event, err := hot.Next(eventCtx)
	is.NoErr(err)
	is.Equal(string(event.Data), `{"reload":true}`)
	// Try again with the new file
	res, err = app.Get("/")
	is.NoErr(err)
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), `Hello Humans!`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.In(stderr.String(), "info: Ready on")
}

func TestEmptyActionWithView(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.NodeModules["svelte"] = version.Svelte
	td.Files["controller/controller.go"] = `
		package controller
		type Controller struct {}
		func (c *Controller) Index() {}
	`
	td.Files["view/index.svelte"] = `<h1>hello</h1>`
	is.NoErr(td.Write(ctx))
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/")
	is.NoErr(err)
	// HTML response
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), `<h1>hello</h1>`)
	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
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
	cli := testcli.New(cli.New(dir))
	app, stdout, stderr, err := cli.Start(ctx, "run")
	is.NoErr(err)
	defer app.Close()
	res, err := app.Get("/about")
	is.NoErr(err)
	// HTML response
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), `about`)
	res, err = app.Get("/users/deactivate")
	is.NoErr(err)
	// HTML response
	diff.TestHTTP(t, res.Headers().String(), `
		HTTP/1.1 200 OK
		Content-Type: text/html
	`)
	is.In(res.Body().String(), `deactivate`)

	// Test stdio
	is.Equal(stdout.String(), "")
	is.Equal(stderr.String(), "")
}
