package radix_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/livebud/bud/pkg/mux/internal/radix"
	"github.com/matryer/is"
	"github.com/matthewmueller/diff"
)

func insertEqual(t *testing.T, tree *radix.Tree, route string, expected string) {
	t.Helper()
	t.Run(route, func(t *testing.T) {
		t.Helper()
		if err := tree.Insert(route, http.NotFoundHandler()); err != nil {
			if err.Error() == expected {
				return
			}
			t.Fatal(err)
		}
		actual := strings.TrimSpace(tree.String())
		expected = strings.ReplaceAll(strings.TrimSpace(expected), "\t", "")
		diff.TestString(t, expected, actual)
	})
}

// https://en.wikipedia.org/wiki/Radix_tree#Insertion
func TestWikipediaInsert(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/test", `
		/test [routable=/test]
	`)
	insertEqual(t, tree, "/slow", `
		/
		•test [routable=/test]
		•slow [routable=/slow]
	`)
	insertEqual(t, tree, "/water", `
		/
		•test [routable=/test]
		•slow [routable=/slow]
		•water [routable=/water]
	`)
	insertEqual(t, tree, "/slower", `
		/
		•test [routable=/test]
		•slow [routable=/slow]
		•••••er [routable=/slower]
		•water [routable=/water]
	`)
	tree = radix.New()
	insertEqual(t, tree, "/tester", `
		/tester [routable=/tester]
	`)
	insertEqual(t, tree, "/test", `
		/test [routable=/test]
		•••••er [routable=/tester]
	`)
	tree = radix.New()
	insertEqual(t, tree, "/test", `
		/test [routable=/test]
	`)
	insertEqual(t, tree, "/team", `
		/te
		•••st [routable=/test]
		•••am [routable=/team]
	`)
	insertEqual(t, tree, "/toast", `
		/t
		••e
		•••st [routable=/test]
		•••am [routable=/team]
		••oast [routable=/toast]
	`)
}

func TestSampleInsert(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/hello/{name}", `
		/hello/{name} [routable=/hello/{name}]
	`)
	insertEqual(t, tree, "/howdy/{name}/", `
		/h
		••ello/{name} [routable=/hello/{name}]
		••owdy/{name} [routable=/howdy/{name}]
	`)
	insertEqual(t, tree, "/hello/{name}/elsewhere", `
		/h
		••ello/{name} [routable=/hello/{name}]
		•••••••••••••/elsewhere [routable=/hello/{name}/elsewhere]
		••owdy/{name} [routable=/howdy/{name}]
	`)
	insertEqual(t, tree, "/hello/{name}/admin/", `
		/h
		••ello/{name} [routable=/hello/{name}]
		•••••••••••••/
		••••••••••••••elsewhere [routable=/hello/{name}/elsewhere]
		••••••••••••••admin [routable=/hello/{name}/admin]
		••owdy/{name} [routable=/howdy/{name}]
	`)
	insertEqual(t, tree, "/hello/{name}/else/", `
		/h
		••ello/{name} [routable=/hello/{name}]
		•••••••••••••/
		••••••••••••••else [routable=/hello/{name}/else]
		••••••••••••••••••where [routable=/hello/{name}/elsewhere]
		••••••••••••••admin [routable=/hello/{name}/admin]
		••owdy/{name} [routable=/howdy/{name}]
	`)
}

func TestEqual(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/hello/{name}", `
		/hello/{name} [routable=/hello/{name}]
	`)
	insertEqual(t, tree, "/hello/{name}", `route already exists "/hello/{name}"`)
	insertEqual(t, tree, "/hello", `
		/hello [routable=/hello]
		••••••/{name} [routable=/hello/{name}]
	`)
	insertEqual(t, tree, "/hello", `route already exists "/hello"`)
	tree = radix.New()
	insertEqual(t, tree, "/{name}", `
		/{name} [routable=/{name}]
	`)
	insertEqual(t, tree, "/{title}", `route "/{title}" is ambiguous with "/{name}"`)
}

func TestDifferentSlots(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/{name}", `
		/{name} [routable=/{name}]
	`)
	insertEqual(t, tree, "/{first}/{last}", `
		/{name} [routable=/{name}]
		•••••••/{last} [routable=/{first}/{last}]
	`)
	insertEqual(t, tree, "/{first}/else", `
		/{name} [routable=/{name}]
		•••••••/
		••••••••else [routable=/{first}/else]
		••••••••{last} [routable=/{first}/{last}]
	`)
	tree = radix.New()
	insertEqual(t, tree, "/{name}", `
		/{name} [routable=/{name}]
	`)
	insertEqual(t, tree, "/else", `
		/
		•else [routable=/else]
		•{name} [routable=/{name}]
	`)
}

func TestPathAfter(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/{name}", `
		/{name} [routable=/{name}]
	`)
	insertEqual(t, tree, "/", `
		/ [routable=/]
		•{name} [routable=/{name}]
	`)
	insertEqual(t, tree, "/first/{name}", `
		/ [routable=/]
		•first/{name} [routable=/first/{name}]
		•{name} [routable=/{name}]
	`)
	insertEqual(t, tree, "/first", `
		/ [routable=/]
		•first [routable=/first]
		••••••/{name} [routable=/first/{name}]
		•{name} [routable=/{name}]
	`)
}

func TestOptionals(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/{name?}", `
		/ [routable=/]
		•{name} [routable=/{name}]
	`)
	insertEqual(t, tree, "/first/{last?}", `
		/ [routable=/]
		•first [routable=/first]
		••••••/{last} [routable=/first/{last}]
		•{name} [routable=/{name}]
	`)
	insertEqual(t, tree, "/{first}/{last}", `
		/ [routable=/]
		•first [routable=/first]
		••••••/{last} [routable=/first/{last}]
		•{name} [routable=/{name}]
		•••••••/{last} [routable=/{first}/{last}]
	`)
	insertEqual(t, tree, "/first/else", `
		/ [routable=/]
		•first [routable=/first]
		••••••/
		•••••••else [routable=/first/else]
		•••••••{last} [routable=/first/{last}]
		•{name} [routable=/{name}]
		•••••••/{last} [routable=/{first}/{last}]
	`)
}

func TestWildcards(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/{name*}", `
		/ [routable=/]
		•{name*} [routable=/{name*}]
	`)
	insertEqual(t, tree, "/first/{last*}", `
		/ [routable=/]
		•first [routable=/first]
		••••••/{last*} [routable=/first/{last*}]
		•{name*} [routable=/{name*}]
	`)
	insertEqual(t, tree, "/{first}/{last}", `
		/ [routable=/]
		•first [routable=/first]
		••••••/{last*} [routable=/first/{last*}]
		•{name*} [routable=/{name*}]
		••••••••/{last} [routable=/{first}/{last}]
	`)
	insertEqual(t, tree, "/first/else", `
		/ [routable=/]
		•first [routable=/first]
		••••••/
		•••••••else [routable=/first/else]
		•••••••{last*} [routable=/first/{last*}]
		•{name*} [routable=/{name*}]
		••••••••/{last} [routable=/{first}/{last}]
	`)
}

func TestInsertRegexp(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/{name|[A-Z]}", `
		/{name|^[A-Z]$} [routable=/{name|^[A-Z]$}]
	`)
	insertEqual(t, tree, "/{path|[0-9]}", `
		/
		•{name|^[A-Z]$} [routable=/{name|^[A-Z]$}]
		•{path|^[0-9]$} [routable=/{path|^[0-9]$}]
	`)
	insertEqual(t, tree, "/{digits|^[0-9]$}", `route "/{digits|^[0-9]$}" is ambiguous with "/{path|^[0-9]$}"`)
	insertEqual(t, tree, "/first/last", `
		/
		•first/last [routable=/first/last]
		•{name|^[A-Z]$} [routable=/{name|^[A-Z]$}]
		•{path|^[0-9]$} [routable=/{path|^[0-9]$}]
	`)
	insertEqual(t, tree, "/{name}", `
		/
		•first/last [routable=/first/last]
		•{name|^[A-Z]$} [routable=/{name|^[A-Z]$}]
		•{path|^[0-9]$} [routable=/{path|^[0-9]$}]
		•{name} [routable=/{name}]
	`)
	insertEqual(t, tree, "/{last*}", `route "/{last*}" is ambiguous with "/{name}"`)
	insertEqual(t, tree, "/first/{last*}", `
		/ [routable=/]
		•first [routable=/first]
		••••••/
		•••••••last [routable=/first/last]
		•••••••{last*} [routable=/first/{last*}]
		•{name|^[A-Z]$} [routable=/{name|^[A-Z]$}]
		•{path|^[0-9]$} [routable=/{path|^[0-9]$}]
		•{name} [routable=/{name}]
	`)
	insertEqual(t, tree, "/{path|[0-9]+}", `
		/ [routable=/]
		•first [routable=/first]
		••••••/
		•••••••last [routable=/first/last]
		•••••••{last*} [routable=/first/{last*}]
		•{name|^[A-Z]$} [routable=/{name|^[A-Z]$}]
		•{path|^[0-9]$} [routable=/{path|^[0-9]$}]
		•{path|^[0-9]+$} [routable=/{path|^[0-9]+$}]
		•{name} [routable=/{name}]
	`)
}

func TestInsertRegexpSlotFirst(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/{name}", `
		/{name} [routable=/{name}]
	`)
	insertEqual(t, tree, "/{path|[A-Z]+}", `
		/
		•{path|^[A-Z]+$} [routable=/{path|^[A-Z]+$}]
		•{name} [routable=/{name}]
	`)
}

func TestRootSwap(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/hello", `
		/hello [routable=/hello]
	`)
	insertEqual(t, tree, "/", `
		/ [routable=/]
		•hello [routable=/hello]
	`)
}

func TestPriority(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/v{version}", `
	/v{version} [routable=/v{version}]
	`)
	insertEqual(t, tree, "/v2", `
	/v
	••2 [routable=/v2]
	••{version} [routable=/v{version}]
	`)
	tree = radix.New()
	insertEqual(t, tree, "/v{version}", `
		/v{version} [routable=/v{version}]
	`)
	insertEqual(t, tree, "/v{major}.{minor}.{patch}", `
		/v{version} [routable=/v{version}]
		•••••••••••.{minor}.{patch} [routable=/v{major}.{minor}.{patch}]
	`)
}

func TestSlotSplit(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/users/{id}/edit", `
		/users/{id}/edit [routable=/users/{id}/edit]
	`)
	insertEqual(t, tree, "/users/settings", `
		/users/
		•••••••settings [routable=/users/settings]
		•••••••{id}/edit [routable=/users/{id}/edit]
	`)
}

func TestInvalidSlot(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/{a}", `
		/{a} [routable=/{a}]
	`)
	insertEqual(t, tree, "/{a}{b}", `slot "a" can't have another slot after`)
}

