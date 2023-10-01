package lexer_test

import (
	"testing"

	"github.com/livebud/bud/mux/internal/lexer"
	"github.com/matthewmueller/diff"
)

func print(input string) string {
	return lexer.Print(input)
}

func equal(t *testing.T, input, expected string) {
	t.Helper()
	t.Run(input, func(t *testing.T) {
		t.Helper()
		actual := print(input)
		diff.TestString(t, expected, actual)
	})
}

func TestEnd(t *testing.T) {
	equal(t, ``, ``)
}

func TestSample(t *testing.T) {
	equal(t, `/{name}`, `/ { slot:"name" }`)
	equal(t, `/{na me}`, `/ { slot:"na" error:"invalid character ' ' in slot" error:"expected '}' but got 'me'" }`)
	equal(t, `/hello/{name}`, `/ path:"hello" / { slot:"name" }`)
	equal(t, `hello/{name}`, `error:"path must start with a slash /" / { slot:"name" }`)
	equal(t, `/hello/{name}/`, `/ path:"hello" / { slot:"name" } /`)
	equal(t, `/hello/{name?}`, `/ path:"hello" / { slot:"name" ? }`)
	equal(t, `/hello/{name*}`, `/ path:"hello" / { slot:"name" * }`)
	equal(t, `/hel_lo/`, `/ path:"hel_lo" /`)
	equal(t, `/hel lo/`, `/ path:"hel" error:"unexpected character ' ' in path" path:"lo" /`)
	equal(t, `/hello/{*name}`, `/ path:"hello" / { error:"slot can't start with '*'" slot:"name" }`)
	equal(t, `/hello/{na*me}`, `/ path:"hello" / { slot:"na" * error:"expected '}' but got 'me'" }`)
	equal(t, `/hello/{name}/admin`, `/ path:"hello" / { slot:"name" } / path:"admin"`)
	equal(t, `/hello/{name}/admin/`, `/ path:"hello" / { slot:"name" } / path:"admin" /`)
	equal(t, `/hello/{name}/{owner}`, `/ path:"hello" / { slot:"name" } / { slot:"owner" }`)
	equal(t, `/hello/{name`, `/ path:"hello" / { slot:"name" error:"unclosed slot"`)
	equal(t, `/hello/{name|[A-Za-z]}`, `/ path:"hello" / { slot:"name" | regexp:"[A-Za-z]" }`)
	equal(t, `/hello/{name|[A-Za-z]{1,3}}`, `/ path:"hello" / { slot:"name" | regexp:"[A-Za-z]{1,3}" }`)
}

