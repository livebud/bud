package linkmap_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/budfs/linkmap"
)

func TestLinkMap(t *testing.T) {
	is := is.New(t)
	linkMap := linkmap.Map{}
	list := linkMap.Scope("bud/internal/app/view/view.go")
	list.Add(func(path string) bool {
		return path == "view/index.svelte"
	})
	list.Add(func(path string) bool {
		return path == "view/about/index.svelte"
	})
	is.True(linkMap["bud/internal/app/view/view.go"].Check("view/about/index.svelte"))
	is.True(linkMap["bud/internal/app/view/view.go"].Check("view/index.svelte"))
	is.True(!linkMap["bud/internal/app/view/view.go"].Check("view"))
}
