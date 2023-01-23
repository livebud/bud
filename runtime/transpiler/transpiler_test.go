package transpiler

import (
	"testing"

	"github.com/livebud/bud/internal/is"
)

func TestSplitRoot(t *testing.T) {
	is := is.New(t)
	root, path := splitRoot("foo/bar/baz.svelte")
	is.Equal(root, "foo")
	is.Equal(path, "bar/baz.svelte")
}
