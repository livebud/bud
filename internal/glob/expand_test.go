package glob_test

import (
	"testing"

	"github.com/livebud/bud/internal/glob"
	"github.com/livebud/bud/internal/is"
)

func TestExpand(t *testing.T) {
	is := is.New(t)
	test := func(input string, expect ...string) {
		is.Helper()
		actual, err := glob.Expand(input)
		if err != nil {
			is.True(len(expect) > 0, "missing expect")
			is.Equal(err.Error(), expect[0])
			return
		}
		is.Equal(actual, expect)
	}
	test(".", ".")
	test(".*", ".*")
	test("a/*/b", "a/*/b")
	test("a*/.*/b", "a*/.*/b")
	test("*/a/b/c", "*/a/b/c")
	test("*", "*")
	test("*/", "*/")
	test("*/*", "*/*")
	test("*/*/", "*/*/")
	test("**", "**")
	test("**/", "**/")
	test("**/*", "**/*")
	test("**/*/", "**/*/")
	test("/*.js", "/*.js")
	test("*.js", "*.js")
	test("**/*.js", "**/*.js")
	test("{a,b}", "a", "b")
	test("/{a,b}", "/a", "/b")
	test("/{a,b}/", "/a/", "/b/")
	test("./{a,b}", "./a", "./b")
	test("path/to/*.js", "path/to/*.js")
	test("/root/path/to/*.js", "/root/path/to/*.js")
	test("chapter/foo [bar]/", "chapter/foo [bar]/")
	test("path/[a-z]", "path/[a-z]")
	test("[a-z]", "[a-z]")
	test("path/{to,from}", "path/to", "path/from")
	test("path/!/foo", "path/!/foo")
	test("path/?/foo", "path/?/foo")
	test("path/+/foo", "path/+/foo")
	test("path/*/foo", "path/*/foo")
	test("path/@/foo", "path/@/foo")
	test("path/!/foo/", "path/!/foo/")
	test("path/?/foo/", "path/?/foo/")
	test("path/+/foo/", "path/+/foo/")
	test("path/*/foo/", "path/*/foo/")
	test("path/@/foo/", "path/@/foo/")
	test("path/**/*", "path/**/*")
	test("path/**/subdir/foo.*", "path/**/subdir/foo.*")
	test("path/subdir/**/foo.js", "path/subdir/**/foo.js")
	test("path/!subdir/foo.js", "path/!subdir/foo.js")
	test("path/{foo,bar}/", "path/foo/", "path/bar/")
	test("{controller/**.go,view/**}", "controller/**.go", "view/**")
	test("{controller/**.go,view/**}", "controller/**.go", "view/**")
	test("{{controller,view}/**.go,view/**}", "controller/**.go", "view/**.go", "view/**")
	test("{controller,controller}", "controller")
	// TODO support: test("{a,b}/{c,d}", "a/c", "a/d", "b/c", "b/d")
}