type Routes []Route

type Route struct {
	Route    string
	Requests Requests
}

type Requests []Request

type Request struct {
	Path   string
	Expect string
}

func matchPath(t *testing.T, tree *radix.Tree, path string, expect string) {
	t.Helper()
	match, err := tree.Match(path)
	if err != nil {
		if err.Error() == expect {
			return
		}
		t.Fatal(err.Error())
	}
	actual := match.String()
	if match.Route != nil && match.Handler == nil {
		t.Fatalf("routes should always have a handler")
	} else if match.Handler != nil && match.Route == nil {
		t.Fatalf("handlers should always have a route")
	}
	diff.TestString(t, expect, actual)
}

func matchEqual(t *testing.T, routes Routes) {
	t.Helper()
	tree := radix.New()
	for _, route := range routes {
		if err := tree.Insert(route.Route, http.NotFoundHandler()); err != nil {
			t.Fatal(err)
		}
		for _, request := range route.Requests {
			t.Run(route.Route, func(t *testing.T) {
				t.Helper()
				matchPath(t, tree, request.Path, request.Expect)
			})
		}
	}
}

func TestSampleMatch(t *testing.T) {
	matchEqual(t, Routes{})
	matchEqual(t, Routes{
		{"/hello", Requests{
			{"/hello", `/hello`},
			{"/hello/world", `no match for "/hello/world"`},
			{"/", `no match for "/"`},
			{"/hello/", `/hello`},
		}},
	})
	matchEqual(t, Routes{
		{"/hello", Requests{
			{"/hello", `/hello`},
			{"/hello/world", `no match for "/hello/world"`},
			{"/hello/", `/hello`},
		}},
		{"/", Requests{
			{"/hello", `/hello`},
			{"/hello/world", `no match for "/hello/world"`},
			{"/", `/`},
			{"/hello/", `/hello`},
		}},
	})
	matchEqual(t, Routes{
		{"/v{version}", Requests{
			{"/v2", "/v{version} version=2"},
		}},
		{"/v{major}.{minor}.{patch}", Requests{
			{"/v2", "/v{version} version=2"},
			{"/v2.0.1", "/v{major}.{minor}.{patch} major=2&minor=0&patch=1"},
		}},
		{"/v1", Requests{
			{"/v2", "/v{version} version=2"},
			{"/v2.0.1", "/v{major}.{minor}.{patch} major=2&minor=0&patch=1"},
			{"/v1", "/v1"},
		}},
		{"/v2.0.0", Requests{
			{"/v2", "/v{version} version=2"},
			{"/v2.0.1", "/v{major}.{minor}.{patch} major=2&minor=0&patch=1"},
			{"/v1", "/v1"},
			{"/v2.0.0", "/v2.0.0"},
		}},
	})
	matchEqual(t, Routes{
		{"/users/{id}/edit", Requests{}},
		{"/users/settings", Requests{}},
		{"/v.{major}.{minor}", Requests{}},
		{"/v.1", Requests{
			{"/v.1.0", `/v.{major}.{minor} major=1&minor=0`},
			{"/users/settings/edit", `/users/{id}/edit id=settings`},
			{"/users/settings", `/users/settings`},
			{"/v.1", `/v.1`},
		}},
	})
}

