package controller_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"testing"
	"testing/fstest"
	"time"

	"github.com/livebud/bud/pkg/controller"
	"github.com/livebud/bud/pkg/di"
	"github.com/livebud/bud/pkg/middleware/dim"
	"github.com/livebud/bud/pkg/mux"
	"github.com/livebud/bud/pkg/view"
	"github.com/livebud/bud/pkg/view/gohtml"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
)

func equal(t testing.TB, h http.Handler, r *http.Request, expect string) {
	t.Helper()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	w := rec.Result()
	dump, err := httputil.DumpResponse(w, true)
	if err != nil {
		if err.Error() != expect {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	diff.TestHTTP(t, expect, string(dump))
}

func emptyViewer() view.Finder {
	return view.New(fstest.MapFS{}, map[string]view.Renderer{})
}

func TestNoActions(t *testing.T) {
	is := is.New(t)
	controller := controller.New(emptyViewer(), mux.New())
	var root struct{}
	is.NoErr(controller.Register("/", &root))
	req := httptest.NewRequest("GET", "/", nil)
	equal(t, controller, req, `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
}

type indexRoot struct{}

type indexRootOut struct {
	Message string `json:"message"`
}

func (c *indexRoot) Index() *indexRootOut {
	return &indexRootOut{"Hello"}
}

type indexAbout struct{}

type indexAboutOut struct {
	Message string `json:"message"`
}

func (c *indexAbout) Index() indexAboutOut {
	return indexAboutOut{"About"}
}

type indexComments struct{}

type indexComment struct {
	PostID  int    `json:"post_id"`
	Message string `json:"message"`
}

func (c *indexComments) Index(in *indexComment) (*indexComment, error) {
	return in, nil
}

func TestIndex(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"index.gohtml":                &fstest.MapFile{Data: []byte("{{.Message}}")},
		"about.gohtml":                &fstest.MapFile{Data: []byte("{{.Message}}")},
		"posts/comments/index.gohtml": &fstest.MapFile{Data: []byte("id={{.PostID}} message={{.Message}}")},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &indexRoot{}))
	is.NoErr(controller.Register("/about", &indexAbout{}))
	is.NoErr(controller.Register("/posts/{post_id}/comments", &indexComments{}))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"message":"Hello"}
	`)
	req = httptest.NewRequest("GET", "/about", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"message":"About"}
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments?message=comments", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":1,"message":"comments"}
	`)
	req = httptest.NewRequest("GET", "/", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		Hello
	`)
	req = httptest.NewRequest("GET", "/about", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		About
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments?message=comments", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		id=1 message=comments
	`)
}

type showRoot struct{}

type showRootIn struct {
	ID int `json:"id"`
}

type showRootOut struct {
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func (c *showRoot) Show(in showRootIn) *showRootOut {
	return &showRootOut{in.ID, "Show"}
}

type showAbout struct{}

type showAboutOut struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (c *showAbout) Show(in struct {
	ID string `json:"id"`
}) showAboutOut {
	return showAboutOut{in.ID, "About"}
}

type showComments struct{}

type showComment struct {
	PostID  int    `json:"post_id"`
	ID      int    `json:"id"`
	Message string `json:"message"`
}

func (c *showComments) Show(in *showComment) (*showComment, error) {
	return in, nil
}

func TestShow(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"show.gohtml":                &fstest.MapFile{Data: []byte("{{.Message}}")},
		"error.gohtml":               &fstest.MapFile{Data: []byte("{{.Error}}")},
		"about/show.gohtml":          &fstest.MapFile{Data: []byte("id={{.ID}} message={{.Message}}")},
		"posts/comments/show.gohtml": &fstest.MapFile{Data: []byte("post_id={{.PostID}} id={{.ID}} message={{.Message}}")},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &showRoot{}))
	is.NoErr(controller.Register("/about", &showAbout{}))
	is.NoErr(controller.Register("/posts/{post_id}/comments", &showComments{}))
	req := httptest.NewRequest("GET", "/10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":10,"message":"Show"}
	`)
	// Bad request
	req = httptest.NewRequest("GET", "/abc", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 400 Bad Request
		Connection: close
		Content-Type: application/json

		{"error":"could not parse int from \"abc\""}
	`)
	req = httptest.NewRequest("GET", "/about/abc", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":"abc","message":"About"}
	`)
	req = httptest.NewRequest("GET", "/posts/10/comments/20?message=comment", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":10,"id":20,"message":"comment"}
	`)
	req = httptest.NewRequest("GET", "/10", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		Show
	`)
	// Bad request
	req = httptest.NewRequest("GET", "/abc", nil)
	equal(t, controller, req, `
		HTTP/1.1 400 Bad Request
		Connection: close
		Content-Type: text/html

		could not parse int from &#34;abc&#34;
	`)
	req = httptest.NewRequest("GET", "/about/abc", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		id=abc message=About
	`)
	req = httptest.NewRequest("GET", "/posts/10/comments/20?message=comment", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		post_id=10 id=20 message=comment
	`)
}

type createRoot struct{}

func (c *createRoot) Create() {
}

type createUsers struct{}

type createUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (c *createUsers) Create(in createUser) (*createUser, error) {
	if in.Name == "" {
		return nil, errors.New("missing name")
	}
	in.ID = 1
	return &in, nil
}

type createComments struct{}

type createComment struct {
	PostID  int    `json:"post_id"`
	ID      int    `json:"id"`
	Comment string `json:"comment"`
}

func (c *createComments) Create(ctx context.Context, in *createComment) (*createComment, error) {
	if in.PostID == 0 {
		return nil, errors.New("missing post_id")
	} else if in.Comment == "" {
		return nil, errors.New("missing comment")
	}
	in.ID = 1
	return in, nil
}

