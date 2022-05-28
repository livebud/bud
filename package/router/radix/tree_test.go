package radix_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/gitchander/permutation"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/router/radix"
)

type test struct {
	inserts  []*insert
	requests []*request
}

type insert struct {
	route string
	err   string
}

type inserts []*insert

func (is inserts) Len() int      { return len(is) }
func (is inserts) Swap(i, j int) { is[i].route, is[j].route = is[j].route, is[i].route }

type request struct {
	path    string
	route   string
	slots   string
	nomatch bool
}

func handler(route string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(route))
	})
}

func ok(t testing.TB, test *test) {
	is := is.New(t)
	tree := radix.New()
	for _, insert := range test.inserts {
		err := tree.Insert(insert.route, handler(insert.route))
		if err != nil {
			is.Equal(insert.err, err.Error())
			continue
		} else if insert.err != "" {
			is.Equal(insert.err, err)
			continue
		}
	}
	for _, request := range test.requests {
		match, ok := tree.Match(request.path)
		if !ok {
			is.Equal(request.nomatch, true)
			continue
		}
		// Test the route
		is.Equal(match.Route, request.route)
		// Test the handler
		req := httptest.NewRequest("GET", "/", nil)
		rec := httptest.NewRecorder()
		match.Handler.ServeHTTP(rec, req)
		res := rec.Result()
		body, err := ioutil.ReadAll(res.Body)
		is.NoErr(err)
		is.Equal(request.route, string(body))
		// Test slots
		if len(match.Slots) > 0 {
			is.Equal(request.slots, qs(match.Slots))
		}
	}
}

// ok checks that all permutations of the routes are ok
// doesn't work with routes that contain "err".
func okp(t testing.TB, test *test) {
	// Short tests shouldn't permutate
	if testing.Short() {
		ok(t, test)
		return
	}
	p := permutation.New(inserts(test.inserts))
	for p.Next() {
		ok(t, test)
	}
}

func qs(slots radix.Slots) string {
	var out []string
	for _, slot := range slots {
		out = append(out, slot.Key+"="+slot.Value)
	}
	sort.Strings(out)
	return strings.Join(out, "&")
}

func TestDuplicates(t *testing.T) {
	ok(t, &test{
		inserts: []*insert{
			{route: "/hi"},
			{route: "/hi", err: `radix: "/hi" is already in the tree`},
		},
	})
	ok(t, &test{
		inserts: []*insert{
			{route: "/users/:id"},
			{route: "/users/:id", err: `radix: "/users/:id" is already in the tree`},
		},
	})
	ok(t, &test{
		inserts: []*insert{
			{route: "/zap/:id/edit"},
			{route: "/zap/:id/edit", err: `radix: "/zap/:id/edit" is already in the tree`},
		},
	})
	okp(t, &test{
		inserts: []*insert{
			{route: "/zap/:id"},
			{route: "/zap/id"},
			{route: "/zap/:id/edit"},
		},
	})
	okp(t, &test{
		inserts: []*insert{
			{route: "/v.1"},
			{route: "/v.2"},
			{route: "/v."},
		},
	})
}

func TestSlots(t *testing.T) {
	// Order shouldn't matter
	okp(t, &test{
		inserts: []*insert{
			{route: "/v.:version"},
			{route: "/v.:major.:minor.:patch"},
			{route: "/v.1"},
			{route: "/v.2.0.0"},
		},
		requests: []*request{
			{path: "/v.2", route: "/v.:version", slots: "version=2"},
			{path: "/v.2.0.1", route: "/v.:major.:minor.:patch", slots: "major=2&minor=0&patch=1"},
			{path: "/v.1", route: "/v.1"},
			{path: "/v.2.0.0", route: "/v.2.0.0"},
		},
	})
}

func TestTreeBackup(t *testing.T) {
	okp(t, &test{
		inserts: []*insert{
			{route: "/users/:id/edit"},
			{route: "/users/settings"},
			{route: "/v.:major.:minor"},
			{route: "/v.1"},
		},
		requests: []*request{
			{path: "/v.1.0", route: "/v.:major.:minor", slots: "major=1&minor=0"},
			{path: "/users/settings/edit", route: "/users/:id/edit", slots: `id=settings`},
			{path: "/users/settings", route: "/users/settings"},
			{path: "/v.1", route: "/v.1"},
		},
	})
}

func TestAmbiguous(t *testing.T) {
	ok(t, &test{
		inserts: []*insert{
			{route: "/:id"},
			{route: "/:ab", err: `radix: ambiguous routes "/:ab" and "/:id"`},
		},
		requests: []*request{
			{path: "/a", route: "/:id", slots: "id=a"},
		},
	})
}

