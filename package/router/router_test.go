package router_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/router"
)

type test struct {
	routes   []*route
	requests []*request
}

type route struct {
	method string
	route  string
	err    string
}

type request struct {
	method string
	path   string

	// response
	status   int
	location string
	body     string
}

// Handler returns the raw query
func handler(route string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.RawQuery))
	})
}

func ok(t testing.TB, test *test) {
	is := is.New(t)
	router := router.New()
	for _, route := range test.routes {
		if route.method == "" {
			route.method = http.MethodGet
		}
		var err error
		switch route.method {
		case http.MethodGet:
			err = router.Get(route.route, handler(route.route))
		case http.MethodPost:
			err = router.Post(route.route, handler(route.route))
		case http.MethodPatch:
			err = router.Patch(route.route, handler(route.route))
		case http.MethodPut:
			err = router.Put(route.route, handler(route.route))
		case http.MethodDelete:
			err = router.Delete(route.route, handler(route.route))
		default:
			err = router.Add(route.method, route.route, handler(route.route))
		}
		if err != nil {
			is.Equal(route.err, err.Error())
			continue
		} else if route.err != "" {
			is.Equal(route.err, err)
			continue
		}
	}
	for _, request := range test.requests {
		// Test the handler
		req := httptest.NewRequest(request.method, request.path, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		res := rec.Result()
		is.Equal(request.status, res.StatusCode)
		if request.location != "" {
			url, err := res.Location()
			is.NoErr(err)
			fmt.Println("location", url.Path)
			is.Equal(request.location, url.Path)
		}
		body, err := ioutil.ReadAll(res.Body)
		is.NoErr(err)
		is.Equal(request.body, string(body))
	}
}

func TestRouter(t *testing.T) {
	ok(t, &test{
		routes: []*route{
			{method: "GET", route: "/"},
			// users
			// TODO: add :format? to every route
			{method: "GET", route: `/users`},
			{method: "GET", route: `/users/new`},
			{method: "POST", route: `/users`},
			{method: "GET", route: `/users/:id.:format?`},
			{method: "GET", route: `/users/:id/edit`},
			{method: "PATCH", route: `/users/:id.:format?`},
			{method: "DELETE", route: `/users/:id.:format?`},
			// comments
			{method: "GET", route: `/posts/:post_id/comments`},
			{method: "GET", route: `/posts/:post_id/comments/new`},
			{method: "POST", route: `/posts/:post_id/comments`},
			{method: "GET", route: `/posts/:post_id/comments/:id.:format?`},
			{method: "GET", route: `/posts/:post_id/comments/:id/edit`},
			{method: "PATCH", route: `/posts/:post_id/comments/:id.:format?`},
			{method: "DELETE", route: `/posts/:post_id/comments/:id.:format?`},
		},
		requests: []*request{
			{method: "GET", path: "/", status: 200},
			// users
			{method: "GET", path: `/users`, status: 200},
			{method: "GET", path: `/users/new`, status: 200},
			{method: "POST", path: `/users`, status: 200},
			{method: "GET", path: `/users/10`, status: 200, body: "id=10"},
			{method: "GET", path: `/users/10.json`, status: 200, body: "format=json&id=10"},
			{method: "GET", path: `/users/10.rss`, status: 200, body: "format=rss&id=10"},
			{method: "GET", path: `/users/10.html`, status: 200, body: "format=html&id=10"},
			{method: "GET", path: `/users/10/edit`, status: 200, body: "id=10"},
			{method: "PATCH", path: `/users/10`, status: 200, body: "id=10"},
			{method: "PATCH", path: `/users/10.json`, status: 200, body: "format=json&id=10"},
			{method: "PATCH", path: `/users/10.rss`, status: 200, body: "format=rss&id=10"},
			{method: "PATCH", path: `/users/10.html`, status: 200, body: "format=html&id=10"},
			{method: "DELETE", path: `/users/10`, status: 200, body: "id=10"},
			{method: "DELETE", path: `/users/10.json`, status: 200, body: "format=json&id=10"},
			{method: "DELETE", path: `/users/10.rss`, status: 200, body: "format=rss&id=10"},
			{method: "DELETE", path: `/users/10.html`, status: 200, body: "format=html&id=10"},
			// comments
			{method: "GET", path: `/posts/1/comments`, status: 200, body: "post_id=1"},
			{method: "GET", path: `/posts/1/comments/new`, status: 200, body: "post_id=1"},
			{method: "POST", path: `/posts/1/comments`, status: 200, body: "post_id=1"},
			{method: "GET", path: `/posts/1/comments/2`, status: 200, body: "id=2&post_id=1"},
			{method: "GET", path: `/posts/1/comments/2.json`, status: 200, body: "format=json&id=2&post_id=1"},
			{method: "GET", path: `/posts/1/comments/2.rss`, status: 200, body: "format=rss&id=2&post_id=1"},
			{method: "GET", path: `/posts/1/comments/2.html`, status: 200, body: "format=html&id=2&post_id=1"},
			{method: "GET", path: `/posts/1/comments/2/edit`, status: 200, body: "id=2&post_id=1"},
			{method: "PATCH", path: `/posts/1/comments/2`, status: 200, body: "id=2&post_id=1"},
			{method: "PATCH", path: `/posts/1/comments/2.json`, status: 200, body: "format=json&id=2&post_id=1"},
			{method: "PATCH", path: `/posts/1/comments/2.rss`, status: 200, body: "format=rss&id=2&post_id=1"},
			{method: "PATCH", path: `/posts/1/comments/2.html`, status: 200, body: "format=html&id=2&post_id=1"},
			{method: "DELETE", path: `/posts/1/comments/2`, status: 200, body: "id=2&post_id=1"},
			{method: "DELETE", path: `/posts/1/comments/2.json`, status: 200, body: "format=json&id=2&post_id=1"},
			{method: "DELETE", path: `/posts/1/comments/2.rss`, status: 200, body: "format=rss&id=2&post_id=1"},
			{method: "DELETE", path: `/posts/1/comments/2.html`, status: 200, body: "format=html&id=2&post_id=1"},
		},
	})
}

func TestQueryPriority(t *testing.T) {
	ok(t, &test{
		routes: []*route{
			{method: "GET", route: "/"},
			{method: "GET", route: "/users/:id.:format?"},
			{method: "GET", route: "/posts/:post_id/comments/:id.:format?"},
		},
		requests: []*request{
			{method: "GET", path: "/?id=10", status: 200, body: "id=10"},
			{method: "GET", path: `/users/10?id=20&format=bin&other=true`, status: 200, body: "format=bin&id=10&other=true"},
			{method: "GET", path: `/users/10.json?id=20&format=bin&other=true`, status: 200, body: "format=json&id=10&other=true"},
			{method: "GET", path: `/posts/1/comments/2?post_id=10&id=20&other=true`, status: 200, body: "id=2&other=true&post_id=1"},
			{method: "GET", path: `/posts/1/comments/2.json?format=bin&post_id=10&id=20&other=true`, status: 200, body: "format=json&id=2&other=true&post_id=1"},
		},
	})
}

func TestTrailingSlash(t *testing.T) {
	ok(t, &test{
		routes: []*route{
			{method: "GET", route: "/"},
			{method: "GET", route: "/hi/", err: `route "/hi/": remove the slash "/" at the end`},
			{method: "GET", route: "/hi"},
		},
		requests: []*request{
			{method: "GET", path: "/", status: 200},
			{method: "GET", path: "/hi/", status: 308, location: "/hi", body: "<a href=\"/hi\">Permanent Redirect</a>.\n\n"},
			{method: "GET", path: "/hi///", status: 308, location: "/hi", body: "<a href=\"/hi\">Permanent Redirect</a>.\n\n"},
		},
	})
}
func TestInsensitive(t *testing.T) {
	ok(t, &test{
		routes: []*route{
			{method: "GET", route: "/HI", err: `route "/HI": uppercase letters are not allowed "H"`},
			{method: "GET", route: "/hi"},
		},
		requests: []*request{
			{method: "GET", path: "/HI", status: 308, location: "/hi", body: "<a href=\"/hi\">Permanent Redirect</a>.\n\n"},
			{method: "GET", path: "/Hi", status: 308, location: "/hi", body: "<a href=\"/hi\">Permanent Redirect</a>.\n\n"},
			{method: "GET", path: "/hI", status: 308, location: "/hi", body: "<a href=\"/hi\">Permanent Redirect</a>.\n\n"},
			{method: "GET", path: "/HI///", status: 308, location: "/hi", body: "<a href=\"/hi\">Permanent Redirect</a>.\n\n"},
		},
	})
}

func TestPut(t *testing.T) {
	is := is.New(t)
	router := router.New()
	is.NoErr(router.Put("/:id", handler("/:id")))
	req := httptest.NewRequest(http.MethodPut, "/10", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(200, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal("id=10", string(body))
}

func TestAdd(t *testing.T) {
	is := is.New(t)
	router := router.New()
	is.NoErr(router.Add(http.MethodHead, "/:id", handler("/:id")))
	req := httptest.NewRequest(http.MethodHead, "/10", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	res := rec.Result()
	is.Equal(200, res.StatusCode)
	body, err := ioutil.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal("id=10", string(body))
}