func TestNonRoutableNoMatch(t *testing.T) {
	is := is.New(t)
	tree := radix.New()
	is.NoErr(tree.Insert("/hello", http.NotFoundHandler()))
	is.NoErr(tree.Insert("/world", http.NotFoundHandler()))
	match, err := tree.Match("/")
	is.True(errors.Is(err, radix.ErrNoMatch))
	is.Equal(match, nil)
}

func TestNoRoutes(t *testing.T) {
	is := is.New(t)
	tree := radix.New()
	match, err := tree.Match("/")
	is.True(errors.Is(err, radix.ErrNoMatch))
	is.Equal(match, nil)
	match, err = tree.Match("/a")
	is.True(errors.Is(err, radix.ErrNoMatch))
	is.Equal(match, nil)
}

func TestAllMatch(t *testing.T) {
	matchEqual(t, Routes{
		{"/hi", Requests{}},
		{"/ab", Requests{}},
		{"/about", Requests{}},
		{"/a", Requests{}},
		{"/α", Requests{}},
		{"/β", Requests{}},
		{"/users", Requests{}},
		{"/users/new", Requests{}},
		{"/users/id", Requests{}},
		{"/users/{id}", Requests{}},
		{"/users/{id}/edit", Requests{}},
		{"/posts/{post_id}/comments", Requests{}},
		{"/posts/{post_id}/comments/new", Requests{}},
		{"/posts/{post_id}/comments/{id}", Requests{}},
		{"/posts/{post_id}/comments/{id}/edit", Requests{}},
		{"/v.{version}", Requests{}},
		{"/v.{major}.{minor}.{patch}", Requests{}},
		{"/v.1", Requests{}},
		{"/v.2.0.0", Requests{}},
		{"/posts/{post_id}.{format}", Requests{}},
		{"/flights/{from}/{to}", Requests{}},
		{"/user/{user}/project/{project}", Requests{}},
		{"/archive/{year}/{month}", Requests{}},
		{"/search/{query}", Requests{
			{"/a", "/a"},
			{"/A", "/a"},
			{"/", `no match for "/"`},
			{"/hi", "/hi"},
			{"/about", "/about"},
			{"/ab", "/ab"},
			{"/abo", `no match for "/abo"`},   // key mismatch
			{"/abou", `no match for "/abou"`}, // key mismatch
			{"/no", `no match for "/no"`},     // no matching child
			{"/α", "/α"},
			{"/β", "/β"},
			{"/αβ", `no match for "/αβ"`},
			{"/users/id", "/users/id"},
			{"/users/10", "/users/{id} id=10"},
			{"/users/1", "/users/{id} id=1"},
			{"/users/a", "/users/{id} id=a"},
			{"/users/-", "/users/{id} id=-"},
			{"/users/_", "/users/{id} id=_"},
			{"/users/abc-d_e", "/users/{id} id=abc-d_e"},
			{"/users/10/edit", "/users/{id}/edit id=10"},
			{"/users/1/edit", "/users/{id}/edit id=1"},
			{"/users/a/edit", "/users/{id}/edit id=a"},
			{"/users/-/edit", "/users/{id}/edit id=-"},
			{"/users/_/edit", "/users/{id}/edit id=_"},
			{"/users/abc-d_e/edit", "/users/{id}/edit id=abc-d_e"},
			{"/posts/1/comments", "/posts/{post_id}/comments post_id=1"},
			{"/posts/10/comments", "/posts/{post_id}/comments post_id=10"},
			{"/posts/a/comments", "/posts/{post_id}/comments post_id=a"},
			{"/posts/-/comments", "/posts/{post_id}/comments post_id=-"},
			{"/posts/_/comments", "/posts/{post_id}/comments post_id=_"},
			{"/posts/abc-d_e/comments", "/posts/{post_id}/comments post_id=abc-d_e"},
			{"/posts/1/comments/2", "/posts/{post_id}/comments/{id} id=2&post_id=1"},
			{"/posts/10/comments/20", "/posts/{post_id}/comments/{id} id=20&post_id=10"},
			{"/posts/a/comments/b", "/posts/{post_id}/comments/{id} id=b&post_id=a"},
			{"/posts/-/comments/-", "/posts/{post_id}/comments/{id} id=-&post_id=-"},
			{"/posts/_/comments/_", "/posts/{post_id}/comments/{id} id=_&post_id=_"},
			{"/posts/abc-d_e/comments/x-y_z", "/posts/{post_id}/comments/{id} id=x-y_z&post_id=abc-d_e"},
			{"/posts/1/comments/2/edit", "/posts/{post_id}/comments/{id}/edit id=2&post_id=1"},
			{"/posts/10/comments/20/edit", "/posts/{post_id}/comments/{id}/edit id=20&post_id=10"},
			{"/posts/a/comments/b/edit", "/posts/{post_id}/comments/{id}/edit id=b&post_id=a"},
			{"/posts/-/comments/-/edit", "/posts/{post_id}/comments/{id}/edit id=-&post_id=-"},
			{"/posts/_/comments/_/edit", "/posts/{post_id}/comments/{id}/edit id=_&post_id=_"},
			{"/posts/abc-d_e/comments/x-y_z/edit", "/posts/{post_id}/comments/{id}/edit id=x-y_z&post_id=abc-d_e"},
			{"/v.1", "/v.1"},
			{"/v.2", "/v.{version} version=2"},
			{"/v.abc", "/v.{version} version=abc"},
			{"/v.2.0.0", "/v.2.0.0"},
			{"/posts/10.json", "/posts/{post_id}.{format} format=json&post_id=10"},
			{"/flights/Berlin/Madison", "/flights/{from}/{to} from=Berlin&to=Madison"},
			{"/archive/2021/2", "/archive/{year}/{month} month=2&year=2021"},
			{"/search/someth!ng+in+ünìcodé", "/search/{query} query=someth%21ng%2Bin%2B%C3%BCn%C3%ACcod%C3%A9"},
			{"/search/with spaces", "/search/{query} query=with+spaces"},
			{"/search/with/slashes", `no match for "/search/with/slashes"`},
		}},
	})
}

