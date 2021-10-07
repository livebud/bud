package router_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-duo/bud/router"
	"github.com/matryer/is"
)

func TestGetRoot(t *testing.T) {
	is := is.New(t)
	r := router.New()
	req := httptest.NewRequest("GET", "/", nil)
	rec := httptest.NewRecorder()
	r.Get("/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("/"))
	}))
	r.ServeHTTP(rec, req)
	w := rec.Result()
	res, err := ioutil.ReadAll(w.Body)
	is.NoErr(err)
	is.Equal("/", string(res))
}

func TestGetAbout(t *testing.T) {
	is := is.New(t)
	r := router.New()
	req := httptest.NewRequest("GET", "/about", nil)
	rec := httptest.NewRecorder()
	r.Get("/about", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("/about"))
	}))
	r.ServeHTTP(rec, req)
	w := rec.Result()
	res, err := ioutil.ReadAll(w.Body)
	is.NoErr(err)
	is.Equal("/about", string(res))
}

func TestGetID(t *testing.T) {
	is := is.New(t)
	r := router.New()
	req := httptest.NewRequest("GET", "/users/10", nil)
	rec := httptest.NewRecorder()
	r.Get("/users/:id", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/" + r.URL.Query().Get("id")))
	}))
	r.ServeHTTP(rec, req)
	w := rec.Result()
	res, err := ioutil.ReadAll(w.Body)
	is.NoErr(err)
	is.Equal("/users/10", string(res))
}

func TestUse(t *testing.T) {
	is := is.New(t)
	r := router.New()
	req := httptest.NewRequest("POST", "/", nil)
	rec := httptest.NewRecorder()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	})
	r.Post("/", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("/"))
	}))
	r.ServeHTTP(rec, req)
	w := rec.Result()
	res, err := ioutil.ReadAll(w.Body)
	is.NoErr(err)
	is.Equal("/", string(res))
}

// func TestRegexp(t *testing.T) {
// is := is.New(t)
// 	// • keys can match anything but /
// 	// • regexps are for additional constraints (default key [^\/]+)
// 	// • ? is 0 or 1 matching path segments
// 	// • + is 1 or more matching path segments
// 	// • * is 0 or more matching path segments
// 	//
// 	// Given:
// 	//
// 	//   /users/:id.:format?
// 	//   /users/:id(v?\\d+\.\\d+\.\\d+).:format?
// 	//   /users/:id+.:format?
// 	//   /users/:id(v?\\d+\.\\d+\.\\d+).:format?
// 	//
// 	// How do we handle?
// 	//
// 	//   easy: /users/10
// 	//   easy: /users/10.json
// 	//   no: /users/10.
// 	//   pattern: /users/v10.2.1
// 	//   pattern: /users/v10.2.1.json
// 	//
// 	// • The trick is to be able to compile non-identifier suffixes into key matching negations (e.g. [])
// 	// • Any non-/ character before an optional key become a optional prefix to the path
// 	// • Patterns can be either named or unnamed. Unnamed parameters do not become keys
// 	// • No need for custom prefixes and suffixes (e.g. {...} in path-to-regexp)
// 	//
// 	// - [x] test in regexr how the other modifiers would compose together
// 	// - [ ] prototype in pegjs
// 	// - [ ] test against path-to-regexps test suite

// 	// /users/:id.:format?
// 	re := regexp.MustCompile(`^\/users\/([^\.\/]+)(?:\.([^\/]+))?$`)
// 	fmt.Println(re.FindAllStringSubmatch("/users/10.", -1))
// 	// /users/:id*.:format?
// 	re3 := regexp.MustCompile(`^\/users(?:\/?([^\.]*))(?:\.([^\/]+))?$`)
// 	fmt.Println(re3.FindAllStringSubmatch("/users/10.", -1))
// 	// /users/:id(v?\\d+\.\\d+\.\\d+).:format?
// 	re2 := regexp.MustCompile(`^\/users\/(v?\d+\.\d+\.\d+)(?:\.([^\/]+))?$`)
// 	fmt.Println(re2.FindAllStringSubmatch("/users/v10.2.1.", -1))
// 	// /users/:id(v?\\d+\.\\d+\.\\d+)*.:format?
// 	re4 := regexp.MustCompile(`^\/users(?:\/(v?\d+\.\d+\.\d+))*(?:\.([^\/]+))?$`)
// 	fmt.Println(re4.FindAllStringSubmatch("/users/v10.2.1.", -1))
// }

// // func TestMatch(t *testing.T) {
// is := is.New(t)
// // 	r := routes.New()
// // 	req := httptest.NewRequest("GET", "/coelho/alchemist", nil)
// // 	rec := httptest.NewRecorder()
// // 	r.Get("/:author/:title", func(w http.ResponseWriter, r *http.Request) {
// // 		query := r.URL.Query()
// // 		fmt.Fprintf(w, "/%s/%s", query.Get(":author"), query.Get(":title"))
// // 	})
// // 	r.ServeHTTP(rec, req)
// // 	w := rec.Result()
// // 	res, err := ioutil.ReadAll(w.Body)
// // 	is.NoErr(err)
// // 	is.Equal("/coelho/alchemist", string(res))
// // }

// // func TestOptional(t *testing.T) {
// is := is.New(t)
// // 	r := routes.New()
// // 	req := httptest.NewRequest("GET", "/coelho", nil)
// // 	rec := httptest.NewRecorder()
// // 	r.Get("/:author/:title?", func(w http.ResponseWriter, r *http.Request) {
// // 		query := r.URL.Query()
// // 		fmt.Fprintf(w, "/%s/%s", query.Get(":author"), query.Get(":title"))
// // 	})
// // 	r.ServeHTTP(rec, req)
// // 	w := rec.Result()
// // 	res, err := ioutil.ReadAll(w.Body)
// // 	is.NoErr(err)
// // 	is.Equal("/coelho/", string(res))
// // }

// // func TestWildcard(t *testing.T) {
// is := is.New(t)
// // 	r := routes.New()
// // 	req := httptest.NewRequest("GET", "/coelho/alchemist/1988", nil)
// // 	rec := httptest.NewRecorder()
// // 	r.Get("/:author/*", func(w http.ResponseWriter, r *http.Request) {
// // 		query := r.URL.Query()
// // 		fmt.Fprintf(w, "/%s/%s", query.Get(":author"), query.Get(":wild"))
// // 	})
// // 	r.ServeHTTP(rec, req)
// // 	w := rec.Result()
// // 	res, err := ioutil.ReadAll(w.Body)
// // 	is.NoErr(err)
// // 	is.Equal("/coelho/alchemist/1988", string(res))
// // }
