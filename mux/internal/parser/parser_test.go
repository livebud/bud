package parser_test

import (
	"testing"

	"github.com/livebud/bud/mux/internal/parser"
	"github.com/matthewmueller/diff"
)

func equal(t *testing.T, input, expected string) {
	t.Helper()
	t.Run(input, func(t *testing.T) {
		t.Helper()
		route, err := parser.Parse(input)
		if err != nil {
			if err.Error() == expected {
				return
			}
			t.Fatal(err)
		}
		actual := route.String()
		diff.TestString(t, expected, actual)
	})
}

func TestEmpty(t *testing.T) {
	equal(t, ``, ``)
}

func TestSample(t *testing.T) {
	equal(t, `/{name}`, `/{name}`)
	equal(t, `/{na me}`, `invalid character ' ' in slot`)
	equal(t, `/hello/{name}`, `/hello/{name}`)
	equal(t, `hello/{name}`, `path must start with a slash /`)
	equal(t, `/hello/{name}/`, `/hello/{name}/`)
	equal(t, `/hello/{name?}`, `/hello/{name?}`)
	equal(t, `/hello/{name*}`, `/hello/{name*}`)
	equal(t, `/hel_lo/`, `/hel_lo/`)
	equal(t, `/hel lo/`, `unexpected character ' ' in path`)
	equal(t, `/hello/{*name}`, `slot can't start with '*'`)
	equal(t, `/hello/{na*me}`, `expected '}' but got 'me'`)
	equal(t, `/hello/{name}/admin`, `/hello/{name}/admin`)
	equal(t, `/hello/{name}/admin/`, `/hello/{name}/admin/`)
	equal(t, `/hello/{name}/{owner}`, `/hello/{name}/{owner}`)
	equal(t, `/hello/{name`, `unclosed slot`)
	equal(t, `/hello/{name|^[A-Za-z]$}`, `/hello/{name|^[A-Za-z]$}`)
	equal(t, `/hello/{name|^[A-Za-z]{1,3}$}`, `/hello/{name|^[A-Za-z]{1,3}$}`)
}

func TestAll(t *testing.T) {
	equal(t, "a", `path must start with a slash /`)
	equal(t, "/a/", `/a/`)
	equal(t, "/", `/`)
	equal(t, "/hi", `/hi`)
	equal(t, "/doc/", `/doc/`)
	equal(t, "/doc/go_faq.html", `/doc/go_faq.html`)
	equal(t, "α", `path must start with a slash /`)
	equal(t, "/α", `/α`)
	equal(t, "/β", `/β`)
	equal(t, "/ ", `unexpected character ' ' in path`)
	equal(t, "/café", `/café`)
	equal(t, "/{", `unclosed slot`)
	equal(t, "/:a", `unexpected character ':' in path`)
	equal(t, "/{a}", `/{a}`)
	equal(t, "/{hi}", `/{hi}`)
	equal(t, "/{hi?}", `/{hi?}`)
	equal(t, "/{hi*}", `/{hi*}`)
	equal(t, "/{hi*?}", `expected '}' but got '?'`)
	equal(t, "/{hi?*}", `expected '}' but got '*'`)
	equal(t, "/a?*", `unexpected character '?*' in path`)
	equal(t, "/{a}/{b}", `/{a}/{b}`)
	equal(t, "/{a}/{b?}", `/{a}/{b?}`)
	equal(t, "/{a}/{b*}", `/{a}/{b*}`)
	equal(t, "/users/{id}.{format?}", `/users/{id}.{format?}`)
	equal(t, "/users/{major}.{minor}", `/users/{major}.{minor}`)
	equal(t, "/users/{major}.{minor}", `/users/{major}.{minor}`)
	equal(t, "/{-a}/{b?}", `slot can't start with '-'`)
	equal(t, "/{a-b}", `invalid character '-' in slot`)
	equal(t, "{a}", `path must start with a slash /`)
	equal(t, "/{a}/", `/{a}/`)
	equal(t, "/{café}", `invalid character 'é' in slot`)
	equal(t, "/about", `/about`)
	equal(t, "/deactivate", `/deactivate`)
	equal(t, "/archive/{year}/{month}", `/archive/{year}/{month}`)
	equal(t, "/users", `/users`)
	equal(t, "/users/{id}", `/users/{id}`)
	equal(t, "/{id}", `/{id}`)
	equal(t, "/v.{version}", `/v.{version}`)
	equal(t, "/{post_id}.{format}", `/{post_id}.{format}`)
	equal(t, "/{from}-{to}", `/{from}-{to}`)
	equal(t, "/{key1}/{key2}", `/{key1}/{key2}`)
	equal(t, "/{id?}", `/{id?}`)
	equal(t, "/v.{version?}", `/v.{version?}`)
	equal(t, "/{post_id}.{format?}", `/{post_id}.{format?}`)
	equal(t, "/{from}-{to?}", `/{from}-{to?}`)
	equal(t, "/{from}/{to?}", `/{from}/{to?}`)
	equal(t, "/{id}/{path*}", `/{id}/{path*}`)
	equal(t, "/v.{version*}", `/v.{version*}`)
	equal(t, "/explore", `/explore`)
	equal(t, "/Explore", `unexpected character 'E' in path`)
	equal(t, "/eXPLORE", `unexpected character 'XPLORE' in path`)
	equal(t, "/explorE", `unexpected character 'E' in path`)
	equal(t, "/explorE/", `unexpected character 'E' in path`)
	equal(t, "/{Slot}", `slot can't start with 'S'`)
	equal(t, "/{sLot}", `invalid character 'L' in slot`)
	equal(t, "/{sloT}", `invalid character 'T' in slot`)
	equal(t, "/{sloT}/", `invalid character 'T' in slot`)
}
