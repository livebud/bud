package cachefs_test

import (
	"errors"
	"io/fs"
	"testing"

	"gitlab.com/mnm/bud/2/singleflight"
	"gitlab.com/mnm/bud/vfs"

	"github.com/matryer/is"

	"gitlab.com/mnm/bud/2/cachefs"
)

func TestRead(t *testing.T) {
	is := is.New(t)
	loader := new(singleflight.Loader)
	vfs := vfs.Map{
		"a.txt": []byte("a"),
	}
	cache := cachefs.New(vfs, loader, cachefs.Cache())
	code, err := fs.ReadFile(cache, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "a")
	vfs["a.txt"] = []byte("b")
	code, err = fs.ReadFile(cache, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "a")
}

func TestStat(t *testing.T) {
	is := is.New(t)
	loader := new(singleflight.Loader)
	vfs := vfs.Map{
		"a.txt": []byte("a"),
	}
	cache := cachefs.New(vfs, loader, cachefs.Cache())
	fi, err := fs.Stat(cache, "a.txt")
	is.NoErr(err)
	is.Equal(fi.Size(), int64(1))
	vfs["a.txt"] = []byte("aa")
	fi, err = fs.Stat(cache, "a.txt")
	is.NoErr(err)
	is.Equal(fi.Size(), int64(1))
}

func TestReadDir(t *testing.T) {
	is := is.New(t)
	loader := new(singleflight.Loader)
	vfs := vfs.Map{
		"a.txt": []byte("a"),
	}
	cache := cachefs.New(vfs, loader, cachefs.Cache())
	des, err := fs.ReadDir(cache, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	vfs["b.txt"] = []byte("b")
	des, err = fs.ReadDir(cache, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
}

func TestUpdate(t *testing.T) {
	is := is.New(t)
	loader := new(singleflight.Loader)
	vfs := vfs.Map{
		"a.txt": []byte("a"),
	}
	store := cachefs.Cache()
	cache := cachefs.New(vfs, loader, store)
	code, err := fs.ReadFile(cache, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "a")
	vfs["a.txt"] = []byte("b")
	store.Update("a.txt")
	code, err = fs.ReadFile(cache, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "b")
}

func TestCreate(t *testing.T) {
	is := is.New(t)
	loader := new(singleflight.Loader)
	vfs := vfs.Map{
		"a.txt": []byte("a"),
	}
	store := cachefs.Cache()
	cache := cachefs.New(vfs, loader, store)
	des, err := fs.ReadDir(cache, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	vfs["b.txt"] = []byte("b")
	store.Create("b.txt")
	des, err = fs.ReadDir(cache, ".")
	is.NoErr(err)
	is.Equal(len(des), 2)
	code, err := fs.ReadFile(cache, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), "b")
}

func TestDelete(t *testing.T) {
	is := is.New(t)
	loader := new(singleflight.Loader)
	vfs := vfs.Map{
		"a.txt": []byte("a"),
		"b.txt": []byte("b"),
	}
	store := cachefs.Cache()
	cache := cachefs.New(vfs, loader, store)
	des, err := fs.ReadDir(cache, ".")
	is.NoErr(err)
	is.Equal(len(des), 2)
	delete(vfs, "b.txt")
	store.Delete("b.txt")
	des, err = fs.ReadDir(cache, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	code, err := fs.ReadFile(cache, "b.txt")
	is.Equal(nil, code)
	is.True(errors.Is(err, fs.ErrNotExist))
}
