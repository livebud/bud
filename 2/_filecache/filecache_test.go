package filecache_test

import (
	"io/fs"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/2/virtual"

	"gitlab.com/mnm/bud/2/filecache"
)

func TestDir(t *testing.T) {
	is := is.New(t)
	fcache := filecache.New()
	// Set the cache
	fcache.Set("view/users", &virtual.Dir{
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
	des, err := fs.ReadDir(fcache, "view/users")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "index.svelte")
	is.Equal(des[1].Name(), "show.svelte")
	// Check again
	des, err = fs.ReadDir(fcache, "view/users")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "index.svelte")
	is.Equal(des[1].Name(), "show.svelte")
	// Adding dir entries doesn't not mean it's still in cache
	info, err := fs.Stat(fcache, "view/users/index.svelte")
	is.Equal(nil, info)
	is.Equal(err, fs.ErrNotExist)
	// Override the directory
	fcache.Set("view/users", &virtual.Dir{
		Name: "users",
		Entries: []fs.DirEntry{
			&virtual.DirEntry{
				Base: "index.svelte",
				Mode: 0644,
			},
		},
	})
	des, err = fs.ReadDir(fcache, "view/users")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "index.svelte")
}
