package action_test

import (
	"context"
	"testing"

	"gitlab.com/mnm/bud/internal/budtest"

	"github.com/matryer/is"
)

func TestIndexString(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := "_tmp"
	bud := budtest.New(dir)
	bud.Files["action/action.go"] = `
		package action

		type Controller struct {
		}

		func (c *Controller) Index() string {
			return "Hello Users!"
		}
	`
	project, err := bud.Compile(ctx)
	is.NoErr(err)
	app, err := project.Build(ctx)
	is.NoErr(err)
	is.NoErr(app.Exists("bud/.app/action/action.go"))
	is.NoErr(app.Exists("bud/.app/main.go"))
	server, err := app.Start(ctx)
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

// func TestAboutIndexString(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/about/action.go"] = []byte(`
// 		package action
// 		type Controller struct {}
// 		func (c *Controller) Index() string { return "About" }
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Get("/about")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		"About"
// 	`)
// }

// func TestCreate302(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		type Controller struct {}
// 		func (c *Controller) Create() {
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Post("/", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /
// 	`)
// }

// func TestIndex204(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		type Controller struct {}
// 		func (c *Controller) Index() {
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Get("/")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 204 No Content
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// }

// func TestIndex500(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		import "errors"
// 		type Controller struct {}
// 		type Post struct {}
// 		func (c *Controller) Index() ([]*Post, error) {
// 			return nil, errors.New("unable to list posts")
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Get("/")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 500 Internal Server Error
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"error":"unable to list posts"}
// 	`)
// }

// func TestIndexList500(t *testing.T) {
// 	t.SkipNow()
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		import "errors"
// 		type Controller struct {}
// 		type Post struct {}
// 		func (c *Controller) Index() (int, string, error) {
// 			return 0, "", errors.New("unable to list posts")
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Get("/")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /
// 	`)
// }

// func TestIndexList200(t *testing.T) {
// 	t.SkipNow()
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		type Controller struct {}
// 		type Post struct {}
// 		func (c *Controller) Index() (int, string, error) {
// 			return 0, "a", nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Get("/")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /
// 	`)
// }

// func TestIndexListObject500(t *testing.T) {
// 	t.SkipNow()
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		import "errors"
// 		type Controller struct {}
// 		type Post struct {}
// 		func (c *Controller) Index() (a int, b string, err error) {
// 			return 0, "a", errors.New("unable to list posts")
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Get("/")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /
// 	`)
// }

// func TestIndexListObject200(t *testing.T) {
// 	t.SkipNow()
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		type Controller struct {}
// 		type Post struct {}
// 		func (c *Controller) Index() (a int, b string, err error) {
// 			return 0, "a", nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Get("/")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /
// 	`)
// }

// func TestIndexStructs200(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		type Controller struct {}
// 		type Post struct {
// 			ID int ` + "`" + `json:"id"` + "`" + `
// 			Title string ` + "`" + `json:"title"` + "`" + `
// 		}
// 		func (c *Controller) Index() (list []*Post, err error) {
// 			return []*Post{{0, "a"}, {1, "b"}}, nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Get("/")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		[{"id":0,"title":"a"},{"id":1,"title":"b"}]
// 	`)
// }

// func TestJSONCreate204(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		type Controller struct {}
// 		func (c *Controller) Create() {}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.PostJSON("/", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 204 No Content
// 		Content-Length: 0
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// }

// func TestJSONCreate500(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		import "errors"
// 		type Controller struct {}
// 		func (c *Controller) Create() (string, error) {
// 			return "", errors.New("Not implemented yet")
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.PostJSON("/", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 500 Internal Server Error
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"error":"Not implemented yet"}
// 	`)
// }

// func TestDependencyHoist(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["postgres/pool.go"] = []byte(`
// 		package postgres
// 		func New() *Pool { return &Pool{1} }
// 		type Pool struct { ID int }
// 	`)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		import "app.com/postgres"
// 		type Controller struct {
// 			Pool *postgres.Pool
// 		}
// 		func (c *Controller) Index() int {
// 			return c.Pool.ID
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Get("/")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		1
// 	`)
// }

// func TestDependencyRequest(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["postgres/pool.go"] = []byte(`
// 		package postgres
// 		import "net/http"
// 		func New(r *http.Request) *Pool { return &Pool{r.URL.Path} }
// 		type Pool struct { Path string }
// 	`)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		import "app.com/postgres"
// 		type Controller struct {
// 			Pool *postgres.Pool
// 		}
// 		func (c *Controller) Index() string {
// 			return c.Pool.Path
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Get("/")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		"/"
// 	`)
// }

// func TestShareStruct(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["article/article.go"] = []byte(`
// 		package article
// 		type Article struct {
// 			ID int ` + "`" + `json:"id"` + "`" + `
// 			Title string ` + "`" + `json:"title"` + "`" + `
// 		}
// 	`)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		import "app.com/article"
// 		type Controller struct {
// 		}
// 		func (c *Controller) Update(a *article.Article) (*article.Article, error) {
// 			return a, nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.PatchJSON("/10", bytes.NewBufferString(`{"id": 1, "title": "a"}`))
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":10,"title":"a"}
// 	`)
// }

// func TestJSONCreateNested(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["postgres/pool.go"] = []byte(`
// 		package postgres
// 		func New(r *http.Request) *Pool { return &Pool{r.URL.Path} }
// 		type Pool struct { Path string }
// 	`)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		type Controller struct {}
// 		type Post struct {
// 			ID int ` + "`" + `json:"id"` + "`" + `
// 			Title string ` + "`" + `json:"title"` + "`" + `
// 		}
// 		func (c *Controller) Create(p *Post) *Post {
// 			return p
// 		}
// 	`)
// 	generator.Files["action/users/users.go"] = []byte(`
// 		package users
// 		type Controller struct {}
// 		type Post struct {
// 			ID int ` + "`" + `json:"id"` + "`" + `
// 			Title string ` + "`" + `json:"title"` + "`" + `
// 		}
// 		func (c *Controller) Create(p *Post) *Post {
// 			return p
// 		}
// 	`)
// 	generator.Files["action/users/admin/admin.go"] = []byte(`
// 		package admin
// 		type Controller struct {}
// 		type Post struct {
// 			UserID int ` + "`" + `json:"user_id"` + "`" + `
// 			ID int ` + "`" + `json:"id"` + "`" + `
// 			Title string ` + "`" + `json:"title"` + "`" + `
// 		}
// 		func (c *Controller) Create(p *Post) *Post {
// 			return p
// 		}
// 	`)
// 	generator.Files["action/articles/articles.go"] = []byte(`
// 		package articles
// 		type Controller struct {}
// 		type Post struct {
// 			ID int ` + "`" + `json:"id"` + "`" + `
// 			Title string ` + "`" + `json:"title"` + "`" + `
// 		}
// 		func (c *Controller) Create(p *Post) *Post {
// 			return p
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.PostJSON("/", bytes.NewBufferString(`{"id": 1, "title": "a"}`))
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":1,"title":"a"}
// 	`)
// 	res, err = server.PostJSON("/users", bytes.NewBufferString(`{"id": 2, "title": "b"}`))
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":2,"title":"b"}
// 	`)
// 	res, err = server.PostJSON("/users/1/admin", bytes.NewBufferString(`{"id": 3, "title": "c"}`))
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"user_id":1,"id":3,"title":"c"}
// 	`)
// 	res, err = server.PostJSON("/articles", bytes.NewBufferString(`{"id": 4, "title": "d"}`))
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":4,"title":"d"}
// 	`)
// }

