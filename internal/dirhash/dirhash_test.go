package dirhash_test

import (
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/dirhash"
	"github.com/livebud/bud/internal/gitignore"
	"github.com/livebud/bud/internal/is"
)

func TestHash(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main`)},
	}
	h1, err := dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(len(h1), 11)
	h2, err := dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(len(h2), 11)
	is.True(h1 == h2)

	fsys["main.go"].Data = []byte(`package main2`)
	h3, err := dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(len(h3), 11)
	is.True(h1 != h3)
	h4, err := dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(len(h4), 11)
	is.True(h3 == h4)

	fsys["another.go"] = &fstest.MapFile{Data: []byte(`package main`)}
	h5, err := dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(len(h5), 11)
	is.True(h4 != h5)

	fsys["something.go"] = &fstest.MapFile{}
	h6, err := dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(len(h6), 11)
	is.True(h5 != h6)
}

func TestSkip(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main`)},
	}
	h1, err := dirhash.Hash(fsys)
	is.NoErr(err)

	fsys["node_modules/svelte/index.js"] = &fstest.MapFile{Data: []byte(`export default svelte`)}
	h2, err := dirhash.Hash(fsys)
	is.NoErr(err)
	is.True(h1 != h2)

	// Skip node_modules
	skipModules := func(path string, isDir bool) bool {
		return isDir && path == "node_modules"
	}
	h3, err := dirhash.Hash(fsys, dirhash.WithSkip(skipModules))
	is.NoErr(err)
	is.True(h1 == h3)
}

func TestGitIgnore(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"main.go":    &fstest.MapFile{Data: []byte(`package main`)},
		".gitignore": &fstest.MapFile{Data: []byte("node_modules\n")},
	}
	h1, err := dirhash.Hash(fsys)
	is.NoErr(err)

	fsys["node_modules/svelte/index.js"] = &fstest.MapFile{Data: []byte(`export default svelte`)}
	h2, err := dirhash.Hash(fsys)
	is.NoErr(err)
	is.True(h1 != h2)

	// Skip node_modules
	h3, err := dirhash.Hash(fsys, dirhash.WithSkip(gitignore.FromFS(fsys)))
	is.NoErr(err)
	is.True(h1 == h3)
}
