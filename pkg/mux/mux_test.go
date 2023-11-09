package mux_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"

	"github.com/livebud/bud/pkg/middleware"
	"github.com/livebud/bud/pkg/mux"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
)

// Handler returns the raw query
func handler(route string) *handlerFunc {
	return &handlerFunc{
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(route + " " + r.URL.RawQuery))
		},
		http.Header{},
	}
}

type handlerFunc struct {
	fn      func(w http.ResponseWriter, r *http.Request)
	headers http.Header
}

func (h *handlerFunc) Set(name, value string) *handlerFunc {
	h.headers.Set(name, value)
	return h
}

func (h *handlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for key := range h.headers {
		w.Header().Set(key, h.headers.Get(key))
	}
	h.fn(w, r)
}

func requestEqual(t testing.TB, router http.Handler, request string, expect string) {
	t.Helper()
	parts := strings.SplitN(request, " ", 2)
	if len(parts) != 2 {
		t.Fatalf("invalid request: %s", request)
	}
	u, err := url.Parse(parts[1])
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(parts[0], u.Path, nil)
	req.URL.RawQuery = u.RawQuery
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	actual, err := httputil.DumpResponse(rec.Result(), true)
	if err != nil {
		if err.Error() == expect {
			return
		}
		t.Fatal(err)
	}
	diff.TestHTTP(t, expect, string(actual))
}

func TestSanity(t *testing.T) {
	router := mux.New()
	router.Get("/", handler("GET /"))
	requestEqual(t, router, "GET /", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /
	`)
	requestEqual(t, router, "GET /hi", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	router.Get("/{name}", handler("GET /{name}"))
	requestEqual(t, router, "GET /anki", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /{name} name=anki
	`)
}

