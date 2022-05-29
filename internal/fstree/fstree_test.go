package fstree_test

import (
	"testing"

	"github.com/livebud/bud/internal/fstree"
	"github.com/livebud/bud/internal/is"
)

func TestTree(t *testing.T) {
	is := is.New(t)
	tree := fstree.New()
	tree.Add("duo")
	tree.Add("duo/view/_index.jsx")
	tree.Add("duo/view/_new.svelte")
	tree.Add("duo/view/_ssr.js")
	tree.Add("duo/view")
	tree.Add("another")
	tree.Add("another/whatever/cool/story.jsx")
	tree.Add("duo/view/index.jsx")
	tree.Add("duo/view/new.svelte")
	tree.Add("duo/controller/show.go")
	tree.Add("duo/controller/new.go")
	tree.Add("duo/controller")
	// lines := strings.Split(tree.String(), "\n")
	expect := `.
├── duo
│   ├── view
│   │   ├── _index.jsx
│   │   ├── _new.svelte
│   │   ├── _ssr.js
│   │   ├── index.jsx
│   │   └── new.svelte
│   └── controller
│       ├── show.go
│       └── new.go
└── another
    └── whatever
        └── cool
            └── story.jsx
`
	is.Equal(tree.String(), expect)
}

func TestOneDeep(t *testing.T) {
	is := is.New(t)
	tree := fstree.New()
	tree.Add("duo/view/_new.svelte")
	expect := `.
└── duo
    └── view
        └── _new.svelte
`
	is.Equal(tree.String(), expect)
}
func TestUncommon(t *testing.T) {
	is := is.New(t)
	tree := fstree.New()
	tree.Add("duo/view/_new.svelte")
	tree.Add("duo/view/interesting/_new.svelte")
	tree.Add("nothing/in/common/woah")
	expect := `.
├── duo
│   └── view
│       ├── _new.svelte
│       └── interesting
│           └── _new.svelte
└── nothing
    └── in
        └── common
            └── woah
`
	is.Equal(tree.String(), expect)
}