func TestMatch(t *testing.T) {
	ok(t, &test{
		inserts: []*insert{
			{route: "/hi"},
			{route: "/ab"},
			{route: "/about"},
			{route: "/a"},
			{route: "/α"},
			{route: "/β"},
			{route: "/users"},
			{route: "/users/new"},
			{route: "/users/id"},
			{route: "/users/:id"},
			{route: "/users/:id/edit"},
			{route: "/posts/:post_id/comments"},
			{route: "/posts/:post_id/comments/new"},
			{route: "/posts/:post_id/comments/:id"},
			{route: "/posts/:post_id/comments/:id/edit"},
			{route: "/v.:version"},
			{route: "/v.:major.:minor.:patch"},
			{route: "/v.1"},
			{route: "/v.2.0.0"},
			{route: "/posts/:post_id.:format"},
			{route: "/flights/:from/:to"},
			{route: "/user/:user/project/:project"},
			{route: "/archive/:year/:month"},
			{route: "/search/:query"},
		},
		requests: []*request{
			{path: "/a", route: "/a"},
			{path: "/", nomatch: true},
			{path: "/hi", route: "/hi"},
			{path: "/about", route: "/about"},
			{path: "/ab", route: "/ab"},
			{path: "/abo", nomatch: true},  // key mismatch
			{path: "/abou", nomatch: true}, // key mismatch
			{path: "/no", nomatch: true},   // no matching child
			{path: "/α", route: "/α"},
			{path: "/β", route: "/β"},
			{path: "/αβ", nomatch: true},
			{path: "/users/id", route: "/users/id"},
			{path: "/users/10", route: "/users/:id", slots: "id=10"},
			{path: "/users/1", route: "/users/:id", slots: "id=1"},
			{path: "/users/a", route: "/users/:id", slots: "id=a"},
			{path: "/users/-", route: "/users/:id", slots: "id=-"},
			{path: "/users/_", route: "/users/:id", slots: "id=_"},
			{path: "/users/abc-d_e", route: "/users/:id", slots: "id=abc-d_e"},

			{path: "/users/10/edit", route: "/users/:id/edit", slots: "id=10"},
			{path: "/users/1/edit", route: "/users/:id/edit", slots: "id=1"},
			{path: "/users/a/edit", route: "/users/:id/edit", slots: "id=a"},
			{path: "/users/-/edit", route: "/users/:id/edit", slots: "id=-"},
			{path: "/users/_/edit", route: "/users/:id/edit", slots: "id=_"},
			{path: "/users/abc-d_e/edit", route: "/users/:id/edit", slots: "id=abc-d_e"},

			{path: "/posts/1/comments", route: "/posts/:post_id/comments", slots: "post_id=1"},
			{path: "/posts/10/comments", route: "/posts/:post_id/comments", slots: "post_id=10"},
			{path: "/posts/a/comments", route: "/posts/:post_id/comments", slots: "post_id=a"},
			{path: "/posts/-/comments", route: "/posts/:post_id/comments", slots: "post_id=-"},
			{path: "/posts/_/comments", route: "/posts/:post_id/comments", slots: "post_id=_"},
			{path: "/posts/abc-d_e/comments", route: "/posts/:post_id/comments", slots: "post_id=abc-d_e"},

			{path: "/posts/1/comments/2", route: "/posts/:post_id/comments/:id", slots: "id=2&post_id=1"},
			{path: "/posts/10/comments/20", route: "/posts/:post_id/comments/:id", slots: "id=20&post_id=10"},
			{path: "/posts/a/comments/b", route: "/posts/:post_id/comments/:id", slots: "id=b&post_id=a"},
			{path: "/posts/-/comments/-", route: "/posts/:post_id/comments/:id", slots: "id=-&post_id=-"},
			{path: "/posts/_/comments/_", route: "/posts/:post_id/comments/:id", slots: "id=_&post_id=_"},
			{path: "/posts/abc-d_e/comments/x-y_z", route: "/posts/:post_id/comments/:id", slots: "id=x-y_z&post_id=abc-d_e"},

			{path: "/posts/1/comments/2/edit", route: "/posts/:post_id/comments/:id/edit", slots: "id=2&post_id=1"},
			{path: "/posts/10/comments/20/edit", route: "/posts/:post_id/comments/:id/edit", slots: "id=20&post_id=10"},
			{path: "/posts/a/comments/b/edit", route: "/posts/:post_id/comments/:id/edit", slots: "id=b&post_id=a"},
			{path: "/posts/-/comments/-/edit", route: "/posts/:post_id/comments/:id/edit", slots: "id=-&post_id=-"},
			{path: "/posts/_/comments/_/edit", route: "/posts/:post_id/comments/:id/edit", slots: "id=_&post_id=_"},
			{path: "/posts/abc-d_e/comments/x-y_z/edit", route: "/posts/:post_id/comments/:id/edit", slots: "id=x-y_z&post_id=abc-d_e"},

			{path: "/v.1", route: "/v.1"},
			{path: "/v.2", route: "/v.:version", slots: "version=2"},
			{path: "/v.abc", route: "/v.:version", slots: "version=abc"},
			{path: "/v.2.0.0", route: "/v.2.0.0"},
			{path: "/posts/10.json", route: "/posts/:post_id.:format", slots: "format=json&post_id=10"},
			{path: "/flights/Berlin/Madison", route: "/flights/:from/:to", slots: "from=Berlin&to=Madison"},
			{path: "/archive/2021/2", route: "/archive/:year/:month", slots: "month=2&year=2021"},

			{path: "/search/someth!ng+in+ünìcodé", route: "/search/:query", slots: "query=someth!ng+in+ünìcodé"},
			{path: "/search/with spaces", route: "/search/:query", slots: "query=with spaces"},
			{path: "/search/with/slashes", nomatch: true},
		},
	})
}