func TestSample(t *testing.T) {
	router := mux.New()
	router.Get("/", handler("GET /"))
	router.Get("/users/{id}", handler("GET /users/{id}"))
	router.Post("/users/{id}.{format}", handler("POST /users/{id}.{format}"))
	router.Get("/posts/{post_id}/comments/{id}", handler("GET /posts/{post_id}/comments/{id}"))
	router.Get("/fly/{from}-{to}", handler("GET /fly/{from}-{to}"))
	router.Get("/v{major|[0-9]+}.{minor|[0-9]+}", handler("GET /v{major|[0-9]+}.{minor|[0-9]+}"))
	router.Get("/{owner}/{repo}/{branch}/{path*}", handler("GET /{owner}/{repo}/{branch}/{path*}"))
	requestEqual(t, router, "GET /", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /
	`)
	requestEqual(t, router, "GET /users/1", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /users/{id} id=1
	`)
	requestEqual(t, router, "POST /users/1.json", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		POST /users/{id}.{format} format=json&id=1
	`)
	requestEqual(t, router, "GET /posts/1/comments/1", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /posts/{post_id}/comments/{id} id=1&post_id=1
	`)
	requestEqual(t, router, "GET /fly/sfo-lax", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /fly/{from}-{to} from=sfo&to=lax
	`)
	requestEqual(t, router, "GET /v1.0", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /v{major|[0-9]+}.{minor|[0-9]+} major=1&minor=0
	`)
	requestEqual(t, router, "GET /1.0", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	requestEqual(t, router, "GET /v1.a", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	requestEqual(t, router, "GET /livebud/mux/main/path/to/file.go", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /{owner}/{repo}/{branch}/{path*} branch=main&owner=livebud&path=path%2Fto%2Ffile.go&repo=mux
	`)
}

func TestREST(t *testing.T) {
	router := mux.New()
	router.Get("/", handler("GET /"))
	router.Get("/users", handler("GET /users"))
	router.Get("/users/new", handler("GET /users/new"))
	router.Post("/users", handler("POST /users"))
	router.Get("/users/{id}.{format?}", handler("GET /users/{id}.{format?}"))
	router.Get("/users/{id}/edit", handler("GET /users/{id}/edit"))
	router.Patch("/users/{id}.{format?}", handler("PATCH /users/{id}.{format?}"))
	router.Put("/users/{id}.{format?}", handler("PUT /users/{id}.{format?}"))
	router.Delete("/users/{id}.{format?}", handler("DELETE /users/{id}.{format?}"))
	router.Get("/posts/{post_id}/comments", handler("GET /posts/{post_id}/comments"))
	router.Get("/posts/{postid}/comments/new", handler("GET /posts/{postid}/comments/new"))
	router.Post("/posts/{post_id}/comments", handler("POST /posts/{post_id}/comments"))
	router.Get("/posts/{post_id}/comments/{id}.{format?}", handler("GET /posts/{post_id}/comments/{id}.{format?}"))
	router.Get("/posts/{post_id}/comments/{id}/edit", handler("GET /posts/{post_id}/comments/{id}/edit"))
	router.Patch("/posts/{post_id}/comments/{id}.{format?}", handler("PATCH /posts/{post_id}/comments/{id}.{format?}"))
	router.Put("/posts/{post_id}/comments/{id}.{format?}", handler("PUT /posts/{post_id}/comments/{id}.{format?}"))
	router.Delete("/posts/{post_id}/comments/{id}.{format?}", handler("DELETE /posts/{post_id}/comments/{id}.{format?}"))

	// requests
	requestEqual(t, router, "GET /", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /
	`)
	requestEqual(t, router, "GET /users", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /users
	`)
	requestEqual(t, router, "GET /users/new", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /users/new
	`)
	requestEqual(t, router, "POST /users", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		POST /users
	`)
	requestEqual(t, router, "GET /users/10", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	requestEqual(t, router, "GET /users/10.json", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /users/{id}.{format?} format=json&id=10
	`)
	requestEqual(t, router, "GET /users/10.rss", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /users/{id}.{format?} format=rss&id=10
	`)
	requestEqual(t, router, "GET /users/10.html", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /users/{id}.{format?} format=html&id=10
	`)
	requestEqual(t, router, "GET /users/10/edit", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /users/{id}/edit id=10
	`)
	requestEqual(t, router, "PATCH /users/10", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	requestEqual(t, router, "PATCH /users/10.json", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PATCH /users/{id}.{format?} format=json&id=10
	`)
	requestEqual(t, router, "PATCH /users/10.rss", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PATCH /users/{id}.{format?} format=rss&id=10
	`)
	requestEqual(t, router, "PATCH /users/10.html", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PATCH /users/{id}.{format?} format=html&id=10
	`)
	requestEqual(t, router, "PUT /users/10", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	requestEqual(t, router, "PUT /users/10.json", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PUT /users/{id}.{format?} format=json&id=10
	`)
	requestEqual(t, router, "PUT /users/10.rss", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PUT /users/{id}.{format?} format=rss&id=10
	`)
	requestEqual(t, router, "PUT /users/10.html", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PUT /users/{id}.{format?} format=html&id=10
	`)
	requestEqual(t, router, "DELETE /users/10", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	requestEqual(t, router, "DELETE /users/10.json", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		DELETE /users/{id}.{format?} format=json&id=10
	`)
	requestEqual(t, router, "DELETE /users/10.rss", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		DELETE /users/{id}.{format?} format=rss&id=10
	`)
	requestEqual(t, router, "DELETE /users/10.html", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		DELETE /users/{id}.{format?} format=html&id=10
	`)
	requestEqual(t, router, "GET /posts/1/comments", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /posts/{post_id}/comments post_id=1
	`)
	requestEqual(t, router, "GET /posts/1/comments/new", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /posts/{postid}/comments/new postid=1
	`)
	requestEqual(t, router, "POST /posts/1/comments", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		POST /posts/{post_id}/comments post_id=1
	`)
	requestEqual(t, router, "GET /posts/1/comments/2", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	requestEqual(t, router, "GET /posts/1/comments/2.json", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /posts/{post_id}/comments/{id}.{format?} format=json&id=2&post_id=1
	`)
	requestEqual(t, router, "GET /posts/1/comments/2.rss", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /posts/{post_id}/comments/{id}.{format?} format=rss&id=2&post_id=1
	`)
	requestEqual(t, router, "GET /posts/1/comments/2.html", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /posts/{post_id}/comments/{id}.{format?} format=html&id=2&post_id=1
	`)
	requestEqual(t, router, "GET /posts/1/comments/2/edit", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /posts/{post_id}/comments/{id}/edit id=2&post_id=1
	`)
	requestEqual(t, router, "PATCH /posts/1/comments/2", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	requestEqual(t, router, "PATCH /posts/1/comments/2.json", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PATCH /posts/{post_id}/comments/{id}.{format?} format=json&id=2&post_id=1
	`)
	requestEqual(t, router, "PATCH /posts/1/comments/2.rss", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PATCH /posts/{post_id}/comments/{id}.{format?} format=rss&id=2&post_id=1
	`)
	requestEqual(t, router, "PATCH /posts/1/comments/2.html", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PATCH /posts/{post_id}/comments/{id}.{format?} format=html&id=2&post_id=1
	`)
	requestEqual(t, router, "PUT /posts/1/comments/2", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	requestEqual(t, router, "PUT /posts/1/comments/2.json", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PUT /posts/{post_id}/comments/{id}.{format?} format=json&id=2&post_id=1
	`)
	requestEqual(t, router, "PUT /posts/1/comments/2.rss", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PUT /posts/{post_id}/comments/{id}.{format?} format=rss&id=2&post_id=1
	`)
	requestEqual(t, router, "PUT /posts/1/comments/2.html", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		PUT /posts/{post_id}/comments/{id}.{format?} format=html&id=2&post_id=1
	`)
	requestEqual(t, router, "DELETE /posts/1/comments/2", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	requestEqual(t, router, "DELETE /posts/1/comments/2.json", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		DELETE /posts/{post_id}/comments/{id}.{format?} format=json&id=2&post_id=1
	`)
	requestEqual(t, router, "DELETE /posts/1/comments/2.rss", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		DELETE /posts/{post_id}/comments/{id}.{format?} format=rss&id=2&post_id=1
	`)
	requestEqual(t, router, "DELETE /posts/1/comments/2.html", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		DELETE /posts/{post_id}/comments/{id}.{format?} format=html&id=2&post_id=1
	`)
}

func TestSlotPriority(t *testing.T) {
	router := mux.New()
	router.Get("/", handler("GET /"))
	router.Get("/users/{id}", handler("GET /users/{id}"))
	router.Get("/users/{id}.{format?}", handler("GET /users/{id}.{format?}"))
	router.Get("/posts/{post_id}/comments/{id}.{format?}", handler("GET /posts/{post_id}/comments/{id}.{format?}"))

	requestEqual(t, router, "GET /?id=10", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET / id=10
	`)
	requestEqual(t, router, `GET /users/10?id=20&format=bin&other=true`, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /users/{id} format=bin&id=10&other=true
	`)
	requestEqual(t, router, `GET /users/10.json?id=20&format=bin&other=true`, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /users/{id}.{format?} format=json&id=10&other=true
	`)
	requestEqual(t, router, `GET /posts/10/comments/20.json?id=30&post_id=30&format=bin&other=true`, `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /posts/{post_id}/comments/{id}.{format?} format=json&id=20&other=true&post_id=10
	`)
}

func TestTrailingSlash(t *testing.T) {
	router := mux.New()
	router.Get("/", handler("GET /"))
	router.Get("/hi/", handler("GET /hi/"))
	router.Get("/hi", handler("GET /hi"))

	requestEqual(t, router, "GET /", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /
	`)
	requestEqual(t, router, "GET /hi/", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /hi/
	`)
	requestEqual(t, router, "GET /hi///", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /hi/
	`)
}

func TestInsensitive(t *testing.T) {
	router := mux.New()
	router.Get("/HI", handler("GET /HI"))
	router.Get("/hi", handler("GET /hi"))
	router.Get("/Hi", handler("GET /Hi"))
	router.Get("/hI", handler("GET /hI"))
	router.Get("/HI/", handler("GET /HI/"))
	router.Get("/hi/", handler("GET /hi/"))
	router.Get("/hI/", handler("GET /hI/"))
	router.Get("/Hi/", handler("GET /Hi/"))

	requestEqual(t, router, "GET /hi", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /hi
	`)
	requestEqual(t, router, "GET /HI", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /hi
	`)
	requestEqual(t, router, "GET /Hi", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /hi
	`)
	requestEqual(t, router, "GET /hi/", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /hi
	`)
	requestEqual(t, router, "GET /HI/", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /hi
	`)
	requestEqual(t, router, "GET /Hi/", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /hi
	`)
	requestEqual(t, router, "GET /hI/", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /hi
	`)
	requestEqual(t, router, "GET /HI////", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /hi
	`)
}

func TestSet(t *testing.T) {
	router := mux.New()
	router.Set(http.MethodHead, "/{id}", handler("HEAD /{id}"))
	requestEqual(t, router, "GET /10", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-Content-Type-Options: nosniff

		404 page not found
	`)
	requestEqual(t, router, "HEAD /10", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		HEAD /{id} id=10
	`)
}

func TestResource(t *testing.T) {
	router := mux.New()
	router.Set(http.MethodGet, "/{id}/edit", handler("GET /{id}/edit"))
	router.Set(http.MethodGet, "/", handler("GET /"))
	router.Set(http.MethodGet, "/new", handler("GET /new"))
	router.Set(http.MethodGet, "/{id}", handler("GET /{id}"))
	requestEqual(t, router, "GET /", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /
	`)
	requestEqual(t, router, "GET /10", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /{id} id=10
	`)
	requestEqual(t, router, "GET /new", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /new
	`)
	requestEqual(t, router, "GET /10/edit", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /{id}/edit id=10
	`)
}

func TestFind(t *testing.T) {
	is := is.New(t)
	router := mux.New()
	h := handler("GET /{id}")
	err := router.Get("/{id}", h)
	is.NoErr(err)
	route, err := router.Find("GET", "/{id}")
	is.NoErr(err)
	is.Equal(route.Method, "GET")
	is.Equal(route.Route, "/{id}")
	is.Equal(route.Handler, h)
	route, err = router.Find("POST", "/{id}")
	is.True(errors.Is(err, mux.ErrNoMatch))
	is.Equal(route, nil)
}

func TestDuplicateError(t *testing.T) {
	is := is.New(t)
	router := mux.New()
	err := router.Get("/", handler("GET /"))
	is.NoErr(err)
	err = router.Get("/", handler("GET /"))
	is.True(errors.Is(err, mux.ErrDuplicate))
}

func TestList(t *testing.T) {
	is := is.New(t)
	router := mux.New()
	is.NoErr(router.Get("/", handler("GET /")))
	is.NoErr(router.Get("/users", handler("GET /users")))
	is.NoErr(router.Get("/users/new", handler("GET /users/new")))
	is.NoErr(router.Post("/users", handler("POST /users")))
	is.NoErr(router.Get("/users/{id}.{format?}", handler("GET /users/{id}.{format?}")))
	is.NoErr(router.Get("/users/{id}/edit", handler("GET /users/{id}/edit")))
	is.NoErr(router.Patch("/users/{id}.{format?}", handler("PATCH /users/{id}.{format?}")))
	is.NoErr(router.Put("/users/{id}.{format?}", handler("PUT /users/{id}.{format?}")))
	is.NoErr(router.Delete("/users/{id}.{format?}", handler("DELETE /users/{id}.{format?}")))
	is.NoErr(router.Get("/posts/{post_id}/comments", handler("GET /posts/{post_id}/comments")))
	is.NoErr(router.Get("/posts/{postid}/comments/new", handler("GET /posts/{postid}/comments/new")))
	is.NoErr(router.Post("/posts/{post_id}/comments", handler("POST /posts/{post_id}/comments")))
	is.NoErr(router.Get("/posts/{post_id}/comments/{id}.{format?}", handler("GET /posts/{post_id}/comments/{id}.{format?}")))
	is.NoErr(router.Get("/posts/{post_id}/comments/{id}/edit", handler("GET /posts/{post_id}/comments/{id}/edit")))
	is.NoErr(router.Patch("/posts/{post_id}/comments/{id}.{format?}", handler("PATCH /posts/{post_id}/comments/{id}.{format?}")))
	is.NoErr(router.Put("/posts/{post_id}/comments/{id}.{format?}", handler("PUT /posts/{post_id}/comments/{id}.{format?}")))
	is.NoErr(router.Delete("/posts/{post_id}/comments/{id}.{format?}", handler("DELETE /posts/{post_id}/comments/{id}.{format?}")))
	routes := router.Routes()
	is.Equal(len(routes), 25)
	is.Equal(routes[0].String(), "GET /")
	is.Equal(routes[1].String(), "GET /posts/{post_id}/comments")
	is.Equal(routes[2].String(), "GET /posts/{post_id}/comments/{id}.")
	is.Equal(routes[3].String(), "GET /posts/{post_id}/comments/{id}.{format}")
	is.Equal(routes[4].String(), "GET /posts/{post_id}/comments/{id}/edit")
	is.Equal(routes[5].String(), "GET /posts/{postid}/comments/new")
	is.Equal(routes[6].String(), "GET /users")
	is.Equal(routes[7].String(), "GET /users/new")
	is.Equal(routes[8].String(), "GET /users/{id}.")
	is.Equal(routes[9].String(), "GET /users/{id}.{format}")
	is.Equal(routes[10].String(), "GET /users/{id}/edit")
	is.Equal(routes[11].String(), "POST /posts/{post_id}/comments")
	is.Equal(routes[12].String(), "POST /users")
	is.Equal(routes[13].String(), "PUT /posts/{post_id}/comments/{id}.")
	is.Equal(routes[14].String(), "PUT /posts/{post_id}/comments/{id}.{format}")
	is.Equal(routes[15].String(), "PUT /users/{id}.")
	is.Equal(routes[16].String(), "PUT /users/{id}.{format}")
	is.Equal(routes[17].String(), "PATCH /posts/{post_id}/comments/{id}.")
	is.Equal(routes[18].String(), "PATCH /posts/{post_id}/comments/{id}.{format}")
	is.Equal(routes[19].String(), "PATCH /users/{id}.")
	is.Equal(routes[20].String(), "PATCH /users/{id}.{format}")
	is.Equal(routes[21].String(), "DELETE /posts/{post_id}/comments/{id}.")
	is.Equal(routes[22].String(), "DELETE /posts/{post_id}/comments/{id}.{format}")
	is.Equal(routes[23].String(), "DELETE /users/{id}.")
	is.Equal(routes[24].String(), "DELETE /users/{id}.{format}")
}

func TestMissingRoot(t *testing.T) {
	router := mux.New()
	router.Get("/signup", handler("GET /signup"))
	router.Post("/users", handler("POST /users"))
	router.Get("/login", handler("GET /login"))
	router.Post("/sessions", handler("POST /sessions"))
	router.Delete("/sessions", handler("DELETE /sessions"))
	router.Get("/", handler("GET /"))

	requestEqual(t, router, "GET /", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8

		GET /
	`)
}

func TestMatch(t *testing.T) {
	is := is.New(t)
	router := mux.New()
	router.Get("/signup", handler("GET /signup"))
	router.Post("/users", handler("POST /users"))
	router.Get("/login", handler("GET /login"))
	router.Post("/sessions", handler("POST /sessions"))
	router.Delete("/sessions", handler("DELETE /sessions"))
	router.Get("/users/{user_id}/posts/{id}", handler("GET /users/{user_id}/posts/{id}"))
	router.Get("/", handler("GET /"))

	match, err := router.Match(http.MethodGet, "/")
	is.NoErr(err)
	is.Equal(match.Route.String(), "/")
	is.Equal(match.Path, "/")
	is.Equal(len(match.Slots), 0)

	match, err = router.Match(http.MethodPost, "/")
	is.True(errors.Is(err, mux.ErrNoMatch))
	is.Equal(match, nil)

	match, err = router.Match(http.MethodGet, "/ok")
	is.True(errors.Is(err, mux.ErrNoMatch))
	is.Equal(match, nil)

	match, err = router.Match(http.MethodPost, "/ok")
	is.True(errors.Is(err, mux.ErrNoMatch))
	is.Equal(match, nil)

	match, err = router.Match(http.MethodPost, "/users")
	is.NoErr(err)
	is.Equal(match.Route.String(), "/users")
	is.Equal(match.Path, "/users")
	is.Equal(len(match.Slots), 0)

	match, err = router.Match(http.MethodPost, "/users/")
	is.NoErr(err)
	is.Equal(match.Route.String(), "/users")
	is.Equal(match.Path, "/users")
	is.Equal(len(match.Slots), 0)

	match, err = router.Match(http.MethodGet, "/users/10/posts/20")
	is.NoErr(err)
	is.Equal(match.Route.String(), "/users/{user_id}/posts/{id}")
	is.Equal(match.Path, "/users/10/posts/20")
	is.Equal(len(match.Slots), 2)
	is.Equal(match.Slots[0].Key, "user_id")
	is.Equal(match.Slots[0].Value, "10")
	is.Equal(match.Slots[1].Key, "id")
	is.Equal(match.Slots[1].Value, "20")
}

func TestLayout(t *testing.T) {
	is := is.New(t)
	router := mux.New()
	root := handler("GET /")
	router.Get("/", root)
	signup := handler("GET /signup")
	router.Get("/signup", signup)
	createUsers := handler("POST /users")
	router.Post("/users", createUsers)
	listUsers := handler("GET /users")
	router.Get("/users", listUsers)
	layout := handler("/")
	router.Layout("/", layout)
	userLayout := handler("/users")
	router.Layout("/users", userLayout)

	// Match /
	match, err := router.Match(http.MethodGet, "/")
	is.NoErr(err)
	is.Equal(match.Method, http.MethodGet)
	is.Equal(match.Route.String(), "/")
	is.Equal(match.Path, "/")
	is.Equal(len(match.Slots), 0)
	is.Equal(match.Handler, root)
	is.True(match.Layout != nil)
	is.Equal(match.Layout.Route, "/")
	is.Equal(match.Layout.Handler, layout)
	is.True(match.Error == nil)

	// Find /
	route, err := router.Find(http.MethodGet, "/")
	is.NoErr(err)
	is.Equal(route.Method, http.MethodGet)
	is.Equal(route.Route, "/")
	is.Equal(route.Handler, root)
	is.True(route.Layout != nil)
	is.Equal(route.Layout.Route, "/")
	is.Equal(route.Layout.Handler, layout)
	is.True(route.Error == nil)

	// Match /signup
	match, err = router.Match(http.MethodGet, "/signup")
	is.NoErr(err)
	is.Equal(match.Method, http.MethodGet)
	is.Equal(match.Route.String(), "/signup")
	is.Equal(match.Path, "/signup")
	is.Equal(len(match.Slots), 0)
	is.Equal(match.Handler, signup)
	is.True(match.Layout != nil)
	is.Equal(match.Layout.Route, "/")
	is.Equal(match.Layout.Handler, layout)
	is.True(match.Error == nil)

	// Find /
	route, err = router.Find(http.MethodGet, "/signup")
	is.NoErr(err)
	is.Equal(route.Method, http.MethodGet)
	is.Equal(route.Route, "/signup")
	is.Equal(route.Handler, signup)
	is.True(route.Layout != nil)
	is.Equal(route.Layout.Route, "/")
	is.Equal(route.Layout.Handler, layout)
	is.True(route.Error == nil)

	// Match POST /users
	match, err = router.Match(http.MethodPost, "/users")
	is.NoErr(err)
	is.Equal(match.Method, http.MethodPost)
	is.Equal(match.Route.String(), "/users")
	is.Equal(match.Path, "/users")
	is.Equal(len(match.Slots), 0)
	is.Equal(match.Handler, createUsers)
	is.True(match.Layout == nil)
	is.True(match.Error == nil)

	// Find /
	route, err = router.Find(http.MethodPost, "/users")
	is.NoErr(err)
	is.Equal(route.Method, http.MethodPost)
	is.Equal(route.Route, "/users")
	is.Equal(route.Handler, createUsers)
	is.True(route.Layout == nil)
	is.True(route.Error == nil)

	// Match GET /users
	match, err = router.Match(http.MethodGet, "/users")
	is.NoErr(err)
	is.Equal(match.Method, http.MethodGet)
	is.Equal(match.Route.String(), "/users")
	is.Equal(match.Path, "/users")
	is.Equal(len(match.Slots), 0)
	is.Equal(match.Handler, listUsers)
	is.True(match.Layout != nil)
	is.Equal(match.Layout.Route, "/users")
	is.Equal(match.Layout.Handler, userLayout)
	is.True(match.Error == nil)

	// Find GET /users
	route, err = router.Find(http.MethodGet, "/users")
	is.NoErr(err)
	is.Equal(route.Method, http.MethodGet)
	is.Equal(route.Route, "/users")
	is.Equal(route.Handler, listUsers)
	is.True(route.Layout != nil)
	is.Equal(route.Layout.Route, "/users")
	is.Equal(route.Layout.Handler, userLayout)
	is.True(route.Error == nil)
}

func TestError(t *testing.T) {
	is := is.New(t)
	router := mux.New()
	root := handler("GET /")
	router.Get("/", root)
	signup := handler("GET /signup")
	router.Get("/signup", signup)
	createUsers := handler("POST /users")
	router.Post("/users", createUsers)
	listUsers := handler("GET /users")
	router.Get("/users", listUsers)
	errorHandler := handler("/")
	router.Error("/", errorHandler)
	userErrorHandler := handler("/users")
	router.Error("/users", userErrorHandler)

	// Match /
	match, err := router.Match(http.MethodGet, "/")
	is.NoErr(err)
	is.Equal(match.Method, http.MethodGet)
	is.Equal(match.Route.String(), "/")
	is.Equal(match.Path, "/")
	is.Equal(len(match.Slots), 0)
	is.Equal(match.Handler, root)
	is.True(match.Layout == nil)
	is.True(match.Error != nil)
	is.Equal(match.Error.Route, "/")
	is.Equal(match.Error.Handler, errorHandler)

	// Find /
	route, err := router.Find(http.MethodGet, "/")
	is.NoErr(err)
	is.Equal(route.Method, http.MethodGet)
	is.Equal(route.Route, "/")
	is.Equal(route.Handler, root)
	is.True(route.Layout == nil)
	is.True(route.Error != nil)
	is.Equal(route.Error.Route, "/")
	is.Equal(route.Error.Handler, errorHandler)

	// Match /signup
	match, err = router.Match(http.MethodGet, "/signup")
	is.NoErr(err)
	is.Equal(match.Method, http.MethodGet)
	is.Equal(match.Route.String(), "/signup")
	is.Equal(match.Path, "/signup")
	is.Equal(len(match.Slots), 0)
	is.Equal(match.Handler, signup)
	is.True(match.Layout == nil)
	is.True(match.Error != nil)
	is.Equal(match.Error.Route, "/")
	is.Equal(match.Error.Handler, errorHandler)

	// Find /signup
	route, err = router.Find(http.MethodGet, "/signup")
	is.NoErr(err)
	is.Equal(route.Method, http.MethodGet)
	is.Equal(route.Route, "/signup")
	is.Equal(route.Handler, signup)
	is.True(route.Layout == nil)
	is.True(route.Error != nil)
	is.Equal(route.Error.Route, "/")
	is.Equal(route.Error.Handler, errorHandler)

	// Match POST /users
	match, err = router.Match(http.MethodPost, "/users")
	is.NoErr(err)
	is.Equal(match.Method, http.MethodPost)
	is.Equal(match.Route.String(), "/users")
	is.Equal(match.Path, "/users")
	is.Equal(len(match.Slots), 0)
	is.Equal(match.Handler, createUsers)
	is.True(match.Layout == nil)
	is.True(match.Error == nil)

	// Find POST /users
	route, err = router.Find(http.MethodPost, "/users")
	is.NoErr(err)
	is.Equal(route.Method, http.MethodPost)
	is.Equal(route.Route, "/users")
	is.Equal(route.Handler, createUsers)
	is.True(route.Layout == nil)
	is.True(route.Error == nil)

	// Match GET /users
	match, err = router.Match(http.MethodGet, "/users")
	is.NoErr(err)
	is.Equal(match.Method, http.MethodGet)
	is.Equal(match.Route.String(), "/users")
	is.Equal(match.Path, "/users")
	is.Equal(len(match.Slots), 0)
	is.Equal(match.Handler, listUsers)
	is.True(match.Layout == nil)
	is.True(match.Error != nil)
	is.Equal(match.Error.Route, "/users")
	is.Equal(match.Error.Handler, userErrorHandler)

	route, err = router.Find(http.MethodGet, "/users")
	is.NoErr(err)
	is.Equal(route.Method, http.MethodGet)
	is.Equal(route.Route, "/users")
	is.Equal(route.Handler, listUsers)
	is.True(route.Layout == nil)
	is.True(route.Error != nil)
	is.Equal(route.Error.Route, "/users")
	is.Equal(route.Error.Handler, userErrorHandler)
}

// func TestLayoutRequest(t *testing.T) {
// 	is := is.New(t)
// 	router := mux.New()
// 	router.Layout("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		b := new(bytes.Buffer)
// 		b.WriteString("<layout>")
// 		slots := slot.Open(w, r)
// 		slot, err := slots.ReadString()
// 		is.NoErr(err)
// 		b.WriteString(slot)
// 		b.WriteString("</layout>")
// 		w.Write(b.Bytes())
// 	}))
// 	router.Get("/posts/{post_id}/comments", handler("GET /posts/{post_id}/comments"))
// 	router.Post("/posts/{post_id}/comments", handler("POST /posts/{post_id}/comments"))
// 	requestEqual(t, router, "GET /posts/10/comments", `
// 		HTTP/1.1 200 OK
// 		Connection: close
// 		Content-Type: text/plain; charset=utf-8

// 		<layout>GET /posts/{post_id}/comments post_id=10</layout>
// 	`)
// 	requestEqual(t, router, "POST /posts/10/comments", `
// 		HTTP/1.1 200 OK
// 		Connection: close
// 		Content-Type: text/plain; charset=utf-8

// 		POST /posts/{post_id}/comments post_id=10
// 	`)
// 	requestEqual(t, router, "GET /", `
// 		HTTP/1.1 404 Not Found
// 		Connection: close
// 		Content-Type: text/plain; charset=utf-8
// 		X-Content-Type-Options: nosniff

// 		404 page not found
// 	`)
// }

func TestLayoutNonHTMLRequest(t *testing.T) {
	is := is.New(t)
	router := mux.New()
	router.Layout("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inner, err := io.ReadAll(r.Body)
		is.NoErr(err)
		w.Write([]byte("<layout>"))
		w.Write(inner)
		w.Write([]byte("</layout>"))
	}))
	router.Get("/posts/{post_id}/comments.js", handler("console.log('hi')").Set("Content-Type", "application/javascript"))
	requestEqual(t, router, "GET /posts/10/comments.js", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: application/javascript

		console.log('hi') post_id=10
	`)
}

func TestMiddleware(t *testing.T) {
	router := mux.New()
	router.Use(middleware.Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-A", "A")
			next.ServeHTTP(w, r)
			// Note: Can't use a header here because we've already written
			w.Write([]byte("(after)"))
		})
	}))
	router.Get("/", handler("GET /"))
	requestEqual(t, router, "GET /", `
		HTTP/1.1 200 OK
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-A: A

		GET / (after)
	`)
}

func TestMiddlewareWrapsNonMatches(t *testing.T) {
	router := mux.New()
	router.Use(middleware.Func(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-A", "A")
			next.ServeHTTP(w, r)
		})
	}))
	router.Get("/", handler("GET /"))
	requestEqual(t, router, "POST /", `
		HTTP/1.1 404 Not Found
		Connection: close
		Content-Type: text/plain; charset=utf-8
		X-A: A
		X-Content-Type-Options: nosniff

		404 page not found
	`)
}

func TestPostBody(t *testing.T) {
	is := is.New(t)
	router := mux.New()
	router.Post("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		is.NoErr(err)
		w.Write(body)
	}))
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"jon"}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	is.Equal(rec.Code, http.StatusOK)
	is.Equal(rec.Body.String(), `{"name":"jon"}`)
}