func TestCreate(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{}
	viewer := view.New(fsys, map[string]view.Renderer{})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &createRoot{}))
	is.NoErr(controller.Register("/users", &createUsers{}))
	is.NoErr(controller.Register("/posts/{post_id}/comments", &createComments{}))
	// JSON creates
	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("POST", "/users?name=max", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":1,"name":"max"}
	`)
	// Error
	req = httptest.NewRequest("POST", "/users", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"missing name"}
	`)
	req = httptest.NewRequest("POST", "/posts/10/comments", bytes.NewBufferString(`{"comment":"howdy"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":10,"id":1,"comment":"howdy"}
	`)
	// Error
	req = httptest.NewRequest("POST", "/posts/20/comments", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"missing comment"}
	`)
	// Error with Referer header
	req = httptest.NewRequest("POST", "/posts/20/comments?comment=howdy", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "/posts/20")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":20,"id":1,"comment":"howdy"}
	`)
	// HTML creates
	req = httptest.NewRequest("POST", "/", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /
	`)
	req = httptest.NewRequest("POST", "/users?name=mark", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users
	`)
	// Create Error
	req = httptest.NewRequest("POST", "/users", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users

		missing name
	`)
	req = httptest.NewRequest("POST", "/posts/10/comments?comment=hi", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/10/comments
	`)
	// Create Error
	req = httptest.NewRequest("POST", "/posts/20/comments", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/20/comments

		missing comment
	`)
	// Create Error with Referer header
	req = httptest.NewRequest("POST", "/posts/20/comments", nil)
	req.Header.Set("Referer", "/posts/20")
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/20

		missing comment
	`)
}

type updateRoot struct{}

func (c *updateRoot) Update() {
}

type updateUsers struct{}

type updateUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *updateUsers) Update(in *updateUser) (updateUser, error) {
	if in.Name == "" {
		return updateUser{}, errors.New("missing name")
	}
	in.Name = in.Name + "y"
	return *in, nil
}

type updateComments struct{}

type updateComment struct {
	PostID int    `json:"post_id"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
}

func (c *updateComments) Update(ctx context.Context, in *updateComment) (*updateComment, error) {
	if in.PostID == 0 {
		return nil, errors.New("invalid post id")
	} else if in.ID == 0 {
		return nil, errors.New("invalid comment id")
	} else if in.Title == "" {
		return nil, errors.New("missing title")
	}
	in.Title += "!"
	return in, nil
}

func TestUpdate(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{}
	viewer := view.New(fsys, map[string]view.Renderer{})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &updateRoot{}))
	is.NoErr(controller.Register("/users", &updateUsers{}))
	is.NoErr(controller.Register("/posts/{post_id}/comments", &updateComments{}))
	// JSON
	req := httptest.NewRequest("PATCH", "/10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("PATCH", "/users/abc?name=max", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":"abc","name":"maxy"}
	`)
	// Coerce to string
	req = httptest.NewRequest("PATCH", "/users/10?name=ank", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":"10","name":"anky"}
	`)
	req = httptest.NewRequest("PATCH", "/posts/20/comments/10?title=cool", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":20,"id":10,"title":"cool!"}
	`)
	// Bad request
	req = httptest.NewRequest("PATCH", "/posts/1/comments/def", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 400 Bad Request
		Connection: close
		Content-Type: application/json

		{"error":"could not parse int from \"def\""}
	`)
	// Error
	req = httptest.NewRequest("PATCH", "/posts/20/comments/10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"missing title"}
	`)
	// Error with Referer header
	req = httptest.NewRequest("PATCH", "/posts/20/comments/10", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "/posts/20")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"missing title"}
	`)
	// HTML
	req = httptest.NewRequest("PATCH", "/10", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /10
	`)
	req = httptest.NewRequest("PATCH", "/users/abc?name=max", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users/abc
	`)
	// Coerce to string
	req = httptest.NewRequest("PATCH", "/users/10?name=ank", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users/10
	`)
	req = httptest.NewRequest("PATCH", "/posts/20/comments/10?title=cool", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/20/comments/10
	`)
	// Bad request
	// req = httptest.NewRequest("PATCH", "/posts/abc/comments/def", nil)
	// equal(t, controller, req, `
	// 	HTTP/1.1 303 See Other
	// 	Connection: close
	// 	Location: /posts/abc/comments/def

	// 	could not parse int from "abc"
	// `)
	// Error
	req = httptest.NewRequest("PATCH", "/posts/20/comments/10", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/20/comments/10

		missing title
	`)
	// Error with Referer header
	req = httptest.NewRequest("PATCH", "/posts/20/comments/10", nil)
	req.Header.Set("Referer", "/posts/20")
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/20

		missing title
	`)
}

type deleteRoot struct{}

func (c *deleteRoot) Delete() {
}

type deleteUsers struct{}

type deleteUser struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

func (c *deleteUsers) Delete(in *deleteUser) (deleteUser, error) {
	if in.Deleted {
		return deleteUser{}, errors.New("already deleted")
	}
	in.Deleted = true
	return *in, nil
}

type deleteComments struct{}

type deleteComment struct {
	PostID  int  `json:"post_id"`
	ID      int  `json:"id"`
	Deleted bool `json:"deleted"`
}

func (c *deleteComments) Delete(ctx context.Context, in *deleteComment) (*deleteComment, error) {
	if in.PostID == 0 {
		return nil, errors.New("invalid post id")
	} else if in.ID == 0 {
		return nil, errors.New("invalid comment id")
	} else if in.Deleted {
		return nil, errors.New("already deleted")
	}
	in.Deleted = true
	return in, nil
}