func TestMatchUnicode(t *testing.T) {
	matchEqual(t, Routes{
		{"/α", Requests{}},
		{"/β", Requests{}},
		{"/δ", Requests{
			{"/α", `/α`},
			{"/β", `/β`},
			{"/δ", `/δ`},
			{"/Δ", `/δ`},
			{"/αβ", `no match for "/αβ"`},
		}},
	})
}

func TestOptional(t *testing.T) {
	matchEqual(t, Routes{
		{"/{id?}", Requests{}},
		{"/users/{id}.{format?}", Requests{}},
		{"/users/v{version?}", Requests{}},
		{"/flights/{from}/{to?}", Requests{
			{"/", "/"},
			{"/10", "/{id} id=10"},
			{"/a", "/{id} id=a"},

			{"/users/10", `no match for "/users/10"`},
			{"/users/10/", `no match for "/users/10"`},
			{"/users/10.json", "/users/{id}.{format} format=json&id=10"},
			{"/users/10.rss", "/users/{id}.{format} format=rss&id=10"},
			{"/users/index.html", "/users/{id}.{format} format=html&id=index"},
			{"/users/ü.html", "/users/{id}.{format} format=html&id=%C3%BC"},
			{"/users/index.html/more", `no match for "/users/index.html/more"`},

			{"/users", "/{id} id=users"},
			{"/users/", `/{id} id=users`},
			{"/users/v10", "/users/v{version} version=10"},
			{"/users/v1", "/users/v{version} version=1"},

			{"/flights/Berlin", "/flights/{from} from=Berlin"},
			{"/flights/Berlin/", `/flights/{from} from=Berlin`},
			{"/flights/Berlin/Madison", "/flights/{from}/{to} from=Berlin&to=Madison"},
		}},
	})
}