// func TestJSONDelete500(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		import "errors"
// 		type Controller struct {}
// 		func (c *Controller) Delete() (string, error) {
// 			return "", errors.New("Not implemented yet")
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.DeleteJSON("/1", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 500 Internal Server Error
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"error":"Not implemented yet"}
// 	`)
// }

// func TestJSONDelete200(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		type Controller struct {}
// 		type Post struct {
// 			ID int
// 			Title string
// 		}
// 		func (c *Controller) Delete(id int) (*Post, error) {
// 			return &Post{id, "a"}, nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.DeleteJSON("/1", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"ID":1,"Title":"a"}
// 	`)
// }

// func TestJSONMultipleActions(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package controller
// 		type Controller struct {}
// 		func (c *Controller) Index() string {
// 			return "hello world"
// 		}
// 		// Show route
// 		func (c *Controller) Show(id int) int {
// 			return id
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.GetJSON("/")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		"hello world"
// 	`)
// 	res, err = server.GetJSON("/10")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		10
// 	`)
// }

// func TestJSONUpdate500(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		import "errors"
// 		type Controller struct {}
// 		func (c *Controller) Update() (string, error) {
// 			return "", errors.New("Not implemented yet")
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.PatchJSON("/1", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 500 Internal Server Error
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"error":"Not implemented yet"}
// 	`)
// }

