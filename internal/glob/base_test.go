package glob_test

import (
	"testing"

	"github.com/livebud/bud/internal/glob"
	"github.com/livebud/bud/internal/is"
)

func TestBase(t *testing.T) {
	is := is.New(t)
	test := func(input, expect string) {
		is.Helper()
		is.Equal(glob.Base(input), expect)
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
	test("{a,b}", ".")
	test("/{a,b}", "/")
	test("/{a,b}/", "/")
	test("{a,b}", ".")
	test("/{a,b}", "/")
	test("./{a,b}", ".")
	test("path/to/*.js", "path/to")
	test("/root/path/to/*.js", "/root/path/to")
	test("chapter/foo [bar]/", "chapter")
	test("path/[a-z]", "path")
	test("[a-z]", ".")
	test("path/{to,from}", "path")
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
	test("path/{foo,bar}/", "path")
	test("{controller/**.go,view/**}", ".")
	test("{controller/**.go,view/**}", ".")
}