func TestDelete(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{}
	viewer := view.New(fsys, map[string]view.Renderer{})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &deleteRoot{}))
	is.NoErr(controller.Register("/users", &deleteUsers{}))
	is.NoErr(controller.Register("/posts/{post_id}/comments", &deleteComments{}))
	// JSON
	req := httptest.NewRequest("DELETE", "/10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("DELETE", "/users/abc?name=max", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":"abc","deleted":true}
	`)
	// Error
	req = httptest.NewRequest("DELETE", "/users/abc?deleted=true", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"already deleted"}
	`)
	req = httptest.NewRequest("DELETE", "/posts/20/comments/10?title=cool", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":20,"id":10,"deleted":true}
	`)
	// Bad request
	req = httptest.NewRequest("DELETE", "/posts/1/comments/def", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 400 Bad Request
		Connection: close
		Content-Type: application/json

		{"error":"could not parse int from \"def\""}
	`)
	// Error
	req = httptest.NewRequest("DELETE", "/posts/20/comments/10?deleted=true", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"already deleted"}
	`)
	// Error with Referer header
	req = httptest.NewRequest("DELETE", "/posts/20/comments/10?deleted=true", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "/posts/20")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"already deleted"}
	`)
	// HTML
	req = httptest.NewRequest("DELETE", "/10", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /10
	`)
	req = httptest.NewRequest("DELETE", "/users/abc", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users/abc
	`)
	// Error
	req = httptest.NewRequest("DELETE", "/users/10?deleted=true", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users/10

		already deleted
	`)
	req = httptest.NewRequest("DELETE", "/posts/20/comments/10", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/20/comments/10
	`)
	// Bad request
	// req = httptest.NewRequest("DELETE", "/posts/abc/comments/def", nil)
	// equal(t, controller, req, `
	// 	HTTP/1.1 303 See Other
	// 	Connection: close
	// 	Location: /posts/abc/comments/def

	// 	could not parse int from "def"
	// `)
	// Error
	req = httptest.NewRequest("DELETE", "/posts/20/comments/10?deleted=true", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/20/comments/10

		already deleted
	`)
	// Error with Referer header
	req = httptest.NewRequest("DELETE", "/posts/20/comments/10?deleted=true", nil)
	req.Header.Set("Referer", "/posts/20")
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/20

		already deleted
	`)
}

type editRoot struct{}

func (c *editRoot) Edit() {
}

type editUsers struct{}

type editUser struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Filter string `json:"filter,omitempty"`
}

func (c *editUsers) Edit(in *editUser) (editUser, error) {
	if in.ID == "20" {
		return editUser{}, errors.New("user not found")
	}
	in.Name = "marc"
	return *in, nil
}

type editComments struct{}

type editComment struct {
	PostID int    `json:"post_id"`
	ID     int    `json:"id"`
	Title  string `json:"title"`
}

func (c *editComments) Edit(ctx context.Context, in *editComment) (*editComment, error) {
	if in.PostID == 0 {
		return nil, errors.New("invalid post id")
	} else if in.ID == 0 {
		return nil, errors.New("invalid comment id")
	} else if in.Title != "" {
		return nil, errors.New("can't preset title")
	}
	in.Title = "hello world"
	return in, nil
}