// func TestJSONUpdate200(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		type Controller struct {}
// 		type Post struct {
// 			ID int ` + "`" + `json:"id"` + "`" + `
// 			Title string ` + "`" + `json:"title"` + "`" + `
// 		}
// 		func (c *Controller) Update(id int) (*Post, error) {
// 			return &Post{id, "a"}, nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.PatchJSON("/1", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":1,"title":"a"}
// 	`)
// }

// func TestReturnKeyedStruct(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/users/users.go"] = []byte(`
// 		package users
// 		type DB struct {}
// 		type Controller struct {
// 			DB *DB
// 		}
// 		type User struct {
// 			ID   int ` + "`" + `json:"id"` + "`" + `
// 			Name string ` + "`" + `json:"name"` + "`" + `
// 			Age  int ` + "`" + `json:"age"` + "`" + `
// 		}
// 		func (c *Controller) Index() (users []*User, err error) {
// 			users = append(users, &User{1, "a", 2})
// 			users = append(users, &User{2, "b", 3})
// 			return users, nil
// 		}
// 		func (c *Controller) New() {}
// 		func (c *Controller) Create(name string, age int) (user *User, err error) {
// 			return &User{3, name, age}, nil
// 		}
// 		func (c *Controller) Show(id int) (user *User, err error) {
// 			return &User{id, "d", 5}, nil
// 		}
// 		func (c *Controller) Edit(id int) {}
// 		func (c *Controller) Update(id int, name *string, age *int) error {
// 			return nil
// 		}
// 		func (c *Controller) Delete(id int) error {
// 			return nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// }

