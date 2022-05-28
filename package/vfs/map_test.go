package vfs_test

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/vfs"
)

func TestMap(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"duo/view/index.svelte": []byte(`<h1>index</h1>`),
	}

	// Read duo/view/index.svelte
	code, err := fs.ReadFile(fsys, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// stat duo/
	stat, err := fs.Stat(fsys, "duo")
	is.NoErr(err)
	is.Equal(stat.Name(), "duo")
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.FileMode(fs.ModeDir))
}
