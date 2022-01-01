package action_test

import (
	"bytes"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/test"
)

func TestIndexString(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action

		type Controller struct {
		}

		func (c *Controller) Index() string {
			return "Hello Users!"
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		"Hello Users!"
	`)
}

func TestAboutIndexString(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/about/action.go"] = `
		package action
		type Controller struct {}
		func (c *Controller) Index() string { return "About" }
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/about")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		"About"
	`)
}

func TestCreate302(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		type Controller struct {}
		func (c *Controller) Create() {
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Post("/", nil)
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 302 Found
		Date: Fri, 31 Dec 2021 00:00:00 GMT
		Location: /
	`)
}

func TestIndex204(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		type Controller struct {}
		func (c *Controller) Index() {
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 204 No Content
		Date: Fri, 31 Dec 2021 00:00:00 GMT
	`)
}

func TestIndex500(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		import "errors"
		type Controller struct {}
		type Post struct {}
		func (c *Controller) Index() ([]*Post, error) {
			return nil, errors.New("unable to list posts")
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		{"error":"unable to list posts"}
	`)
}

func TestIndexList500(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		import "errors"
		type Controller struct {}
		type Post struct {}
		func (c *Controller) Index() (int, string, error) {
			return 0, "", errors.New("unable to list posts")
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 302 Found
		Date: Fri, 31 Dec 2021 00:00:00 GMT
		Location: /
	`)
}

func TestIndexList200(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		type Controller struct {}
		type Post struct {}
		func (c *Controller) Index() (int, string, error) {
			return 0, "a", nil
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 302 Found
		Date: Fri, 31 Dec 2021 00:00:00 GMT
		Location: /
	`)
}

func TestIndexListObject500(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		import "errors"
		type Controller struct {}
		type Post struct {}
		func (c *Controller) Index() (a int, b string, err error) {
			return 0, "a", errors.New("unable to list posts")
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 302 Found
		Date: Fri, 31 Dec 2021 00:00:00 GMT
		Location: /
	`)
}

func TestIndexListObject200(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		type Controller struct {}
		type Post struct {}
		func (c *Controller) Index() (a int, b string, err error) {
			return 0, "a", nil
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 302 Found
		Date: Fri, 31 Dec 2021 00:00:00 GMT
		Location: /
	`)
}

func TestIndexStructs200(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		type Controller struct {}
		type Post struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
		func (c *Controller) Index() (list []*Post, err error) {
			return []*Post{{0, "a"}, {1, "b"}}, nil
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		[{"id":0,"title":"a"},{"id":1,"title":"b"}]
	`)
}

func TestJSONCreate204(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		type Controller struct {}
		func (c *Controller) Create() {}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.PostJSON("/", nil)
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 204 No Content
		Content-Length: 0
		Date: Fri, 31 Dec 2021 00:00:00 GMT
	`)
}