func TestOptional(t *testing.T) {
	okp(t, &test{
		inserts: []*insert{
			{route: "/:id?"},
			{route: "/users/:id.:format?"},
			{route: "/users/v:version?"},
			{route: "/flights/:from/:to?"},
		},
		requests: []*request{
			{path: "/", route: "/:id?"},
			{path: "/10", route: "/:id?", slots: `id=10`},
			{path: "/a", route: "/:id?", slots: `id=a`},

			{path: "/users/10", route: "/users/:id.:format?", slots: `id=10`},
			{path: "/users/10/", nomatch: true},
			{path: "/users/10.json", route: "/users/:id.:format?", slots: `format=json&id=10`},
			{path: "/users/10.rss", route: "/users/:id.:format?", slots: `format=rss&id=10`},
			{path: "/users/index.html", route: "/users/:id.:format?", slots: `format=html&id=index`},
			{path: "/users/ü.html", route: "/users/:id.:format?", slots: `format=html&id=ü`},
			{path: "/users/index.html/more", nomatch: true},

			{path: "/users", route: "/users/v:version?"},
			{path: "/users/", nomatch: true},
			{path: "/users/v10", route: "/users/v:version?", slots: `version=10`},
			{path: "/users/v1", route: "/users/v:version?", slots: `version=1`},

			{path: "/flights/Berlin", route: "/flights/:from/:to?", slots: `from=Berlin`},
			{path: "/flights/Berlin/", nomatch: true},
			{path: "/flights/Berlin/Madison", route: "/flights/:from/:to?", slots: `from=Berlin&to=Madison`},
		},
	})
	ok(t, &test{
		inserts: []*insert{
			{route: "/:id?"},
			{route: "/:a.:b?", err: `radix: ambiguous routes "/:a.:b?" and "/:id?"`},
			{route: "/x:id?", err: `radix: "/x:id?" is already in the tree`},
			{route: "/not/:last?/path", err: `route "/not/:last?/path": optional "?" must be at the end`},
			{route: "/slash/:last?/", err: `route "/slash/:last?/": optional "?" must be at the end`},
		},
	})
}

func TestWildcard(t *testing.T) {
	okp(t, &test{
		inserts: []*insert{
			{route: "/:path*"},
			{route: "/users/:id/:file*"},
			{route: "/api/v.:version*"},
		},
		requests: []*request{
			{path: "/10", route: "/:path*", slots: `path=10`},
			{path: "/10/20", route: "/:path*", slots: `path=10/20`},
			{path: "/users/10/dir/file.json", route: "/users/:id/:file*", slots: `file=dir/file.json&id=10`},
			{path: "/api/v.2/1", route: "/api/v.:version*", slots: `version=2/1`},
			{path: "/api/v.2.1", route: "/api/v.:version*", slots: `version=2.1`},
		},
	})
	ok(t, &test{
		inserts: []*insert{
			{route: "/not/:last*/path", err: `route "/not/:last*/path": wildcard "*" must be at the end`},
			{route: "/slash/:last*/", err: `route "/slash/:last*/": wildcard "*" must be at the end`},
		},
	})
}

func TestNoRoutes(t *testing.T) {
	ok(t, &test{
		requests: []*request{
			{path: "/", nomatch: true},
			{path: "/a", nomatch: true},
		},
	})
}
