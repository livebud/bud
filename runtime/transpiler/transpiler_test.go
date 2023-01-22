package transpiler_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/runtime/transpiler"
)

func TestSplitRoot(t *testing.T) {
	is := is.New(t)
	root, path := transpiler.SplitRoot("foo/bar/baz.svelte")
	is.Equal(root, "foo")
	is.Equal(path, "bar/baz.svelte")
}