func TestEdit(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"edit.gohtml":                &fstest.MapFile{Data: []byte("Edit Root")},
		"users/edit.gohtml":          &fstest.MapFile{Data: []byte(`Edit {{ .ID }} name={{.Name}}`)},
		"posts/comments/edit.gohtml": &fstest.MapFile{Data: []byte(`Edit {{.PostID}}/{{.ID}} title={{.Title}}`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &editRoot{}))
	is.NoErr(controller.Register("/users", &editUsers{}))
	is.NoErr(controller.Register("/posts/{post_id}/comments", &editComments{}))
	// JSON
	req := httptest.NewRequest("GET", "/10/edit", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("GET", "/users/abc/edit?filter=customer", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":"abc","name":"marc","filter":"customer"}
	`)
	// Omit empty
	req = httptest.NewRequest("GET", "/users/10/edit", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":"10","name":"marc"}
	`)
	// Error
	req = httptest.NewRequest("GET", "/users/20/edit", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"user not found"}
	`)
	req = httptest.NewRequest("GET", "/posts/20/comments/10/edit", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":20,"id":10,"title":"hello world"}
	`)
	// Bad request
	req = httptest.NewRequest("GET", "/posts/1/comments/def/edit", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 400 Bad Request
		Connection: close
		Content-Type: application/json

		{"error":"could not parse int from \"def\""}
	`)
	// Error
	req = httptest.NewRequest("GET", "/posts/20/comments/10/edit?title=set", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"can't preset title"}
	`)
	// Error with Referer header
	req = httptest.NewRequest("GET", "/posts/20/comments/10/edit?title=set", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "/posts/20")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"can't preset title"}
	`)
	// HTML
	req = httptest.NewRequest("GET", "/10/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		Edit Root
	`)
	req = httptest.NewRequest("GET", "/users/abc/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		Edit abc name=marc
	`)
	// Coerce to string
	req = httptest.NewRequest("GET", "/users/10/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		Edit 10 name=marc
	`)
	// Error
	req = httptest.NewRequest("GET", "/users/20/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		user not found
	`)
	req = httptest.NewRequest("GET", "/posts/20/comments/10/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		Edit 20/10 title=hello world
	`)
	// Bad request
	req = httptest.NewRequest("GET", "/posts/10/comments/def/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 400 Bad Request
		Connection: close
		Content-Type: text/html

		could not parse int from "def"
	`)
	// Error
	req = httptest.NewRequest("GET", "/posts/20/comments/10/edit?title=set", nil)
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		can't preset title
	`)
	// Error with Referer header
	req = httptest.NewRequest("GET", "/posts/20/comments/10/edit?title=set", nil)
	req.Header.Set("Referer", "/posts/20")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		can't preset title
	`)
}

type newRoot struct{}

func (c *newRoot) New() {
}

type newUsers struct{}

type newUser struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (c *newUsers) New(in *newUser) (newUser, error) {
	if in.Type == "" {
		return newUser{}, errors.New("missing type")
	}
	in.Name = "enter name"
	return *in, nil
}

type newComments struct{}

type newComment struct {
	PostID int    `json:"post_id"`
	Title  string `json:"title"`
}

func (c *newComments) New(ctx context.Context, in *newComment) (*newComment, error) {
	if in.PostID == 0 {
		return nil, errors.New("invalid post id")
	}
	in.Title = "hello world"
	return in, nil
}

func TestNew(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"new.gohtml":                &fstest.MapFile{Data: []byte("New root")},
		"users/new.gohtml":          &fstest.MapFile{Data: []byte(`New name={{.Name}} type={{.Type}}`)},
		"posts/comments/new.gohtml": &fstest.MapFile{Data: []byte(`New {{.PostID}} title={{.Title}}`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &newRoot{}))
	is.NoErr(controller.Register("/users", &newUsers{}))
	is.NoErr(controller.Register("/posts/{post_id}/comments", &newComments{}))
	// JSON
	req := httptest.NewRequest("GET", "/new", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("GET", "/users/new?type=customer", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"name":"enter name","type":"customer"}
	`)
	// Error
	req = httptest.NewRequest("GET", "/users/new", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"missing type"}
	`)
	req = httptest.NewRequest("GET", "/posts/20/comments/new", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":20,"title":"hello world"}
	`)
	// Bad request
	req = httptest.NewRequest("GET", "/posts/abc/comments/new", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 400 Bad Request
		Connection: close
		Content-Type: application/json

		{"error":"could not parse int from \"abc\""}
	`)
	// Error
	req = httptest.NewRequest("GET", "/posts/0/comments/new", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"invalid post id"}
	`)
	// Error with Referer header
	req = httptest.NewRequest("GET", "/posts/0/comments/new", nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", "/posts/20")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"invalid post id"}
	`)
	// HTML
	req = httptest.NewRequest("GET", "/new", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		New root
	`)
	req = httptest.NewRequest("GET", "/users/new?type=customer", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		New name=enter name type=customer
	`)
	// Error
	req = httptest.NewRequest("GET", "/users/new", nil)
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		missing type
	`)
	req = httptest.NewRequest("GET", "/posts/20/comments/new", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		New 20 title=hello world
	`)
	// Bad request
	// req = httptest.NewRequest("GET", "/posts/abc/comments/new", nil)
	// equal(t, controller, req, `
	// 	HTTP/1.1 400 Bad Request
	// 	Connection: close
	// 	Content-Type: text/html

	// 	could not parse int from "abc"
	// `)
	// Error
	req = httptest.NewRequest("GET", "/posts/0/comments/new", nil)
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		invalid post id
	`)
	// Error with Referer header
	req = httptest.NewRequest("GET", "/posts/0/comments/new", nil)
	req.Header.Set("Referer", "/posts")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		invalid post id
	`)
}

type noContentRoot struct{}

func (c *noContentRoot) Index() {
}
func (c *noContentRoot) Show(in struct {
	ID string `json:"id"`
}) {
}
func (c *noContentRoot) New() {}
func (c *noContentRoot) Edit(in struct {
	ID string `json:"id"`
}) {
}

type noContentComments struct{}

func (c *noContentComments) Index() {}
func (c *noContentComments) Show(in struct {
	PostID string `json:"post_id"`
	ID     string `json:"id"`
}) {
}
func (c *noContentComments) New() {}
func (c *noContentComments) Edit(in struct {
	PostID string `json:"post_id"`
	ID     string `json:"id"`
}) {
}

func TestNoContent(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"index.html":                &fstest.MapFile{Data: []byte("root index")},
		"show.html":                 &fstest.MapFile{Data: []byte("root show")},
		"new.html":                  &fstest.MapFile{Data: []byte("root new")},
		"edit.html":                 &fstest.MapFile{Data: []byte("root edit")},
		"posts/comments/index.html": &fstest.MapFile{Data: []byte("comments index")},
		"posts/comments/show.html":  &fstest.MapFile{Data: []byte("comments show")},
		"posts/comments/new.html":   &fstest.MapFile{Data: []byte("comments new")},
		"posts/comments/edit.html":  &fstest.MapFile{Data: []byte("comments edit")},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".html": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &noContentRoot{}))
	is.NoErr(controller.Register("/posts/{post_id}/comments", &noContentComments{}))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("GET", "/10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("GET", "/new", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("GET", "/10/edit", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/new", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/10/edit", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("GET", "/", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		root index
	`)
	req = httptest.NewRequest("GET", "/10", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		root show
	`)
	req = httptest.NewRequest("GET", "/new", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		root new
	`)
	req = httptest.NewRequest("GET", "/10/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		root edit
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		comments index
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/5", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		comments show
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/new", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		comments new
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/10/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		comments edit
	`)
}

func TestNoView(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{}
	viewer := view.New(fsys, map[string]view.Renderer{})
	controller := controller.New(viewer, mux.New())
	err := controller.Get("/", func() {})
	is.True(err != nil)
	is.True(errors.Is(err, view.ErrNotFound))
}

type listPosts struct{}

type listPostsOut struct {
	Posts []listPost `json:"posts"`
}

type listPost struct {
	Title string `json:"title"`
}

func (c *listPosts) Index() listPostsOut {
	return listPostsOut{
		Posts: []listPost{
			{Title: "First Post"},
			{Title: "Second Post"},
		},
	}
}

type listComments struct{}

type listCommentsIn struct {
	PostID string `json:"postid"`
}

type listCommentsOut struct {
	Comments []*listComment `json:"comments"`
}

