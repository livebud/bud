package dag_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/package/virt"

	"github.com/livebud/bud/internal/dag"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testsub"
)

const dbPath = ":memory:"

func TestSetGet(t *testing.T) {
	is := is.New(t)
	fsys := virt.Map{}
	cache, err := dag.Load(fsys, dbPath)
	is.NoErr(err)
	defer cache.Close()
	is.NoErr(cache.Set("a.txt", &virt.File{
		Data: []byte("a.txt"),
		Mode: 0644,
	}))
	file, err := cache.Get("a.txt")
	is.NoErr(err)
	is.Equal(file.Data, []byte("a.txt"))
	is.Equal(file.Mode, fs.FileMode(0644))
}

func TestGetNotFound(t *testing.T) {
	is := is.New(t)
	fsys := virt.Map{}
	cache, err := dag.Load(fsys, dbPath)
	is.NoErr(err)
	defer cache.Close()
	file, err := cache.Get("a.txt")
	is.True(errors.Is(err, dag.ErrNotFound))
	is.Equal(file, nil)
}

func seed(is *is.I, cache *dag.Cache) {
	is.NoErr(cache.Set("a.txt", &virt.File{
		Data: []byte("a.txt"),
		Mode: 0644,
	}))
	is.NoErr(cache.Link("a.txt", "b.txt", "c.txt"))
	is.NoErr(cache.Set("b.txt", &virt.File{
		Data: []byte("b.txt"),
		Mode: 0644,
	}))
	is.NoErr(cache.Link("b.txt", "e.txt"))
	is.NoErr(cache.Set("e.txt", &virt.File{
		Data: []byte("e.txt"),
		Mode: 0644,
	}))
	is.NoErr(cache.Set("c.txt", &virt.File{
		Data: []byte("c.txt"),
		Mode: 0644,
	}))
	is.NoErr(cache.Link("c.txt", "d.txt"))
	is.NoErr(cache.Set("d.txt", &virt.File{
		Data: []byte("d.txt"),
		Mode: 0644,
	}))
	is.NoErr(cache.Link("d.txt", "f.txt"))
	is.NoErr(cache.Set("f.txt", &virt.File{
		Data: []byte("f.txt"),
		Mode: 0644,
	}))
}

func TestAncestors(t *testing.T) {
	is := is.New(t)
	fsys := virt.Map{}
	cache, err := dag.Load(fsys, dbPath)
	is.NoErr(err)
	defer cache.Close()
	seed(is, cache)
	paths, err := cache.Ancestors("a.txt")
	is.NoErr(err)
	is.Equal(len(paths), 0)
	paths, err = cache.Ancestors("c.txt")
	is.NoErr(err)
	is.Equal(paths, []string{"a.txt"})
	paths, err = cache.Ancestors("c.txt", "e.txt")
	is.NoErr(err)
	is.Equal(paths, []string{"a.txt", "b.txt"})
	paths, err = cache.Ancestors("f.txt", "e.txt")
	is.NoErr(err)
	is.Equal(paths, []string{"a.txt", "b.txt", "c.txt", "d.txt"})
}

func TestDelete(t *testing.T) {
	is := is.New(t)
	fsys := virt.Map{}
	cache, err := dag.Load(fsys, dbPath)
	is.NoErr(err)
	defer cache.Close()
	seed(is, cache)
	// Changing c.txt deletes c.txt and a.txt
	err = cache.Delete("c.txt")
	is.NoErr(err)

	file, err := cache.Get("a.txt")
	is.True(errors.Is(err, dag.ErrNotFound))
	is.Equal(file, nil)
	file, err = cache.Get("c.txt")
	is.True(errors.Is(err, dag.ErrNotFound))
	is.Equal(file, nil)

	files, err := cache.Files()
	is.NoErr(err)
	is.Equal(len(files), 4)
	is.Equal(files[0].Path, "b.txt")
	is.Equal(files[1].Path, "d.txt")
	is.Equal(files[2].Path, "e.txt")
	is.Equal(files[3].Path, "f.txt")

	links, err := cache.Links()
	is.NoErr(err)
	is.Equal(len(links), 2)
	is.Equal(links[0].From, "b.txt")
	is.Equal(links[0].To, "e.txt")
	is.Equal(links[1].From, "d.txt")
	is.Equal(links[1].To, "f.txt")
}

func TestPrint(t *testing.T) {
	is := is.New(t)
	fsys := virt.Map{}
	cache, err := dag.Load(fsys, dbPath)
	is.NoErr(err)
	defer cache.Close()
	seed(is, cache)
	dot := new(bytes.Buffer)
	err = cache.Print(dot)
	is.NoErr(err)
	is.Equal(dot.String(), `digraph dag {
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
	fsys := virt.Map{}
	cache, err := dag.Load(fsys, dbPath)
	is.NoErr(err)
	defer cache.Close()
	for i := 0; i < 100; i++ {
		is.NoErr(cache.Set(fmt.Sprintf("%d.txt", i), &virt.File{
			Data: []byte(fmt.Sprintf("%d.txt", i)),
		}))
	}
}

func TestConcurrentWrites(t *testing.T) {
	is := is.New(t)
	fsys := virt.Map{}
	parent := func(t testing.TB, cmd *exec.Cmd) {
		dbPath := filepath.Join(t.TempDir(), "test_concurrent_writes.db")
		cache, err := dag.Load(fsys, dbPath)
		is.NoErr(err)
		defer cache.Close()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(cmd.Env, "DBPATH="+dbPath)
		is.NoErr(cmd.Start())
		for i := 0; i < 100; i++ {
			is.NoErr(cache.Set(fmt.Sprintf("%d.txt", i), &virt.File{
				Data: []byte(fmt.Sprintf("%d.txt", i)),
			}))
		}
		is.NoErr(cmd.Wait())
	}
	child := func(t testing.TB) {
		dbPath := os.Getenv("DBPATH")
		cache, err := dag.Load(fsys, dbPath)
		is.NoErr(err)
		defer cache.Close()
		for i := 100; i < 200; i++ {
			is.NoErr(cache.Set(fmt.Sprintf("%d.txt", i), &virt.File{
				Data: []byte(fmt.Sprintf("%d.txt", i)),
			}))
		}
	}
	testsub.Run(t, parent, child)
}