func TestAll(t *testing.T) {
	equal(t, "a", `error:"path must start with a slash /"`)
	equal(t, "/a/", `/ path:"a" /`)
	equal(t, "/", `/`)
	equal(t, "/hi", `/ path:"hi"`)
	equal(t, "/doc/", `/ path:"doc" /`)
	equal(t, "/doc/go_faq.html", `/ path:"doc" / path:"go_faq.html"`)
	equal(t, "α", `error:"path must start with a slash /"`)
	equal(t, "/α", `/ path:"α"`)
	equal(t, "/β", `/ path:"β"`)
	equal(t, "/ ", `/ error:"unexpected character ' ' in path"`)
	equal(t, "/café", `/ path:"café"`)
	equal(t, "/{", `/ { error:"unclosed slot"`)
	equal(t, "/:a", `/ error:"unexpected character ':' in path" path:"a"`)
	equal(t, "/{a}", `/ { slot:"a" }`)
	equal(t, "/{hi}", `/ { slot:"hi" }`)
	equal(t, "/{hi?}", `/ { slot:"hi" ? }`)
	equal(t, "/{hi*}", `/ { slot:"hi" * }`)
	equal(t, "/{hi*?}", `/ { slot:"hi" * error:"expected '}' but got '?'" }`)
	equal(t, "/{hi?*}", `/ { slot:"hi" ? error:"expected '}' but got '*'" }`)
	equal(t, "/a?*", `/ path:"a" error:"unexpected character '?*' in path"`)
	equal(t, "/{a}/{b}", `/ { slot:"a" } / { slot:"b" }`)
	equal(t, "/{a}/{b?}", `/ { slot:"a" } / { slot:"b" ? }`)
	equal(t, "/{a}/{b*}", `/ { slot:"a" } / { slot:"b" * }`)
	equal(t, "/users/{id}.{format?}", `/ path:"users" / { slot:"id" } path:"." { slot:"format" ? }`)
	equal(t, "/users/{major}.{minor}", `/ path:"users" / { slot:"major" } path:"." { slot:"minor" }`)
	equal(t, "/users/{major}.{minor}", `/ path:"users" / { slot:"major" } path:"." { slot:"minor" }`)
	equal(t, "/{-a}/{b?}", `/ { error:"slot can't start with '-'" slot:"a" } / { slot:"b" ? }`)
	equal(t, "/{a-b}", `/ { slot:"a" error:"invalid character '-' in slot" error:"expected '}' but got 'b'" }`)
	equal(t, "{a}", `error:"path must start with a slash /"`)
	equal(t, "/{a}/", `/ { slot:"a" } /`)
	equal(t, "/{café}", `/ { slot:"caf" error:"invalid character 'é' in slot" }`)
	equal(t, "/about", `/ path:"about"`)
	equal(t, "/deactivate", `/ path:"deactivate"`)
	equal(t, "/archive/{year}/{month}", `/ path:"archive" / { slot:"year" } / { slot:"month" }`)
	equal(t, "/users", `/ path:"users"`)
	equal(t, "/users/{id}", `/ path:"users" / { slot:"id" }`)
	equal(t, "/{id}", `/ { slot:"id" }`)
	equal(t, "/v.{version}", `/ path:"v." { slot:"version" }`)
	equal(t, "/{post_id}.{format}", `/ { slot:"post_id" } path:"." { slot:"format" }`)
	equal(t, "/{from}-{to}", `/ { slot:"from" } path:"-" { slot:"to" }`)
	equal(t, "/{key1}/{key2}", `/ { slot:"key1" } / { slot:"key2" }`)
	equal(t, "/{id?}", `/ { slot:"id" ? }`)
	equal(t, "/v.{version?}", `/ path:"v." { slot:"version" ? }`)
	equal(t, "/{post_id}.{format?}", `/ { slot:"post_id" } path:"." { slot:"format" ? }`)
	equal(t, "/{from}-{to?}", `/ { slot:"from" } path:"-" { slot:"to" ? }`)
	equal(t, "/{from}/{to?}", `/ { slot:"from" } / { slot:"to" ? }`)
	equal(t, "/{id}/{path*}", `/ { slot:"id" } / { slot:"path" * }`)
	equal(t, "/v.{version*}", `/ path:"v." { slot:"version" * }`)
	equal(t, "/explore", `/ path:"explore"`)
	equal(t, "/Explore", `/ error:"unexpected character 'E' in path" path:"xplore"`)
	equal(t, "/eXPLORE", `/ path:"e" error:"unexpected character 'XPLORE' in path"`)
	equal(t, "/explorE", `/ path:"explor" error:"unexpected character 'E' in path"`)
	equal(t, "/explorE/", `/ path:"explor" error:"unexpected character 'E' in path" /`)
	equal(t, "/{Slot}", `/ { error:"slot can't start with 'S'" slot:"lot" }`)
	equal(t, "/{sLot}", `/ { slot:"s" error:"invalid character 'L' in slot" error:"expected '}' but got 'ot'" }`)
	equal(t, "/{sloT}", `/ { slot:"slo" error:"invalid character 'T' in slot" }`)
	equal(t, "/{sloT}/", `/ { slot:"slo" error:"invalid character 'T' in slot" } /`)
}