type listComment struct {
	Comment string `json:"comment"`
}

func (c *listComments) Index(ctx context.Context, in *listCommentsIn) (*listCommentsOut, error) {
	if in.PostID != "1" {
		return nil, errors.New("post not found")
	}
	return &listCommentsOut{
		Comments: []*listComment{
			{Comment: "First Comment"},
			{Comment: "Second Comment"},
		},
	}, nil
}

func TestList(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"posts/index.gohtml":          &fstest.MapFile{Data: []byte(`{{range .Posts}}<p>{{.Title}}</p>{{end}}`)},
		"posts/comments/index.gohtml": &fstest.MapFile{Data: []byte(`{{range .Comments}}<p>{{.Comment}}</p>{{end}}`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/posts", &listPosts{}))
	is.NoErr(controller.Register("/posts/{postid}/comments", &listComments{}))
	req := httptest.NewRequest("GET", "/posts", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"posts":[{"title":"First Post"},{"title":"Second Post"}]}
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"comments":[{"comment":"First Comment"},{"comment":"Second Comment"}]}
	`)
	req = httptest.NewRequest("GET", "/posts/2/comments", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"post not found"}
	`)
	// HTML
	req = httptest.NewRequest("GET", "/posts", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		<p>First Post</p><p>Second Post</p>
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		<p>First Comment</p><p>Second Comment</p>
	`)
	req = httptest.NewRequest("GET", "/posts/2/comments", nil)
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		post not found
	`)
}

type depRoot struct {
}

type depPool struct {
	URL string
}

type depContext struct {
	context.Context
	Pool *depPool
}

type depOut struct {
	URL string
}

func (c *depRoot) Index(ctx *depContext) depOut {
	return depOut{URL: ctx.Pool.URL}
}

func TestDepUsingRequest(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"index.gohtml": &fstest.MapFile{Data: []byte(`{{ .URL }}`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	in := di.New()
	di.Provide[*depPool](in, func(req *http.Request) (*depPool, error) {
		return &depPool{
			URL: req.URL.Path,
		}, nil
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &depRoot{}))
	req := httptest.NewRequest("GET", "/", nil)
	equal(t, dim.Provide(in).Middleware(controller), req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		/
	`)
}

type shareArticle struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type shareRoot struct {
}

func (c *shareRoot) Index(a *shareArticle) *shareArticle {
	return a
}

func TestShareStruct(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"index.gohtml": &fstest.MapFile{Data: []byte(`index`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &shareRoot{}))
	req := httptest.NewRequest("GET", "/?id=10&title=first", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":10,"title":"first"}
	`)
}

type usersResource struct{}

type userResource struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (c *usersResource) Index() ([]*userResource, error) {
	return []*userResource{{1, "a", 2}, {2, "b", 3}}, nil
}

func (c *usersResource) New() {}

type usersCreate struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (c *usersResource) Create(in *usersCreate) (*userResource, error) {
	if in.Name == "" {
		return nil, errors.New("missing name")
	} else if in.Age == 0 {
		return nil, errors.New("missing age")
	}
	return &userResource{3, in.Name, in.Age}, nil
}

type usersShow struct {
	ID int `json:"id"`
}

func (c *usersResource) Show(in *usersShow) (*userResource, error) {
	return &userResource{in.ID, "d", 5}, nil
}

type usersEdit struct {
	ID int `json:"id"`
}

func (c *usersResource) Edit(in *usersEdit) error {
	if in.ID == 0 {
		return fmt.Errorf("id not found")
	}
	return nil
}

type usersUpdate struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
	Age  *int    `json:"age"`
}

func (c *usersResource) Update(in *usersUpdate) (*userResource, error) {
	if in.ID == 0 {
		return nil, errors.New("missing id")
	}
	user := new(userResource)
	user.ID = in.ID
	if in.Age != nil {
		user.Age = *in.Age
	}
	if in.Name != nil {
		user.Name = *in.Name
	}
	return user, nil
}

type usersDelete struct {
	ID int `json:"id"`
}

func (c *usersResource) Delete(in *usersDelete) error {
	if in.ID == 0 {
		return errors.New("missing id")
	}
	return nil
}

type helloOut struct {
	Message string
}

func (c *usersResource) Hello() (out helloOut) {
	out.Message = "hello"
	return out
}

func TestResource(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"users/index.gohtml": &fstest.MapFile{Data: []byte(`index`)},
		"users/show.gohtml":  &fstest.MapFile{Data: []byte(`show`)},
		"users/new.gohtml":   &fstest.MapFile{Data: []byte(`new`)},
		"users/edit.gohtml":  &fstest.MapFile{Data: []byte(`edit`)},
		"users/hello.gohtml": &fstest.MapFile{Data: []byte(`hello`)},
		"users/error.gohtml": &fstest.MapFile{Data: []byte(`error page: {{ .Error }}`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/users", &usersResource{}))
	req := httptest.NewRequest("GET", "/users", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		[{"id":1,"name":"a","age":2},{"id":2,"name":"b","age":3}]
	`)
	req = httptest.NewRequest("GET", "/users/new", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("POST", "/users?name=matt&age=10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":3,"name":"matt","age":10}
	`)
	req = httptest.NewRequest("GET", "/users/10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":10,"name":"d","age":5}
	`)
	req = httptest.NewRequest("GET", "/users/10/edit", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("PATCH", "/users/10?name=matt&age=10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"id":10,"name":"matt","age":10}
	`)
	req = httptest.NewRequest("PATCH", "/users/0", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"missing id"}
	`)
	req = httptest.NewRequest("DELETE", "/users/10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("DELETE", "/users/0", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"missing id"}
	`)
	req = httptest.NewRequest("GET", "/users/hello", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"Message":"hello"}
	`)
	// HTML
	req = httptest.NewRequest("GET", "/users", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		index
	`)
	req = httptest.NewRequest("GET", "/users/new", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		new
	`)
	req = httptest.NewRequest("POST", "/users?name=matt&age=10", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users
	`)
	req = httptest.NewRequest("POST", "/users?name=matt", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users

		missing age
	`)
	req = httptest.NewRequest("GET", "/users/10", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		show
	`)
	req = httptest.NewRequest("GET", "/users/10/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		edit
	`)
	req = httptest.NewRequest("GET", "/users/0/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		error page: id not found
	`)
	req = httptest.NewRequest("PATCH", "/users/10?name=matt&age=10", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users/10
	`)
	req = httptest.NewRequest("PATCH", "/users/0", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users/0

		missing id
	`)
	req = httptest.NewRequest("DELETE", "/users/10", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users/10
	`)
	req = httptest.NewRequest("DELETE", "/users/0", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /users/0

		missing id
	`)
}

