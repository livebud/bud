package glob_test

import (
	"testing"

	"github.com/livebud/bud/internal/glob"
	"github.com/livebud/bud/internal/is"
)

func eq(is *is.I, input, expect string) {
	is.Helper()
	is.Equal(glob.Base(input), expect)
}

func TestBase(t *testing.T) {
	is := is.New(t)
	eq(is, ".", ".")
	eq(is, ".*", ".")
	eq(is, "a/*/b", "a")
	eq(is, "a*/.*/b", ".")
	eq(is, "*/a/b/c", ".")
	eq(is, "*", ".")
	eq(is, "*/", ".")
	eq(is, "*/*", ".")
	eq(is, "*/*/", ".")
	eq(is, "**", ".")
	eq(is, "**/", ".")
	eq(is, "**/*", ".")
	eq(is, "**/*/", ".")
	eq(is, "/*.js", "/")
	eq(is, "*.js", ".")
	eq(is, "**/*.js", ".")
	eq(is, "{a,b}", ".")
	eq(is, "/{a,b}", "/")
	eq(is, "/{a,b}/", "/")
	eq(is, "{a,b}", ".")
	eq(is, "/{a,b}", "/")
	eq(is, "./{a,b}", ".")
	eq(is, "path/to/*.js", "path/to")
	eq(is, "/root/path/to/*.js", "/root/path/to")
	eq(is, "chapter/foo [bar]/", "chapter")
	eq(is, "path/[a-z]", "path")
	eq(is, "[a-z]", ".")
	eq(is, "path/{to,from}", "path")
	eq(is, "path/!/foo", "path/!/foo")
	eq(is, "path/?/foo", "path")
	eq(is, "path/+/foo", "path/+/foo")
	eq(is, "path/*/foo", "path")
	eq(is, "path/@/foo", "path/@/foo")
	eq(is, "path/!/foo/", "path/!/foo")
	eq(is, "path/?/foo/", "path")
	eq(is, "path/+/foo/", "path/+/foo")
	eq(is, "path/*/foo/", "path")
	eq(is, "path/@/foo/", "path/@/foo")
	eq(is, "path/**/*", "path")
	eq(is, "path/**/subdir/foo.*", "path")
	eq(is, "path/subdir/**/foo.js", "path/subdir")
	eq(is, "path/!subdir/foo.js", "path/!subdir/foo.js")
	eq(is, "path/{foo,bar}/", "path")
	eq(is, "{controller/**.go,view/**}", ".")
}
