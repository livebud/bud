package virt_test

import (
	"testing"

	"github.com/livebud/bud/pkg/virt"
	"github.com/matryer/is"
)

func TestTree(t *testing.T) {
	is := is.New(t)
	tree := virt.Tree{}
	tree["duo/view/_index.jsx"] = &virt.File{}
	tree["duo/view/_new.svelte"] = &virt.File{}
	tree["duo/view/_ssr.js"] = &virt.File{}
	tree["another/whatever/cool/story.jsx"] = &virt.File{}
	tree["duo/view/index.jsx"] = &virt.File{}
	tree["duo/view/new.svelte"] = &virt.File{}
	tree["duo/controller/show.go"] = &virt.File{}
	tree["duo/controller/new.go"] = &virt.File{}
	actual, err := virt.Print(tree)
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
	tree := virt.Tree{}
	tree["duo/view/_new.svelte"] = &virt.File{}
	actual, err := virt.Print(tree)
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
	tree := virt.Tree{}
	tree["duo/view/_new.svelte"] = &virt.File{}
	tree["duo/view/interesting/_new.svelte"] = &virt.File{}
	tree["nothing/in/common/woah"] = &virt.File{}
	actual, err := virt.Print(tree)
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