type commentsResource struct{}

type commentsIndex struct {
	PostID int `json:"post_id"`
}

type commentResource struct {
	PostID    int       `json:"post_id"`
	ID        int       `json:"id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *commentsResource) Index(in *commentsIndex) ([]*commentResource, error) {
	t1 := time.Date(2023, 9, 4, 12, 8, 0, 0, time.UTC)
	t2 := time.Date(2023, 9, 5, 12, 8, 0, 0, time.UTC)
	return []*commentResource{{in.PostID, 1, "a", t1}, {in.PostID, 2, "b", t2}}, nil
}

func (c *commentsResource) New() {}

type commentsCreate struct {
	PostID  int    `json:"post_id"`
	Comment string `json:"comment"`
}

func (c *commentsResource) Create(in *commentsCreate) (*commentResource, error) {
	if in.PostID == 0 {
		return nil, errors.New("missing post_id")
	} else if in.Comment == "" {
		return nil, errors.New("missing comment")
	}
	t1 := time.Date(2023, 9, 4, 12, 1, 0, 0, time.UTC)
	return &commentResource{in.PostID, 3, in.Comment, t1}, nil
}

type commentsShow struct {
	PostID int `json:"post_id"`
	ID     int `json:"id"`
}

func (c *commentsResource) Show(in *commentsShow) (*commentResource, error) {
	if in.PostID == 0 {
		return nil, fmt.Errorf("post id not found")
	} else if in.ID == 0 {
		return nil, fmt.Errorf("id not found")
	}
	t1 := time.Date(2023, 9, 4, 12, 8, 0, 0, time.UTC)
	return &commentResource{in.PostID, in.ID, "a", t1}, nil
}

type commentsEdit struct {
	PostID int `json:"post_id"`
	ID     int `json:"id"`
}

func (c *commentsResource) Edit(in *commentsEdit) error {
	if in.PostID == 0 {
		return fmt.Errorf("post id not found")
	} else if in.ID == 0 {
		return fmt.Errorf("id not found")
	}
	return nil
}

type commentsUpdate struct {
	PostID  int     `json:"post_id"`
	ID      int     `json:"id"`
	Comment *string `json:"comment"`
}

func (c *commentsResource) Update(in *commentsUpdate) (*commentResource, error) {
	if in.PostID == 0 {
		return nil, fmt.Errorf("post id not found")
	} else if in.ID == 0 {
		return nil, errors.New("missing id")
	}
	comment := new(commentResource)
	comment.PostID = in.PostID
	comment.ID = in.ID
	if in.Comment != nil {
		comment.Comment = *in.Comment
	} else {
		comment.Comment = "a"
	}
	return comment, nil
}

type commentsDelete struct {
	PostID int `json:"post_id"`
	ID     int `json:"id"`
}

func (c *commentsResource) Delete(in *commentsDelete) error {
	if in.PostID == 0 {
		return fmt.Errorf("post id not found")
	} else if in.ID == 0 {
		return errors.New("missing id")
	}
	return nil
}

type previewOut struct {
	Comment string
}

func (c *commentsResource) Preview() (out previewOut) {
	out.Comment = "hello"
	return out
}

func TestDeepResource(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"posts/comments/index.gohtml":   &fstest.MapFile{Data: []byte(`{{range $post := . }}<div data-post-id="{{$post.PostID}}" data-id="{{$post.ID}}">{{$post.Comment}}</div>{{ end }}`)},
		"posts/comments/show.gohtml":    &fstest.MapFile{Data: []byte(`show: post_id={{.PostID}} id={{.ID}} comment={{.Comment}}`)},
		"posts/comments/new.gohtml":     &fstest.MapFile{Data: []byte(`new`)},
		"posts/comments/edit.gohtml":    &fstest.MapFile{Data: []byte(`edit`)},
		"posts/comments/preview.gohtml": &fstest.MapFile{Data: []byte(`preview`)},
		"posts/comments/error.gohtml":   &fstest.MapFile{Data: []byte(`error page: {{ .Error }}`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/posts/{post_id}/comments", &commentsResource{}))
	req := httptest.NewRequest("GET", "/posts/1/comments", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		[{"post_id":1,"id":1,"comment":"a","created_at":"2023-09-04T12:08:00Z"},{"post_id":1,"id":2,"comment":"b","created_at":"2023-09-05T12:08:00Z"}]
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/new", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("POST", "/posts/1/comments?comment=first", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":1,"id":3,"comment":"first","created_at":"2023-09-04T12:01:00Z"}
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":1,"id":10,"comment":"a","created_at":"2023-09-04T12:08:00Z"}
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/10/edit", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("GET", "/posts/0/comments/10/edit", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"post id not found"}
	`)
	req = httptest.NewRequest("PATCH", "/posts/1/comments/10?name=matt&age=10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":1,"id":10,"comment":"a","created_at":"0001-01-01T00:00:00Z"}
	`)
	req = httptest.NewRequest("PATCH", "/posts/1/comments/10?comment=b", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"post_id":1,"id":10,"comment":"b","created_at":"0001-01-01T00:00:00Z"}
	`)
	req = httptest.NewRequest("PATCH", "/posts/1/comments/0", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"missing id"}
	`)
	req = httptest.NewRequest("DELETE", "/posts/1/comments/10", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
	req = httptest.NewRequest("DELETE", "/posts/1/comments/0", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"missing id"}
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/preview", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"Comment":"hello"}
	`)
	// HTML
	req = httptest.NewRequest("GET", "/posts/1/comments", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		<div data-post-id="1" data-id="1">a</div><div data-post-id="1" data-id="2">b</div>
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/new", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		new
	`)
	req = httptest.NewRequest("POST", "/posts/1/comments?name=matt&age=10", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/1/comments

		missing comment
	`)
	req = httptest.NewRequest("POST", "/posts/1/comments?name=matt", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/1/comments

		missing comment
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/10", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		show: post_id=1 id=10 comment=a
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/10/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		edit
	`)
	req = httptest.NewRequest("GET", "/posts/1/comments/0/edit", nil)
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		error page: id not found
	`)
	req = httptest.NewRequest("PATCH", "/posts/1/comments/10?name=matt&age=10", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/1/comments/10
	`)
	req = httptest.NewRequest("PATCH", "/posts/1/comments/0", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/1/comments/0

		missing id
	`)
	req = httptest.NewRequest("DELETE", "/posts/1/comments/10", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/1/comments/10
	`)
	req = httptest.NewRequest("DELETE", "/posts/1/comments/0", nil)
	equal(t, controller, req, `
		HTTP/1.1 303 See Other
		Connection: close
		Location: /posts/1/comments/0

		missing id
	`)
}

type ambiguousResource struct {
}

func (a *ambiguousResource) Users() {}

type ambiguousUsersResource struct{}

func (*ambiguousUsersResource) Index() {
}

func TestAmbiguousAction(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"users/index.gohtml": &fstest.MapFile{Data: []byte(`index`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	router := controller.New(viewer, mux.New())
	is.NoErr(router.Register("/", &ambiguousResource{}))
	err := router.Register("/users", &ambiguousUsersResource{})
	is.True(errors.Is(err, mux.ErrDuplicate))
}

type handlerResource struct {
}

func (h *handlerResource) Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("index"))
}

func TestHandler(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{}
	viewer := view.New(fsys, map[string]view.Renderer{})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &handlerResource{}))
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		index
	`)
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		index
	`)
}

type escapeController struct{}

type escapeResource struct {
	HTML string
}

func (e *escapeController) Index() escapeResource {
	return escapeResource{HTML: `<b>hello<script type="text/javascript">alert('xss!')</script></b>`}
}

func TestEscapeProps(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"index.gohtml": &fstest.MapFile{Data: []byte(`{{ .HTML }}`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &escapeController{}))
	req := httptest.NewRequest("GET", "/", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		&lt;b&gt;hello&lt;script type=&#34;text/javascript&#34;&gt;alert(&#39;xss!&#39;)&lt;/script&gt;&lt;/b&gt;
	`)
}