// func TestNestedResource(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/users/users.go"] = []byte(`
// 		package users
// 		type DB struct {}
// 		type Controller struct {
// 			DB *DB
// 		}
// 		type User struct {
// 			ID   int ` + "`" + `json:"id"` + "`" + `
// 			Name string ` + "`" + `json:"name"` + "`" + `
// 			Age  int ` + "`" + `json:"age"` + "`" + `
// 		}
// 		func (c *Controller) Index() ([]*User, error) {
// 			return []*User{{1, "a", 2}, {2, "b", 3}}, nil
// 		}
// 		func (c *Controller) New() {}
// 		func (c *Controller) Create(name string, age int) (*User, error) {
// 			return &User{3, name, age}, nil
// 		}
// 		func (c *Controller) Show(id int) (*User, error) {
// 			return &User{id, "d", 5}, nil
// 		}
// 		func (c *Controller) Edit(id int) {}
// 		func (c *Controller) Update(id int, name *string, age *int) error {
// 			return nil
// 		}
// 		func (c *Controller) Delete(id int) error {
// 			return nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.GetJSON("/users")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		[{"id":1,"name":"a","age":2},{"id":2,"name":"b","age":3}]
// 	`)
// 	res, err = server.GetJSON("/users/new")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 204 No Content
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	res, err = server.PostJSON("/users?name=matt&age=10", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":3,"name":"matt","age":10}
// 	`)
// 	res, err = server.GetJSON("/users/10")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":10,"name":"d","age":5}
// 	`)
// 	res, err = server.GetJSON("/users/10/edit")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 204 No Content
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	res, err = server.PatchJSON("/users/10", bytes.NewBufferString(`{"name": "matt", "age": 10}`))
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	res, err = server.DeleteJSON("/users/10", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// }

// func TestDeepNestedResource(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/posts/comments/comments.go"] = []byte(`
// 		package comments
// 		type DB struct {}
// 		type Controller struct {
// 			DB *DB
// 		}
// 		type Comment struct {
// 			ID     int ` + "`" + `json:"id,omitempty"` + "`" + `
// 			PostID int ` + "`" + `json:"post_id,omitempty"` + "`" + `
// 			Title  string ` + "`" + `json:"title,omitempty"` + "`" + `
// 		}
// 		func (c *Controller) Index(postID int) ([]*Comment, error) {
// 			return []*Comment{{2, postID, "a"}, {3, postID, "b"}}, nil
// 		}
// 		func (c *Controller) New(postID int) {}
// 		func (c *Controller) Create(postID int, title string) (*Comment, error) {
// 			return &Comment{1, postID, title}, nil
// 		}
// 		func (c *Controller) Show(postID, id int) (*Comment, error) {
// 			return &Comment{id, postID, "a"}, nil
// 		}
// 		func (c *Controller) Edit(postID, id int) {}
// 		func (c *Controller) Update(postID, id int, title *string) (*Comment, error) {
// 			if title == nil {
// 				return &Comment{postID, id, ""}, nil
// 			}
// 			return &Comment{postID, id, *title}, nil
// 		}
// 		func (c *Controller) Delete(postID, id int) (*Comment, error) {
// 			return &Comment{postID, id, ""}, nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.GetJSON("/posts/1/comments")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		[{"id":2,"post_id":1,"title":"a"},{"id":3,"post_id":1,"title":"b"}]
// 	`)
// 	res, err = server.GetJSON("/posts/1/comments/new")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 204 No Content
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	res, err = server.PostJSON("/posts/1/comments", bytes.NewBufferString(`{"title":"1st"}`))
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":1,"post_id":1,"title":"1st"}
// 	`)
// 	res, err = server.GetJSON("/posts/1/comments/2")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":2,"post_id":1,"title":"a"}
// 	`)
// 	res, err = server.GetJSON("/posts/1/comments/2/edit")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 204 No Content
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	res, err = server.PatchJSON("/posts/1/comments/2", bytes.NewBufferString(`{"title":"1st"}`))
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":1,"post_id":2,"title":"1st"}
// 	`)
// 	res, err = server.PatchJSON("/posts/1/comments/2", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":1,"post_id":2}
// 	`)
// 	res, err = server.DeleteJSON("/posts/1/comments/2", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":1,"post_id":2}
// 	`)
// }

// func TestRedirectRootResource(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		type Controller struct {
// 		}
// 		type Post struct {
// 			ID     int ` + "`" + `json:"id,omitempty"` + "`" + `
// 			Title  string ` + "`" + `json:"title,omitempty"` + "`" + `
// 		}
// 		func (c *Controller) Create(title string) (*Post, error) {
// 			return &Post{2, title}, nil
// 		}
// 		func (c *Controller) Update(id int, title *string) (*Post, error) {
// 			if title == nil {
// 				return &Post{id, ""}, nil
// 			}
// 			return &Post{id, *title}, nil
// 		}
// 		func (c *Controller) Delete(id int) (*Post, error) {
// 			return &Post{id, ""}, nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Post("/", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /2
// 	`)
// 	res, err = server.Patch("/1", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /1
// 	`)
// 	res, err = server.Delete("/1", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /
// 	`)
// }

// func TestRedirectNestedResource(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/posts/posts.go"] = []byte(`
// 		package posts
// 		type Controller struct {
// 		}
// 		type Post struct {
// 			ID     int ` + "`" + `json:"id,omitempty"` + "`" + `
// 			Title  string ` + "`" + `json:"title,omitempty"` + "`" + `
// 		}
// 		func (c *Controller) Create(title string) (*Post, error) {
// 			return &Post{2, title}, nil
// 		}
// 		func (c *Controller) Update(id int, title *string) (*Post, error) {
// 			if title == nil {
// 				return &Post{id, ""}, nil
// 			}
// 			return &Post{id, *title}, nil
// 		}
// 		func (c *Controller) Delete(id int) (*Post, error) {
// 			return &Post{id, ""}, nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Post("/posts", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /posts/2
// 	`)
// 	res, err = server.Patch("/posts/1", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /posts/1
// 	`)
// 	res, err = server.Delete("/posts/1", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /posts
// 	`)
// }

// func TestRedirectDeepNestedResource(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/posts/comments/comments.go"] = []byte(`
// 		package comments
// 		type DB struct {}
// 		type Controller struct {
// 			DB *DB
// 		}
// 		type Comment struct {
// 			ID     int ` + "`" + `json:"id,omitempty"` + "`" + `
// 			PostID int ` + "`" + `json:"post_id,omitempty"` + "`" + `
// 			Title  string ` + "`" + `json:"title,omitempty"` + "`" + `
// 		}
// 		func (c *Controller) Index(postID int) ([]*Comment, error) {
// 			return []*Comment{{2, postID, "a"}, {3, postID, "b"}}, nil
// 		}
// 		func (c *Controller) New(postID int) {}
// 		func (c *Controller) Create(postID int, title string) (*Comment, error) {
// 			return &Comment{2, postID, title}, nil
// 		}
// 		func (c *Controller) Show(postID, id int) (*Comment, error) {
// 			return &Comment{id, postID, "a"}, nil
// 		}
// 		func (c *Controller) Edit(postID, id int) {}
// 		func (c *Controller) Update(postID, id int, title *string) (*Comment, error) {
// 			if title == nil {
// 				return &Comment{postID, id, ""}, nil
// 			}
// 			return &Comment{postID, id, *title}, nil
// 		}
// 		func (c *Controller) Delete(postID, id int) (*Comment, error) {
// 			return &Comment{postID, id, ""}, nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Post("/posts/1/comments", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /posts/1/comments/2
// 	`)
// 	res, err = server.Patch("/posts/1/comments/2", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /posts/1/comments/2
// 	`)
// 	res, err = server.Delete("/posts/1/comments/2", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 302 Found
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 		Location: /posts/1/comments
// 	`)
// }

// func TestViewRootResourceUnkeyed(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["view/index.svelte"] = []byte(`
// 		<script>
// 			export let props = []
// 		</script>
// 		{#each props as user}
// 		<h1>index: {user.id} {user.name}</h1>
// 		{/each}
// 	`)
// 	generator.Files["view/new.svelte"] = []byte(`
// 		<script>
// 			export let props = {}
// 		</script>
// 		<h1>new: {props.id} {props.name}</h1>
// 	`)
// 	generator.Files["view/show.svelte"] = []byte(`
// 		<script>
// 			export let props = {}
// 		</script>
// 		<h1>show: {props.id} {props.name}</h1>
// 	`)
// 	generator.Files["view/edit.svelte"] = []byte(`
// 		<script>
// 			export let props = {}
// 		</script>
// 		<h1>edit: {props.id} {props.name}</h1>
// 	`)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action
// 		type Controller struct {}
// 		type User struct {
// 			ID int ` + "`" + `json:"id"` + "`" + `
// 			Name string ` + "`" + `json:"name"` + "`" + `
// 		}
// 		func (c *Controller) Index() []*User {
// 			return []*User{{1, "a"}, {2, "b"}}
// 		}
// 		func (c *Controller) New() *User {
// 			return &User{3, "c"}
// 		}
// 		func (c *Controller) Show(id int) *User {
// 			return &User{id, "s"}
// 		}
// 		func (c *Controller) Edit(id int) *User {
// 			return &User{id, "e"}
// 		}
// 	`)
// 	// Generate the app
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()

// 	//
// 	res, err := server.GetJSON("")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		[{"id":1,"name":"a"},{"id":2,"name":"b"}]
// 	`)
// 	res, err = server.Get("")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el := res.Query("#bud_target")
// 	html, err := el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>index: 1 a</h1><h1>index: 2 b</h1>`, html)

// 	// /new
// 	res, err = server.GetJSON("/new")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":3,"name":"c"}
// 	`)
// 	res, err = server.Get("/new")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el = res.Query("#bud_target")
// 	html, err = el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>new: 3 c</h1>`, html)

// 	// /:id
// 	res, err = server.GetJSON("/10")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":10,"name":"s"}
// 	`)
// 	res, err = server.Get("/10")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el = res.Query("#bud_target")
// 	html, err = el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>show: 10 s</h1>`, html)

// 	// /:id/edit
// 	res, err = server.GetJSON("/10/edit")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":10,"name":"e"}
// 	`)
// 	res, err = server.Get("/10/edit")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el = res.Query("#bud_target")
// 	html, err = el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>edit: 10 e</h1>`, html)
// }

// func TestViewNestedResourceUnkeyed(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["view/users/index.svelte"] = []byte(`
// 		<script>
// 			export let props = []
// 		</script>
// 		{#each props as user}
// 		<h1>index: {user.id} {user.name}</h1>
// 		{/each}
// 	`)
// 	generator.Files["view/users/new.svelte"] = []byte(`
// 		<script>
// 			export let props = {}
// 		</script>
// 		<h1>new: {props.id} {props.name}</h1>
// 	`)
// 	generator.Files["view/users/show.svelte"] = []byte(`
// 		<script>
// 			export let props = {}
// 		</script>
// 		<h1>show: {props.id} {props.name}</h1>
// 	`)
// 	generator.Files["view/users/edit.svelte"] = []byte(`
// 		<script>
// 			export let props = {}
// 		</script>
// 		<h1>edit: {props.id} {props.name}</h1>
// 	`)
// 	generator.Files["action/users/users.go"] = []byte(`
// 		package users
// 		type Controller struct {}
// 		type User struct {
// 			ID int ` + "`" + `json:"id"` + "`" + `
// 			Name string ` + "`" + `json:"name"` + "`" + `
// 		}
// 		func (c *Controller) Index() []*User {
// 			return []*User{{1, "a"}, {2, "b"}}
// 		}
// 		func (c *Controller) New() *User {
// 			return &User{3, "c"}
// 		}
// 		func (c *Controller) Show(id int) *User {
// 			return &User{id, "s"}
// 		}
// 		func (c *Controller) Edit(id int) *User {
// 			return &User{id, "e"}
// 		}
// 	`)
// 	// Generate the app
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()

// 	// /users
// 	res, err := server.GetJSON("/users")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		[{"id":1,"name":"a"},{"id":2,"name":"b"}]
// 	`)
// 	res, err = server.Get("/users")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el := res.Query("#bud_target")
// 	html, err := el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>index: 1 a</h1><h1>index: 2 b</h1>`, html)

// 	// /users/new
// 	res, err = server.GetJSON("/users/new")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":3,"name":"c"}
// 	`)
// 	res, err = server.Get("/users/new")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el = res.Query("#bud_target")
// 	html, err = el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>new: 3 c</h1>`, html)

// 	// /users/:id
// 	res, err = server.GetJSON("/users/10")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":10,"name":"s"}
// 	`)
// 	res, err = server.Get("/users/10")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el = res.Query("#bud_target")
// 	html, err = el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>show: 10 s</h1>`, html)

// 	// /users/:id/edit
// 	res, err = server.GetJSON("/users/10/edit")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":10,"name":"e"}
// 	`)
// 	res, err = server.Get("/users/10/edit")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el = res.Query("#bud_target")
// 	html, err = el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>edit: 10 e</h1>`, html)
// }

// func TestViewDeepResourceUnkeyed(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["view/teams/users/index.svelte"] = []byte(`
// 		<script>
// 			export let props = []
// 		</script>
// 		{#each props as user}
// 		<h1>index: {user.id} {user.name}</h1>
// 		{/each}
// 	`)
// 	generator.Files["view/teams/users/new.svelte"] = []byte(`
// 		<script>
// 			export let props = {}
// 		</script>
// 		<h1>new: {props.id} {props.name}</h1>
// 	`)
// 	generator.Files["view/teams/users/show.svelte"] = []byte(`
// 		<script>
// 			export let props = {}
// 		</script>
// 		<h1>show: {props.id} {props.name}</h1>
// 	`)
// 	generator.Files["view/teams/users/edit.svelte"] = []byte(`
// 		<script>
// 			export let props = {}
// 		</script>
// 		<h1>edit: {props.id} {props.name}</h1>
// 	`)
// 	generator.Files["action/teams/users/users.go"] = []byte(`
// 		package users
// 		type Controller struct {}
// 		type User struct {
// 			ID int ` + "`" + `json:"id"` + "`" + `
// 			Name string ` + "`" + `json:"name"` + "`" + `
// 		}
// 		func (c *Controller) Index() []*User {
// 			return []*User{{1, "a"}, {2, "b"}}
// 		}
// 		func (c *Controller) New() *User {
// 			return &User{3, "c"}
// 		}
// 		func (c *Controller) Show(id int) *User {
// 			return &User{id, "s"}
// 		}
// 		func (c *Controller) Edit(id int) *User {
// 			return &User{id, "e"}
// 		}
// 	`)
// 	// Generate the app
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()

// 	// /teams/:team_id/users
// 	res, err := server.GetJSON("/teams/5/users")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		[{"id":1,"name":"a"},{"id":2,"name":"b"}]
// 	`)
// 	res, err = server.Get("/teams/5/users")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el := res.Query("#bud_target")
// 	html, err := el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>index: 1 a</h1><h1>index: 2 b</h1>`, html)

// 	// /teams/:team_id/users/new
// 	res, err = server.GetJSON("/teams/5/users/new")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":3,"name":"c"}
// 	`)
// 	res, err = server.Get("/teams/5/users/new")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el = res.Query("#bud_target")
// 	html, err = el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>new: 3 c</h1>`, html)

// 	// /teams/:team_id/users/:id
// 	res, err = server.GetJSON("/teams/5/users/10")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":10,"name":"s"}
// 	`)
// 	res, err = server.Get("/teams/5/users/10")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el = res.Query("#bud_target")
// 	html, err = el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>show: 10 s</h1>`, html)

// 	// /teams/:team_id/users/:id/edit
// 	res, err = server.GetJSON("/teams/5/users/10/edit")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":10,"name":"e"}
// 	`)
// 	res, err = server.Get("/teams/5/users/10/edit")
// 	is.NoErr(err)
// 	res.ExpectHeaders(`
// 		HTTP/1.1 200 OK
// 		Content-Type: text/html
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT
// 	`)
// 	el = res.Query("#bud_target")
// 	html, err = el.Html()
// 	is.NoErr(err)
// 	is.Equal(`<h1>edit: 10 e</h1>`, html)
// }

// func TestResourceContext(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/users/users.go"] = []byte(`
// 		package users
// 		import contexts "context"
// 		type User struct {
// 			ID int ` + "`" + `json:"id"` + "`" + `
// 			Name string ` + "`" + `json:"name"` + "`" + `
// 		}
// 		type Controller struct {
// 		}
// 		func (c *Controller) Index(ctx contexts.Context) []*User {
// 			return []*User{{1, "a"}, {2, "b"}}
// 		}
// 		func (c *Controller) Create(ctx contexts.Context, name string) *User {
// 			return &User{1, name}
// 		}
// 		func (c *Controller) Show(ctx contexts.Context, id int) *User {
// 			return &User{id, "a"}
// 		}
// 		func (c *Controller) Edit(ctx contexts.Context, id int) *User {
// 			return &User{id, "a"}
// 		}
// 		func (c *Controller) Update(ctx contexts.Context, id int, name string) (*User, error) {
// 			return &User{id, name}, nil
// 		}
// 		func (c *Controller) Delete(ctx contexts.Context, id int) (*User, error) {
// 			return &User{id, "a"}, nil
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.GetJSON("/users")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		[{"id":1,"name":"a"},{"id":2,"name":"b"}]
// 	`)
// 	res, err = server.PostJSON("/users", bytes.NewBufferString(`{"name":"b"}`))
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":1,"name":"b"}
// 	`)
// 	res, err = server.GetJSON("/users/2")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":2,"name":"a"}
// 	`)
// 	res, err = server.GetJSON("/users/2/edit")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":2,"name":"a"}
// 	`)
// 	res, err = server.PatchJSON("/users/2", bytes.NewBufferString(`{"name":"b"}`))
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":2,"name":"b"}
// 	`)
// 	res, err = server.DeleteJSON("/users/2", nil)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		{"id":2,"name":"a"}
// 	`)
// }

// func TestWorkingThenWorking(t *testing.T) {
// 	is := is.New(t)
// 	generator := test.Generator(t)
// 	generator.Files["action/action.go"] = []byte(`
// 		package action

// 		type Controller struct {
// 		}

// 		func (c *Controller) Index() string {
// 			return "Hello Users!"
// 		}
// 	`)
// 	fsCache := fscache.New()
// 	app, err := generator.Generate(fsCache)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	server, err := app.Start()
// 	is.NoErr(err)
// 	defer server.Close()
// 	res, err := server.Get("/")
// 	is.NoErr(err)
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		"Hello Users!"
// 	`)
// 	ctx := context.Background()
// 	err = app.Regenerate(ctx)
// 	is.NoErr(err)
// 	is.True(app.Exists("bud/action/action.go"))
// 	is.True(app.Exists("bud/main.go"))
// 	err = server.Restart()
// 	is.NoErr(err)
// 	res, err = server.Get("/")
// 	is.NoErr(err)
// 	res.Expect(`
// 		HTTP/1.1 200 OK
// 		Content-Type: application/json
// 		Date: Fri, 31 Dec 2021 00:00:00 GMT

// 		"Hello Users!"
// 	`)
// }
