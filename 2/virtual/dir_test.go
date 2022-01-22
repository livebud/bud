package virtual_test

import (
	"io/fs"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/2/virtual"
)

func TestDir(t *testing.T) {
	is := is.New(t)
	fmap := virtual.FileMap()
	// Set the cache
	fmap.Set("view/users", &virtual.Dir{
		Name: "users",
		Entries: []fs.DirEntry{
			&virtual.DirEntry{
				Base: "show.svelte",
				Mode: 0644,
			},
			&virtual.DirEntry{
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
	fmap.Set("view/users", &virtual.Dir{
		Name: "users",
		Entries: []fs.DirEntry{
			&virtual.DirEntry{
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
	fmap := virtual.FileMap()
	fmap.Set("view", &virtual.Dir{
		Name: "view",
		Mode: 0755 &^ fs.ModeDir,
	})
	fi, err := fs.Stat(fmap, "view")
	is.NoErr(err)
	is.Equal(true, fi.IsDir())
}
