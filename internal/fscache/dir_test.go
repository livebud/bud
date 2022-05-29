package fscache_test

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/fscache"
	"github.com/livebud/bud/internal/is"
)

func TestDir(t *testing.T) {
	is := is.New(t)
	fmap := fscache.New()
	// Set the cache
	fmap.Set("view/users", &fscache.Dir{
		Name: "users",
		Entries: []fs.DirEntry{
			&fscache.DirEntry{
				Base: "show.svelte",
				Mode: 0644,
			},
			&fscache.DirEntry{
				Base: "index.svelte",
				Mode: 0644,
			},
		},
	})
	// Test the the directory has entries
	des, err := fs.ReadDir(fmap, "view/users")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "index.svelte")
	is.Equal(des[1].Name(), "show.svelte")
	// Check again
	des, err = fs.ReadDir(fmap, "view/users")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "index.svelte")
	is.Equal(des[1].Name(), "show.svelte")
	// Adding dir entries doesn't not mean it's still in cache
	info, err := fs.Stat(fmap, "view/users/index.svelte")
	is.Equal(nil, info)
	is.Equal(err, fs.ErrNotExist)
	// Override the directory
	fmap.Set("view/users", &fscache.Dir{
		Name: "users",
		Entries: []fs.DirEntry{
			&fscache.DirEntry{
				Base: "index.svelte",
				Mode: 0644,
			},
		},
	})
	des, err = fs.ReadDir(fmap, "view/users")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "index.svelte")
}

func TestFakeFile(t *testing.T) {
	is := is.New(t)
	fmap := fscache.New()
	fmap.Set("view", &fscache.Dir{
		Name: "view",
		Mode: 0755 &^ fs.ModeDir,
	})
	fi, err := fs.Stat(fmap, "view")
	is.NoErr(err)
	is.Equal(true, fi.IsDir())
}
