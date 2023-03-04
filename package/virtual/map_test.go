package virtual_test

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/virtual"
)

func TestMap(t *testing.T) {
	is := is.New(t)
	fsys := virtual.Map{
		"bud/view/index.svelte":        `<h1>index</h1>`,
		"bud/controller/controller.go": `package controller`,
	}

	// Read bud/view/index.svelte
	code, err := fs.ReadFile(fsys, "bud/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// stat bud/
	stat, err := fs.Stat(fsys, "bud")
	is.NoErr(err)
	is.Equal(stat.Name(), "bud")
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(fs.ModeDir))

	// Test reading the directory
	des, err := fs.ReadDir(fsys, "bud")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "controller")
	is.Equal(des[1].Name(), "view")
}