func TestLastOptional(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/slash/{last?}/", `
		/slash [routable=/slash]
		••••••/{last} [routable=/slash/{last}]
	`)
	insertEqual(t, tree, "/not/{last?}/path", `optional slots must be at the end of the path`)
}

func TestWildcard(t *testing.T) {
	matchEqual(t, Routes{
		{"/{path*}", Requests{}},
		{"/users/{id}/{file*}", Requests{}},
		{"/api/v.{version*}", Requests{
			{"/", `/`},
			{"/10", `/{path*} path=10`},
			{"/10/20", `/{path*} path=10%2F20`},
			{"/users/10/dir/file.json", `/users/{id}/{file*} file=dir%2Ffile.json&id=10`},
			{"/users/10/dir", `/users/{id}/{file*} file=dir&id=10`},
			{"/users/10", `/users/{id} id=10`},
			{"/api/v.2/1", `/api/v.{version*} version=2%2F1`},
			{"/api/v.2.1", `/api/v.{version*} version=2.1`},
			{"/api/v.", `/api/v.`},
			{"/api/v", `/{path*} path=api%2Fv`},
		}},
	})
}

func TestLastWildcard(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/slash/{last*}/", `
		/slash [routable=/slash]
		••••••/{last*} [routable=/slash/{last*}]
	`)
	insertEqual(t, tree, "/not/{last*}/path", `wildcard slots must be at the end of the path`)
}

func TestMatchDashedSlots(t *testing.T) {
	matchEqual(t, Routes{
		{"/{a}-{b}", Requests{
			{"/hello-world", `/{a}-{b} a=hello&b=world`},
			{"/a-b", `/{a}-{b} a=a&b=b`},
			{"/A-B", `/{a}-{b} a=A&b=B`},
			{"/AB", `no match for "/AB"`},
		}},
	})
}

func TestBackupTree(t *testing.T) {
	is := is.New(t)
	tree := radix.New()
	insertEqual(t, tree, "/{post_id}/comments", `
		/{post_id}/comments [routable=/{post_id}/comments]
	`)
	insertEqual(t, tree, "/{post_id}.{format}", `
		/{post_id}
		••••••••••/comments [routable=/{post_id}/comments]
		••••••••••.{format} [routable=/{post_id}.{format}]
	`)
	match, err := tree.Match("/10/comments")
	is.NoErr(err)
	is.Equal(match.String(), `/{post_id}/comments post_id=10`)
	match, err = tree.Match("/10.json")
	is.NoErr(err)
	is.Equal(match.String(), `/{post_id}.{format} format=json&post_id=10`)
}

func TestToRoutable(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/last", `
		/last [routable=/last]
	`)
	insertEqual(t, tree, "/first", `
		/
		•last [routable=/last]
		•first [routable=/first]
	`)
	insertEqual(t, tree, "/{last*}", `
		/ [routable=/]
		•last [routable=/last]
		•first [routable=/first]
		•{last*} [routable=/{last*}]
	`)
}

func TestMatchRegexp(t *testing.T) {
	matchEqual(t, Routes{
		{"/{path|[A-Z]}", Requests{
			{"/A", `/{path|^[A-Z]$} path=A`},
			{"/B", `/{path|^[A-Z]$} path=B`},
			{"/Z", `/{path|^[A-Z]$} path=Z`},
			{"/AB", `no match for "/AB"`},
		}},
	})
	matchEqual(t, Routes{
		{"/{path|[A-Z]}", Requests{}},
		{"/{path|[0-9]}", Requests{
			{"/A", `/{path|^[A-Z]$} path=A`},
			{"/0", `/{path|^[0-9]$} path=0`},
			{"/9", `/{path|^[0-9]$} path=9`},
			{"/09", `no match for "/09"`},
		}},
	})
	matchEqual(t, Routes{
		{"/{path|[A-Z]}", Requests{}},
		{"/{path|[0-9]}", Requests{}},
		{"/{path|[A-Z]+}", Requests{
			{"/A", `/{path|^[A-Z]$} path=A`},
			{"/0", `/{path|^[0-9]$} path=0`},
			{"/9", `/{path|^[0-9]$} path=9`},
			{"/09", `no match for "/09"`},
			{"/AB", `/{path|^[A-Z]+$} path=AB`},
		}},
	})
	matchEqual(t, Routes{
		{"/{name}", Requests{}},
		{"/{path|[A-Z]}", Requests{}},
		{"/{path|[0-9]}", Requests{}},
		{"/{path|[A-Z]+}", Requests{
			{"/A", `/{path|^[A-Z]$} path=A`},
			{"/0", `/{path|^[0-9]$} path=0`},
			{"/9", `/{path|^[0-9]$} path=9`},
			{"/09", `/{name} name=09`},
			{"/AB", `/{path|^[A-Z]+$} path=AB`},
		}},
	})
	matchEqual(t, Routes{
		{"/{name}", Requests{}},
		{"/{path|[A-Z]}", Requests{}},
		{"/{path|[0-9]}", Requests{}},
		{"/first", Requests{}},
		{"/{path|[A-Z]+}", Requests{
			{"/A", `/{path|^[A-Z]$} path=A`},
			{"/0", `/{path|^[0-9]$} path=0`},
			{"/9", `/{path|^[0-9]$} path=9`},
			{"/09", `/{name} name=09`},
			{"/AB", `/{path|^[A-Z]+$} path=AB`},
			{"/first", `/first`},
			{"/second", `/{name} name=second`},
		}},
	})
	matchEqual(t, Routes{
		{"/v{version}", Requests{}},
		{"/v{major|[0-9]}.{minor|[0-9]}", Requests{}},
		{"/v{major|[0-9]}.{minor|[0-9]}.{patch|[0-9]}", Requests{
			{"/v1.2.3", `/v{major|^[0-9]$}.{minor|^[0-9]$}.{patch|^[0-9]$} major=1&minor=2&patch=3`},
			{"/v1.2", `/v{major|^[0-9]$}.{minor|^[0-9]$} major=1&minor=2`},
			{"/v1", `/v{version} version=1`},
			{"/valpha.beta.omega", `/v{version} version=alpha.beta.omega`},
		}},
	})
}

func TestResource(t *testing.T) {
	tree := radix.New()
	insertEqual(t, tree, "/{id}/edit", `
		/{id}/edit [routable=/{id}/edit]
	`)
	insertEqual(t, tree, "/", `
		/ [routable=/]
		•{id}/edit [routable=/{id}/edit]
	`)
	matchEqual(t, Routes{
		{"/{id}/edit", Requests{}},
		{"/", Requests{
			{"/", `/`},
			{"/2/edit", `/{id}/edit id=2`},
			{"/3/edit", `/{id}/edit id=3`},
		}},
	})
}

func TestFind(t *testing.T) {
	is := is.New(t)
	tree := radix.New()
	a := http.NotFoundHandler()
	b := http.NotFoundHandler()
	is.NoErr(tree.Insert(`/{post_id}/comments`, a))
	is.NoErr(tree.Insert(`/{post_id}.{format}`, b))
	an, err := tree.Find(`/{post_id}/comments`)
	is.NoErr(err)
	is.Equal(an.Route.String(), `/{post_id}/comments`)
	is.Equal(an.Handler, a)
	bn, err := tree.Find(`/{post_id}.{format}`)
	is.NoErr(err)
	is.Equal(bn.Route.String(), `/{post_id}.{format}`)
	is.Equal(bn.Handler, a)
	cn, err := tree.Find(`/`)
	is.True(errors.Is(err, radix.ErrNoMatch))
	is.Equal(cn, nil)
}

func TestFindByPrefix(t *testing.T) {
	is := is.New(t)
	tree := radix.New()
	a := http.NotFoundHandler()
	is.NoErr(tree.Insert("/", a))
	node, err := tree.FindByPrefix("/{post_id}/comments")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/")
	node, err = tree.FindByPrefix("/{post_id}/")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/")
	node, err = tree.FindByPrefix("/{post_id}")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/")
	node, err = tree.FindByPrefix("/")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/")
	node, err = tree.FindByPrefix("/a")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/")

	// Add nested layout
	is.NoErr(tree.Insert("/{post_id}/comments", a))
	node, err = tree.FindByPrefix("/{post_id}/comments")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/{post_id}/comments")
	node, err = tree.FindByPrefix("/{post_id}/")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/")
	node, err = tree.FindByPrefix("/{post_id}")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/")
	node, err = tree.FindByPrefix("/")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/")

	// No root initially
	tree = radix.New()
	is.NoErr(tree.Insert("/{post_id}/comments", a))
	node, err = tree.FindByPrefix("/{post_id}/comments")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/{post_id}/comments")
	node, err = tree.FindByPrefix("/{post_id}/")
	is.True(errors.Is(err, radix.ErrNoMatch))
	is.Equal(node, nil)
	node, err = tree.FindByPrefix("/{post_id}")
	is.True(errors.Is(err, radix.ErrNoMatch))
	is.Equal(node, nil)
	node, err = tree.FindByPrefix("/")
	is.True(errors.Is(err, radix.ErrNoMatch))
	is.Equal(node, nil)
	node, err = tree.FindByPrefix("/a")
	is.True(errors.Is(err, radix.ErrNoMatch))
	is.Equal(node, nil)
	// Add a subpath
	is.NoErr(tree.Insert("/{post_id}", a))
	node, err = tree.FindByPrefix("/{post_id}/comments")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/{post_id}/comments")
	node, err = tree.FindByPrefix("/{post_id}/")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/{post_id}")
	node, err = tree.FindByPrefix("/{post_id}")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/{post_id}")
	node, err = tree.FindByPrefix("/")
	is.True(errors.Is(err, radix.ErrNoMatch))
	is.Equal(node, nil)
	node, err = tree.FindByPrefix("/a")
	is.True(errors.Is(err, radix.ErrNoMatch))
	is.Equal(node, nil)
	// Add the root
	is.NoErr(tree.Insert("/", a))
	node, err = tree.FindByPrefix("/{post_id}/comments")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/{post_id}/comments")
	node, err = tree.FindByPrefix("/{post_id}/")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/{post_id}")
	node, err = tree.FindByPrefix("/{post_id}")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/{post_id}")
	node, err = tree.FindByPrefix("/")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/")
	node, err = tree.FindByPrefix("/a")
	is.NoErr(err)
	is.Equal(node.Route.String(), "/")
}
