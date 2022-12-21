package dcache_test

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/internal/testsub"

	"github.com/livebud/bud/internal/dcache"
	"github.com/livebud/bud/internal/is"
)

const dbPath = ":memory:"

func TestPutGet(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cache, err := dcache.Load(ctx, dbPath)
	is.NoErr(err)
	defer cache.Close()
	is.NoErr(cache.Put(ctx, "a.txt", &dcache.File{
		Data: []byte("a.txt"),
		Mode: 0644,
	}))
	file, err := cache.Get(ctx, "a.txt")
	is.NoErr(err)
	is.Equal(file.Data, []byte("a.txt"))
	is.Equal(file.Mode, fs.FileMode(0644))
}

func TestGetNotFound(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cache, err := dcache.Load(ctx, dbPath)
	is.NoErr(err)
	defer cache.Close()
	file, err := cache.Get(ctx, "a.txt")
	is.True(errors.Is(err, dcache.ErrNotFound))
	is.Equal(file, nil)
}

func seed(ctx context.Context, is *is.I, cache dcache.Cache) {
	is.NoErr(cache.Put(ctx, "a.txt", &dcache.File{
		Data:  []byte("a.txt"),
		Mode:  0644,
		Links: []string{"b.txt", "c.txt"},
	}))
	is.NoErr(cache.Put(ctx, "b.txt", &dcache.File{
		Data:  []byte("b.txt"),
		Mode:  0644,
		Links: []string{"e.txt"},
	}))
	is.NoErr(cache.Put(ctx, "e.txt", &dcache.File{
		Data:  []byte("e.txt"),
		Mode:  0644,
		Links: []string{},
	}))
	is.NoErr(cache.Put(ctx, "c.txt", &dcache.File{
		Data:  []byte("c.txt"),
		Mode:  0644,
		Links: []string{"d.txt"},
	}))
	is.NoErr(cache.Put(ctx, "d.txt", &dcache.File{
		Data:  []byte("d.txt"),
		Mode:  0644,
		Links: []string{"f.txt"},
	}))
	is.NoErr(cache.Put(ctx, "f.txt", &dcache.File{
		Data:  []byte("f.txt"),
		Mode:  0644,
		Links: []string{},
	}))
}

func TestAncestors(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cache, err := dcache.Load(ctx, dbPath)
	is.NoErr(err)
	defer cache.Close()
	seed(ctx, is, cache)
	paths, err := cache.Ancestors(ctx, "a.txt")
	is.NoErr(err)
	is.Equal(len(paths), 0)
	paths, err = cache.Ancestors(ctx, "c.txt")
	is.NoErr(err)
	is.Equal(paths, []string{"a.txt"})
	paths, err = cache.Ancestors(ctx, "c.txt", "e.txt")
	is.NoErr(err)
	is.Equal(paths, []string{"a.txt", "b.txt"})
	paths, err = cache.Ancestors(ctx, "f.txt", "e.txt")
	is.NoErr(err)
	is.Equal(paths, []string{"a.txt", "b.txt", "c.txt", "d.txt"})
}

func TestDelete(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cache, err := dcache.Load(ctx, dbPath)
	is.NoErr(err)
	defer cache.Close()
	seed(ctx, is, cache)
	// Changing c.txt deletes c.txt and a.txt
	err = cache.Delete(ctx, "c.txt")
	is.NoErr(err)

	file, err := cache.Get(ctx, "a.txt")
	is.True(errors.Is(err, dcache.ErrNotFound))
	is.Equal(file, nil)
	file, err = cache.Get(ctx, "c.txt")
	is.True(errors.Is(err, dcache.ErrNotFound))
	is.Equal(file, nil)

	files, err := cache.Files(ctx)
	is.NoErr(err)
	is.Equal(len(files), 4)
	is.Equal(files[0].Path, "b.txt")
	is.Equal(files[1].Path, "d.txt")
	is.Equal(files[2].Path, "e.txt")
	is.Equal(files[3].Path, "f.txt")

	links, err := cache.Links(ctx)
	is.NoErr(err)
	is.Equal(len(links), 2)
	is.Equal(links[0].From, "b.txt")
	is.Equal(links[0].To, "e.txt")
	is.Equal(links[1].From, "d.txt")
	is.Equal(links[1].To, "f.txt")
}

func TestPrint(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cache, err := dcache.Load(ctx, dbPath)
	is.NoErr(err)
	defer cache.Close()
	seed(ctx, is, cache)
	dot, err := cache.Print(ctx)
	is.NoErr(err)
	is.Equal(dot, `digraph dag {
	"a.txt"
	"b.txt"
	"c.txt"
	"d.txt"
	"e.txt"
	"f.txt"
	"a.txt" -> "b.txt"
	"a.txt" -> "c.txt"
	"b.txt" -> "e.txt"
	"c.txt" -> "d.txt"
	"d.txt" -> "f.txt"
}
`)
}

func TestManyWrites(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	cache, err := dcache.Load(ctx, dbPath)
	is.NoErr(err)
	defer cache.Close()
	for i := 0; i < 100; i++ {
		is.NoErr(cache.Put(ctx, fmt.Sprintf("%d.txt", i), &dcache.File{
			Data: []byte(fmt.Sprintf("%d.txt", i)),
		}))
	}
}

func TestConcurrentWrites(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	parent := func(t testing.TB, cmd *exec.Cmd) {
		dbPath := filepath.Join(t.TempDir(), "test_concurrent_writes.db")
		cache, err := dcache.Load(ctx, dbPath)
		is.NoErr(err)
		defer cache.Close()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(cmd.Env, "DBPATH="+dbPath)
		is.NoErr(cmd.Start())
		for i := 0; i < 100; i++ {
			is.NoErr(cache.Put(ctx, fmt.Sprintf("%d.txt", i), &dcache.File{
				Data: []byte(fmt.Sprintf("%d.txt", i)),
			}))
		}
		is.NoErr(cmd.Wait())
	}
	child := func(t testing.TB) {
		dbPath := os.Getenv("DBPATH")
		cache, err := dcache.Load(ctx, dbPath)
		is.NoErr(err)
		defer cache.Close()
		for i := 100; i < 200; i++ {
			is.NoErr(cache.Put(ctx, fmt.Sprintf("%d.txt", i), &dcache.File{
				Data: []byte(fmt.Sprintf("%d.txt", i)),
			}))
		}
	}
	testsub.Run(t, parent, child)
}