type layoutFrameController struct {
}

type layoutFrameIndex struct {
	Message string
}

func (c *layoutFrameController) Index() layoutFrameIndex {
	return layoutFrameIndex{Message: "hello"}
}

func TestLayoutFrame(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"layout.gohtml": &fstest.MapFile{Data: []byte(`<html><body>{{ $.Slot }}</body></html>`)},
		"frame.gohtml":  &fstest.MapFile{Data: []byte(`<main>{{ $.Slot }}</main>`)},
		"index.gohtml":  &fstest.MapFile{Data: []byte(`<p>{{ .Message }}</p>`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &layoutFrameController{}))
	// HTML
	req := httptest.NewRequest("GET", "/", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		<html><body><main><p>hello</p></main></body></html>
	`)
	// JSON
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"Message":"hello"}
	`)
}

type layoutFrameErrorController struct {
}

func (c *layoutFrameErrorController) Index() error {
	return errors.New("unable to load index")
}

func TestLayoutFrameError(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"layout.gohtml": &fstest.MapFile{Data: []byte(`<html><body>{{ $.Slot }}</body></html>`)},
		"frame.gohtml":  &fstest.MapFile{Data: []byte(`<main>{{ $.Slot }}</main>`)},
		"index.gohtml":  &fstest.MapFile{Data: []byte(`<p>{{ .Message }}</p>`)},
		"error.gohtml":  &fstest.MapFile{Data: []byte(`<p>{{ .Error }}</p>`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &layoutFrameErrorController{}))
	// HTML
	req := httptest.NewRequest("GET", "/", nil)
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		<html><body><main><p>unable to load index</p></main></body></html>
	`)
	// JSON
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"unable to load index"}
	`)
}

type frameErrorController struct {
}

func (c *frameErrorController) Index() {
}

func TestFrameError(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"layout.gohtml": &fstest.MapFile{Data: []byte(`<html><body>{{ $.Slot }}</body></html>`)},
		"frame.gohtml":  &fstest.MapFile{Data: []byte(`<main>{{ blah }}</main>`)},
		"index.gohtml":  &fstest.MapFile{Data: []byte(`<p></p>`)},
		"error.gohtml":  &fstest.MapFile{Data: []byte(`<p>{{ .Error }}</p>`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &frameErrorController{}))
	// HTML
	req := httptest.NewRequest("GET", "/", nil)
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		<p>template: frame.gohtml:1: function &#34;blah&#34; not defined</p>
	`)
	// JSON
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 204 No Content
		Connection: close
	`)
}

type frameIndexErrorController struct {
}

