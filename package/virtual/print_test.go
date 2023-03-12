package virtual_test

import (
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/virtual"
)

func TestTree(t *testing.T) {
	is := is.New(t)
	tree := virtual.Tree{}
	tree["duo/view/_index.jsx"] = &virtual.File{}
	tree["duo/view/_new.svelte"] = &virtual.File{}
	tree["duo/view/_ssr.js"] = &virtual.File{}
	tree["another/whatever/cool/story.jsx"] = &virtual.File{}
	tree["duo/view/index.jsx"] = &virtual.File{}
	tree["duo/view/new.svelte"] = &virtual.File{}
	tree["duo/controller/show.go"] = &virtual.File{}
	tree["duo/controller/new.go"] = &virtual.File{}
	actual, err := virtual.Print(tree)
	is.NoErr(err)
	expect := `.
├── another
│   └── whatever
│       └── cool
│           └── story.jsx
└── duo
    ├── controller
    │   ├── new.go
    │   └── show.go
    └── view
        ├── _index.jsx
        ├── _new.svelte
        ├── _ssr.js
        ├── index.jsx
        └── new.svelte
`
	is.Equal(actual, expect)
}

func TestOneDeep(t *testing.T) {
	is := is.New(t)
	tree := virtual.Tree{}
	tree["duo/view/_new.svelte"] = &virtual.File{}
	actual, err := virtual.Print(tree)
	is.NoErr(err)
	expect := `.
└── duo
    └── view
        └── _new.svelte
`
	is.Equal(actual, expect)
}

func TestUncommon(t *testing.T) {
	is := is.New(t)
	tree := virtual.Tree{}
	tree["duo/view/_new.svelte"] = &virtual.File{}
	tree["duo/view/interesting/_new.svelte"] = &virtual.File{}
	tree["nothing/in/common/woah"] = &virtual.File{}
	actual, err := virtual.Print(tree)
	is.NoErr(err)
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
	is.Equal(actual, expect)
}
