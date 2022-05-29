package fscache_test

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/fscache"
	"github.com/livebud/bud/internal/is"
)

func TestFakeDir(t *testing.T) {
	is := is.New(t)
	fmap := fscache.New()
	fmap.Set("index.svelte", &fscache.File{
		Name: "index.svelte",
		Data: []byte(`<h1>hello</h1>`),
		Mode: 0644 | fs.ModeDir,
	})
	fi, err := fs.Stat(fmap, "index.svelte")
	is.NoErr(err)
	is.Equal(false, fi.IsDir())
}

func TestFile(t *testing.T) {
	is := is.New(t)
	fmap := fscache.New()
	fmap.Set("index.svelte", &fscache.File{
		Name: "index.svelte",
		Data: []byte(`<h1>hello</h1>`),
		Mode: 0644,
	})
	code, err := fs.ReadFile(fmap, "index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>hello</h1>`)
	// Read again to test that it's been reset
	code, err = fs.ReadFile(fmap, "index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>hello</h1>`)
}
