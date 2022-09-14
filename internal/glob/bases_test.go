package glob_test

import (
	"testing"

	"github.com/livebud/bud/internal/glob"
	"github.com/livebud/bud/internal/is"
)

func TestBases(t *testing.T) {
	is := is.New(t)
	test := func(input string, expect ...string) {
		is.Helper()
		actual, err := glob.Bases(input)
		if err != nil {
			is.True(len(expect) > 0, "missing expect")
			is.Equal(err.Error(), expect[0])
			return
		}
		is.Equal(actual, expect)
	}
	test(".", ".")
	test(".*", ".")
	test("a/*/b", "a")
	test("a*/.*/b", ".")
	test("*/a/b/c", ".")
	test("*", ".")
	test("*/", ".")
	test("*/*", ".")
	test("*/*/", ".")
	test("**", ".")
	test("**/", ".")
	test("**/*", ".")
	test("**/*/", ".")
	test("/*.js", "/")
	test("*.js", ".")
	test("**/*.js", ".")
	test("{a,b}", "a", "b")
	test("/{a,b}", "/a", "/b")
	test("/{a,b}/", "/a", "/b")
	test("./{a,b}", "a", "b")
	test("path/to/*.js", "path/to")
	test("/root/path/to/*.js", "/root/path/to")
	test("chapter/foo [bar]/", "chapter")
	test("path/[a-z]", "path")
	test("[a-z]", ".")
	test("path/{to,from}", "path/to", "path/from")
	test("path/!/foo", "path/!/foo")
	test("path/?/foo", "path")
	test("path/+/foo", "path/+/foo")
	test("path/*/foo", "path")
	test("path/@/foo", "path/@/foo")
	test("path/!/foo/", "path/!/foo")
	test("path/?/foo/", "path")
	test("path/+/foo/", "path/+/foo")
	test("path/*/foo/", "path")
	test("path/@/foo/", "path/@/foo")
	test("path/**/*", "path")
	test("path/**/subdir/foo.*", "path")
	test("path/subdir/**/foo.js", "path/subdir")
	test("path/!subdir/foo.js", "path/!subdir/foo.js")
	test("path/{foo,bar}/", "path/foo", "path/bar")
	test("{controller/**.go,view/**}", "controller", "view")
	test("{controller/**.go,view/**}", "controller", "view")
	test("{{controller,view}/**.go,view/**}", "controller", "view")
	test("{to/*.js,to/*.go}", "to")
}
