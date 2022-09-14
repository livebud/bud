package linkmap_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/budfs/linkmap"
	"golang.org/x/sync/errgroup"
)

func TestLinkMap(t *testing.T) {
	is := is.New(t)
	linkMap := linkmap.New()
	list := linkMap.Scope("bud/view.go")
	list.Add(func(path string) bool {
		return path == "view/index.svelte"
	})
	list.Add(func(path string) bool {
		return path == "view/about/index.svelte"
	})
	list, ok := linkMap.Get("bud/view.go")
	is.True(ok)
	is.True(list.Check("view/about/index.svelte"))
	is.True(list.Check("view/index.svelte"))
	is.True(!list.Check("view"))
}

func TestLinkSafe(t *testing.T) {
	is := is.New(t)
	eg := new(errgroup.Group)
	linkMap := linkmap.New()
	eg.Go(func() error {
		list := linkMap.Scope("bud/view.go")
		list.Add(func(path string) bool {
			return path == "view/index.svelte"
		})
		list.Add(func(path string) bool {
			return path == "view/about/index.svelte"
		})
		list, ok := linkMap.Get("bud/view.go")
		is.True(ok)
		is.True(list.Check("view/about/index.svelte"))
		is.True(list.Check("view/index.svelte"))
		is.True(!list.Check("view"))
		return nil
	})
	eg.Go(func() error {
		list := linkMap.Scope("bud/view.go")
		list.Add(func(path string) bool {
			return path == "view/index.svelte"
		})
		list.Add(func(path string) bool {
			return path == "view/about/index.svelte"
		})
		list, ok := linkMap.Get("bud/view.go")
		is.True(ok)
		is.True(list.Check("view/about/index.svelte"))
		is.True(list.Check("view/index.svelte"))
		is.True(!list.Check("view"))
		return nil
	})
	is.NoErr(eg.Wait())
}

func TestLinkRange(t *testing.T) {
	is := is.New(t)
	linkMap := linkmap.New()
	expect := map[string]bool{
		"bud/view.go":       false,
		"bud/controller.go": false,
	}
	l1 := linkMap.Scope("bud/view.go")
	l1.Add(func(path string) bool {
		return path == "view/index.svelte"
	})
	l2 := linkMap.Scope("bud/controller.go")
	l2.Add(func(path string) bool {
		return path == "view/index.svelte"
	})
	linkMap.Range(func(path string, list *linkmap.List) bool {
		expect[path] = true
		is.True(list.Check("view/index.svelte"))
		return true
	})
	is.True(expect["bud/view.go"])
	is.True(expect["bud/controller.go"])
}
