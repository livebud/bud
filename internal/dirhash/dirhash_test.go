package dirhash_test

import (
	"testing"
	"testing/fstest"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/dirhash"
	"gitlab.com/mnm/bud/internal/gitignore"
)

func TestHash(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main`)},
	}
	hash, err := dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(hash, "Fiyf9IKN3Y0")
	hash, err = dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(hash, "Fiyf9IKN3Y0")

	fsys["main.go"].Data = []byte(`package main2`)
	hash, err = dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(hash, "0yqocFBaFWI")
	hash, err = dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(hash, "0yqocFBaFWI")

	fsys["another.go"] = &fstest.MapFile{Data: []byte(`package main`)}
	hash, err = dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(hash, "FD_ugjQCgCs")

	fsys["something.go"] = &fstest.MapFile{}
	hash, err = dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(hash, "_9y2svTYSsc")
}

func TestSkip(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"main.go": &fstest.MapFile{Data: []byte(`package main`)},
	}
	hash, err := dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(hash, "Fiyf9IKN3Y0")

	fsys["node_modules/svelte/index.js"] = &fstest.MapFile{Data: []byte(`export default svelte`)}
	hash, err = dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(hash, "BJFOfNZj8vk")

	// Skip node_modules
	skipModules := func(path string, isDir bool) bool {
		return isDir && path == "node_modules"
	}
	hash, err = dirhash.Hash(fsys, dirhash.WithSkip(skipModules))
	is.NoErr(err)
	is.Equal(hash, "Fiyf9IKN3Y0")
}

func TestGitIgnore(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{
		"main.go":    &fstest.MapFile{Data: []byte(`package main`)},
		".gitignore": &fstest.MapFile{Data: []byte("node_modules\n")},
	}
	hash, err := dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(hash, "vs70i6JSubE")

	fsys["node_modules/svelte/index.js"] = &fstest.MapFile{Data: []byte(`export default svelte`)}
	hash, err = dirhash.Hash(fsys)
	is.NoErr(err)
	is.Equal(hash, "eU_TkhGsqxk")

	// Skip node_modules
	hash, err = dirhash.Hash(fsys, dirhash.WithSkip(gitignore.FromFS(fsys)))
	is.NoErr(err)
	is.Equal(hash, "vs70i6JSubE")
}