func (c *frameIndexErrorController) Index() error {
	return errors.New("oh noz")
}

func TestFrameIndexError(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"layout.gohtml": &fstest.MapFile{Data: []byte(`<html><body>{{ $.Slot }}</body></html>`)},
		"frame.gohtml":  &fstest.MapFile{Data: []byte(`<main>{{ blah }}</main>`)},
		"index.gohtml":  &fstest.MapFile{Data: []byte(`<p>{{ .Message }}</p>`)},
		"error.gohtml":  &fstest.MapFile{Data: []byte(`<p>{{ .Error }}</p>`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &frameIndexErrorController{}))
	// HTML
	req := httptest.NewRequest("GET", "/", nil)
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: text/html

		<p>oh noz</p>
	`)
	// JSON
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 500 Internal Server Error
		Connection: close
		Content-Type: application/json

		{"error":"oh noz"}
	`)
}

// TODO:
// - Layout, Frame & Error handlers

type viewHandlers struct {
}

type viewHandler struct {
	Message string
}

func (v *viewHandlers) Index() viewHandler {
	return viewHandler{Message: "index"}
}

func (v *viewHandlers) Frame() viewHandler {
	return viewHandler{Message: "frame"}
}

func (v *viewHandlers) Layout() viewHandler {
	return viewHandler{Message: "layout"}
}

func TestViewHandlers(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"layout.gohtml": &fstest.MapFile{Data: []byte(`{{ .Message }} {{ $.Slot }}`)},
		"frame.gohtml":  &fstest.MapFile{Data: []byte(`{{ .Message }} {{ $.Slot }}`)},
		"index.gohtml":  &fstest.MapFile{Data: []byte(`{{ .Message }}`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	router := controller.New(viewer, mux.New())
	is.NoErr(router.Register("/", &viewHandlers{}))
	// HTML
	req := httptest.NewRequest("GET", "/", nil)
	equal(t, router, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		layout frame index
	`)
	// JSON
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, router, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"Message":"index"}
	`)
}

type nestedFrame struct{}

func (n *nestedFrame) Index() viewHandler {
	return viewHandler{Message: "index"}
}

func (n *nestedFrame) Layout() viewHandler {
	return viewHandler{Message: "layout"}
}

type nestedFramePosts struct{}

func (n *nestedFramePosts) Index() viewHandler {
	return viewHandler{Message: "posts/index"}
}

func (n *nestedFramePosts) Frame() viewHandler {
	return viewHandler{Message: "posts/frame"}
}

type nestedFrameSessions struct{}

func (n *nestedFrameSessions) Index() viewHandler {
	return viewHandler{Message: "sessions/index"}
}

func TestNestedFrame(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"layout.gohtml":         &fstest.MapFile{Data: []byte(`{{ .Message }} {{ $.Slot }}`)},
		"posts/frame.gohtml":    &fstest.MapFile{Data: []byte(`{{ .Message }} {{ $.Slot }}`)},
		"posts/index.gohtml":    &fstest.MapFile{Data: []byte(`{{ .Message }} {{ $.Slot }}`)},
		"sessions/index.gohtml": &fstest.MapFile{Data: []byte(`{{ .Message }}`)},
		"index.gohtml":          &fstest.MapFile{Data: []byte(`{{ .Message }}`)},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	router := controller.New(viewer, mux.New())
	is.NoErr(router.Register("/", &nestedFrame{}))
	is.NoErr(router.Register("/posts", &nestedFramePosts{}))
	is.NoErr(router.Register("/sessions", &nestedFrameSessions{}))
	// HTML
	// req := httptest.NewRequest("GET", "/", nil)
	// equal(t, router, req, `
	// 	HTTP/1.1 200 OK
	// 	Connection: close
	// 	Content-Type: text/html

	// 	layout index
	// `)
	// // JSON
	// req = httptest.NewRequest("GET", "/", nil)
	// req.Header.Set("Accept", "application/json")
	// equal(t, router, req, `
	// 	HTTP/1.1 200 OK
	// 	Connection: close
	// 	Content-Type: application/json

	// 	{"Message":"index"}
	// `)
	// // HTML
	// req = httptest.NewRequest("GET", "/posts", nil)
	// equal(t, router, req, `
	// 	HTTP/1.1 200 OK
	// 	Connection: close
	// 	Content-Type: text/html

	// 	layout posts/frame posts/index
	// `)
	// // JSON
	// req = httptest.NewRequest("GET", "/posts", nil)
	// req.Header.Set("Accept", "application/json")
	// equal(t, router, req, `
	// 	HTTP/1.1 200 OK
	// 	Connection: close
	// 	Content-Type: application/json

	// 	{"Message":"posts/index"}
	// `)
	// HTML
	req := httptest.NewRequest("GET", "/sessions", nil)
	equal(t, router, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		layout sessions/index
	`)
	// JSON
	req = httptest.NewRequest("GET", "/sessions", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, router, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"Message":"sessions/index"}
	`)
}

// TODO: test underscores (e.g. ExecuteQuery)

type underscores struct {
}

func (u *underscores) QueryBuilder() struct{ Message string } {
	return struct{ Message string }{Message: "Hello"}
}

func TestUnderscores(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"query_builder.gohtml": &fstest.MapFile{Data: []byte("{{.Message}}")},
	}
	viewer := view.New(fsys, map[string]view.Renderer{
		".gohtml": gohtml.New(),
	})
	controller := controller.New(viewer, mux.New())
	is.NoErr(controller.Register("/", &underscores{}))
	req := httptest.NewRequest("GET", "/query_builder", nil)
	req.Header.Set("Accept", "application/json")
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/json

		{"Message":"Hello"}
	`)
	req = httptest.NewRequest("GET", "/query_builder", nil)
	equal(t, controller, req, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/html

		Hello
	`)
}