func TestJSONCreate500(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		import "errors"
		type Controller struct {}
		func (c *Controller) Create() (string, error) {
			return "", errors.New("Not implemented yet")
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.PostJSON("/", nil)
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		{"error":"Not implemented yet"}
	`)
}

func TestDependencyHoist(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["postgres/pool.go"] = `
		package postgres
		func New() *Pool { return &Pool{1} }
		type Pool struct { ID int }
	`
	generator.Files["action/action.go"] = `
		package action
		import "app.com/postgres"
		type Controller struct {
			Pool *postgres.Pool
		}
		func (c *Controller) Index() int {
			return c.Pool.ID
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		1
	`)
}

func TestDependencyRequest(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["postgres/pool.go"] = `
		package postgres
		import "net/http"
		func New(r *http.Request) *Pool { return &Pool{r.URL.Path} }
		type Pool struct { Path string }
	`
	generator.Files["action/action.go"] = `
		package action
		import "app.com/postgres"
		type Controller struct {
			Pool *postgres.Pool
		}
		func (c *Controller) Index() string {
			return c.Pool.Path
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.Get("/")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		"/"
	`)
}

func TestShareStruct(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["article/article.go"] = `
		package article
		type Article struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
	`
	generator.Files["action/action.go"] = `
		package action
		import "app.com/article"
		type Controller struct {
		}
		func (c *Controller) Update(a *article.Article) (*article.Article, error) {
			return a, nil
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.PatchJSON("/10", bytes.NewBufferString(`{"id": 1, "title": "a"}`))
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		{"id":10,"title":"a"}
	`)
}

func TestJSONCreateNested(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["postgres/pool.go"] = `
		package postgres
		func New(r *http.Request) *Pool { return &Pool{r.URL.Path} }
		type Pool struct { Path string }
	`
	generator.Files["action/action.go"] = `
		package action
		type Controller struct {}
		type Post struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
		func (c *Controller) Create(p *Post) *Post {
			return p
		}
	`
	generator.Files["action/users/users.go"] = `
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
	generator.Files["action/users/admin/admin.go"] = `
		package admin
		type Controller struct {}
		type Post struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
		func (c *Controller) Create(p *Post) *Post {
			return p
		}
	`
	generator.Files["action/articles/articles.go"] = `
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
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.PostJSON("/", bytes.NewBufferString(`{"id": "1", title: "a"}`))
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 302 Found
		Date: Fri, 31 Dec 2021 00:00:00 GMT
		Location: /
	`)
	res, err = server.PostJSON("/users", bytes.NewBufferString(`{"id": "2", title: "b"}`))
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 302 Found
		Date: Fri, 31 Dec 2021 00:00:00 GMT
		Location: /
	`)
	res, err = server.PostJSON("/users/admin", bytes.NewBufferString(`{"id": "3", title: "c"}`))
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 302 Found
		Date: Fri, 31 Dec 2021 00:00:00 GMT
		Location: /
	`)
	res, err = server.PostJSON("/articles", bytes.NewBufferString(`{"id": "4", title: "d"}`))
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 302 Found
		Date: Fri, 31 Dec 2021 00:00:00 GMT
		Location: /
	`)
}

func TestJSONDelete500(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		import "errors"
		type Controller struct {}
		func (c *Controller) Delete() (string, error) {
			return "", errors.New("Not implemented yet")
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.DeleteJSON("/1", nil)
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		{"error":"Not implemented yet"}
	`)
}

func TestJSONDelete200(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		type Controller struct {}
		type Post struct {
			ID int
			Title string
		}
		func (c *Controller) Delete(id int) (*Post, error) {
			return &Post{id, "a"}, nil
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.DeleteJSON("/1", nil)
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		{"ID":1,"Title":"a"}
	`)
}

func TestJSONMultipleActions(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
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
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.GetJSON("/")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		"hello world"
	`)
	res, err = server.GetJSON("/10")
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		10
	`)
}

func TestJSONUpdate500(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		import "errors"
		type Controller struct {}
		func (c *Controller) Update() (string, error) {
			return "", errors.New("Not implemented yet")
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.PatchJSON("/1", nil)
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 500 Internal Server Error
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		{"error":"Not implemented yet"}
	`)
}

func TestJSONUpdate200(t *testing.T) {
	is := is.New(t)
	generator := test.Generator(t)
	generator.Files["action/action.go"] = `
		package action
		type Controller struct {}
		type Post struct {
			ID int ` + "`" + `json:"id"` + "`" + `
			Title string ` + "`" + `json:"title"` + "`" + `
		}
		func (c *Controller) Update(id int) (*Post, error) {
			return &Post{id, "a"}, nil
		}
	`
	app, err := generator.Generate()
	is.NoErr(err)
	is.True(app.Exists("bud/action/action.go"))
	is.True(app.Exists("bud/main.go"))
	server, err := app.Start()
	is.NoErr(err)
	defer server.Close()
	res, err := server.PatchJSON("/1", nil)
	is.NoErr(err)
	res.Expect(`
		HTTP/1.1 200 OK
		Content-Type: application/json
		Date: Fri, 31 Dec 2021 00:00:00 GMT

		{"id":1,"title":"a"}
	`)
}
