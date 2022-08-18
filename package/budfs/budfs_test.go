package budfs_test

import (
	"context"
	"io"
	"io/fs"
	"net/http/httptest"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/testsub"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/log/testlog"
	"github.com/livebud/bud/package/remotefs"
)

func TestReadFsys(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	bfs := budfs.New(fsys, log)
	des, err := fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
}

func TestGeneratorPriority(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	bfs := budfs.New(fsys, log)
	bfs.FileGenerator("a.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
		file.Data = []byte("b")
		return nil
	}))
	code, err := fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "b")
}

func TestCaching(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{}
	bfs := budfs.New(fsys, log)
	count := 1
	bfs.FileGenerator("a.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
		file.Data = []byte(strconv.Itoa(count))
		count++
		return nil
	}))
	code, err := fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	code, err = fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
}

func TestFileFSUpdate(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	count := 1
	fsys := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte(strconv.Itoa(count))},
	}
	bfs := budfs.New(fsys, log)
	bfs.FileGenerator("b.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
		code, err := fs.ReadFile(fsys, "a.txt")
		if err != nil {
			return err
		}
		file.Data = []byte(code)
		return nil
	}))
	code, err := fs.ReadFile(bfs, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	// Update the count
	count++
	fsys["a.txt"] = &fstest.MapFile{Data: []byte(strconv.Itoa(count))}
	// Cached
	code, err = fs.ReadFile(bfs, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	// Update the file
	bfs.Update("a.txt")
	// Read again
	code, err = fs.ReadFile(bfs, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), "2")
}

func TestFileGenUpdate(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{}
	bfs := budfs.New(fsys, log)
	count := 1
	bfs.FileGenerator("a.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
		file.Data = []byte(strconv.Itoa(count))
		count++
		return nil
	}))
	bfs.FileGenerator("b.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
		code, err := fs.ReadFile(fsys, "a.txt")
		if err != nil {
			return err
		}
		file.Data = []byte(code)
		return nil
	}))
	code, err := fs.ReadFile(bfs, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	// Cached
	code, err = fs.ReadFile(bfs, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	// Update the file
	bfs.Update("a.txt")
	// Read again
	code, err = fs.ReadFile(bfs, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), "2")
}

func TestDirFSCreate(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	bfs := budfs.New(fsys, log)
	des, err := fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	fsys["b.txt"] = &fstest.MapFile{Data: []byte("b")}
	// Cached
	des, err = fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	// Create the file
	bfs.Create("b.txt")
	// Try again
	des, err = fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 2)
}

func TestDirGenCreate(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{
		"docs/a.txt": &fstest.MapFile{Data: []byte("a")},
	}
	bfs := budfs.New(fsys, log)
	bfs.FileGenerator("bud/docs.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
		des, err := fs.ReadDir(fsys, "docs")
		if err != nil {
			return err
		}
		file.Data = []byte(strconv.Itoa(len(des)))
		return nil
	}))
	code, err := fs.ReadFile(bfs, "bud/docs.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	// Add a file
	fsys["docs/b.txt"] = &fstest.MapFile{Data: []byte("b")}
	// Cached
	code, err = fs.ReadFile(bfs, "bud/docs.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	// Create the file
	bfs.Create("docs/b.txt")
	// Try again
	code, err = fs.ReadFile(bfs, "bud/docs.txt")
	is.NoErr(err)
	is.Equal(string(code), "2")
}

func TestDirFSDelete(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a")},
		"b.txt": &fstest.MapFile{Data: []byte("b")},
	}
	bfs := budfs.New(fsys, log)
	des, err := fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 2)
	delete(fsys, "b.txt")
	// Cached
	des, err = fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 2)
	// Create the file
	bfs.Delete("b.txt")
	// Try again
	des, err = fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
}

func TestDirGenDelete(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{
		"docs/a.txt": &fstest.MapFile{Data: []byte("a")},
		"docs/b.txt": &fstest.MapFile{Data: []byte("b")},
	}
	bfs := budfs.New(fsys, log)
	bfs.FileGenerator("bud/docs.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
		des, err := fs.ReadDir(fsys, "docs")
		if err != nil {
			return err
		}
		file.Data = []byte(strconv.Itoa(len(des)))
		return nil
	}))
	code, err := fs.ReadFile(bfs, "bud/docs.txt")
	is.NoErr(err)
	is.Equal(string(code), "2")
	// Add a file
	delete(fsys, "docs/b.txt")
	// Cached
	code, err = fs.ReadFile(bfs, "bud/docs.txt")
	is.NoErr(err)
	is.Equal(string(code), "2")
	// Create the file
	bfs.Delete("docs/b.txt")
	// Try again
	code, err = fs.ReadFile(bfs, "bud/docs.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
}

func TestMount(t *testing.T) {
	is := is.New(t)
	fsys := fstest.MapFS{}
	log := testlog.New()
	tailwind := fstest.MapFS{
		"tailwind/tailwind.css":  &fstest.MapFile{Data: []byte("/** tailwind **/")},
		"tailwind/preflight.css": &fstest.MapFile{Data: []byte("/** preflight **/")},
	}
	bfs := budfs.New(fsys, log)
	bfs.Mount("generator", tailwind)
	code, err := fs.ReadFile(bfs, "generator/tailwind/tailwind.css")
	is.NoErr(err)
	is.Equal(string(code), "/** tailwind **/")
	// Update the file
	tailwind["tailwind/tailwind.css"] = &fstest.MapFile{Data: []byte("/** css **/")}
	// Read the cached file
	code, err = fs.ReadFile(bfs, "generator/tailwind/tailwind.css")
	is.NoErr(err)
	is.Equal(string(code), "/** tailwind **/")
	// Mark the file as being updated
	bfs.Update("generator/tailwind/tailwind.css")
	// Read the file again
	code, err = fs.ReadFile(bfs, "generator/tailwind/tailwind.css")
	is.NoErr(err)
	is.Equal(string(code), "/** css **/")
}

func TestServer(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{
		"view/index.svelte": &fstest.MapFile{Data: []byte("<h1>index</h1>")},
	}
	bfs := budfs.New(fsys, log)
	// Request the file
	r := httptest.NewRequest("GET", "/view/index.svelte", nil)
	w := httptest.NewRecorder()
	bfs.ServeHTTP(w, r)
	res := w.Result()
	is.Equal(res.StatusCode, 200)
	code, err := io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(code), "<h1>index</h1>")
	// Change the file
	fsys["view/index.svelte"] = &fstest.MapFile{Data: []byte("<h1>index!</h1>")}
	// Request the file again
	r = httptest.NewRequest("GET", "/view/index.svelte", nil)
	w = httptest.NewRecorder()
	bfs.ServeHTTP(w, r)
	res = w.Result()
	is.Equal(res.StatusCode, 200)
	code, err = io.ReadAll(res.Body)
	is.NoErr(err)
	is.Equal(string(code), "<h1>index!</h1>")
}

func TestRemoteFS(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		log := testlog.New()
		fsys := fstest.MapFS{}
		ctx := context.Background()
		bfs := budfs.New(fsys, log)
		count := 1
		bfs.DirGenerator("bud/generator", budfs.GenerateDir(func(fsys budfs.FS, dir *budfs.Dir) error {
			dir.GenerateFile(dir.Relative(), func(fsys budfs.FS, file *budfs.File) error {
				command := remotefs.Command{
					Env:    cmd.Env,
					Stderr: os.Stderr,
					Stdout: os.Stdout,
				}
				remotefs, err := command.Start(ctx, cmd.Path, cmd.Args...)
				if err != nil {
					return err
				}
				defer remotefs.Close()
				data, err := fs.ReadFile(remotefs, dir.Relative())
				if err != nil {
					return err
				}
				file.Data = []byte(strings.Repeat(string(data), count))
				count++
				return nil
			})
			return nil
		}))
		code, err := fs.ReadFile(bfs, "bud/generator/a.txt")
		is.NoErr(err)
		is.Equal(string(code), "a")
		// Cached
		code, err = fs.ReadFile(bfs, "bud/generator/a.txt")
		is.NoErr(err)
		is.Equal(string(code), "a")
		// Read new path (uncached)
		code, err = fs.ReadFile(bfs, "bud/generator/b.txt")
		is.NoErr(err)
		is.Equal(string(code), "bb")
		// Update the file
		bfs.Update("bud/generator/a.txt")
		// Read again
		code, err = fs.ReadFile(bfs, "bud/generator/a.txt")
		is.NoErr(err)
		is.Equal(string(code), "aaa")
	}
	child := func(t testing.TB) {
		ctx := context.Background()
		fsys := fstest.MapFS{
			"a.txt": &fstest.MapFile{Data: []byte("a")},
			"b.txt": &fstest.MapFile{Data: []byte("b")},
		}
		err := remotefs.ServeFrom(ctx, fsys, "")
		is.NoErr(err)
	}
	testsub.Run(t, parent, child)
}

func TestMountRemoteFS(t *testing.T) {
	is := is.New(t)
	log := testlog.New()
	fsys := fstest.MapFS{}
	ctx := context.Background()
	parent := func(t testing.TB, cmd *exec.Cmd) {
		bfs := budfs.New(fsys, log)
		command := remotefs.Command{
			Env:    cmd.Env,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
		remotefs, err := command.Start(ctx, cmd.Path, cmd.Args...)
		is.NoErr(err)
		defer remotefs.Close()
		bfs.Mount("bud/generator", remotefs)
		code, err := fs.ReadFile(bfs, "bud/generator/a.txt")
		is.NoErr(err)
		is.Equal(string(code), "a")
		// Cached
		code, err = fs.ReadFile(bfs, "bud/generator/a.txt")
		is.NoErr(err)
		is.Equal(string(code), "a")
		// Read new path (uncached)
		code, err = fs.ReadFile(bfs, "bud/generator/b.txt")
		is.NoErr(err)
		is.Equal(string(code), "b")
		// Update the file
		bfs.Update("bud/generator/a.txt")
		// Read again
		code, err = fs.ReadFile(bfs, "bud/generator/a.txt")
		is.NoErr(err)
		is.Equal(string(code), "a")
	}
	child := func(t testing.TB) {
		count := 1
		bfs := budfs.New(fsys, log)
		bfs.FileGenerator("a.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
			file.Data = []byte(strings.Repeat(string("a"), count))
			count++
			return nil
		}))
		bfs.FileGenerator("b.txt", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
			file.Data = []byte(strings.Repeat(string("b"), count))
			count++
			return nil
		}))
		err := remotefs.ServeFrom(ctx, bfs, "")
		is.NoErr(err)
	}
	testsub.Run(t, parent, child)
}
