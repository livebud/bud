package budfs_test

import (
	"bytes"
	"context"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/remotefs"
	"github.com/livebud/bud/package/vfs"
	"github.com/livebud/bud/package/virtual/vcache"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/testsub"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/log/testlog"
)

func TestReadFsys(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["a.txt"] = "a"
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	code, err := fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "a")
}

func TestGeneratorPriority(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["a.txt"] = "a"
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	bfs.GenerateFile("a.txt", func(fsys *budfs.FS, file *budfs.File) error {
		file.Data = []byte("b")
		return nil
	})
	code, err := fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "b")
}

func TestCaching(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	count := 1
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	bfs.GenerateFile("a.txt", func(fsys *budfs.FS, file *budfs.File) error {
		file.Data = []byte(strconv.Itoa(count))
		count++
		return nil
	})
	code, err := fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	code, err = fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
}

func TestCachingThroughDir(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	dirCount := 0
	fileCount := 0
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	bfs.GenerateDir("bud/public", func(fsys *budfs.FS, dir *budfs.Dir) error {
		dirCount++
		dir.GenerateFile("public.go", func(fsys *budfs.FS, file *budfs.File) error {
			fileCount++
			file.Data = []byte("public")
			return nil
		})
		return nil
	})
	code, err := fs.ReadFile(bfs, "bud/public/public.go")
	is.NoErr(err)
	is.Equal(string(code), "public")
	des, err := fs.ReadDir(bfs, "bud/public")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(fileCount, 1)
	is.Equal(dirCount, 1)
	// Try calling again to make sure it doesn't regenerate
	code, err = fs.ReadFile(bfs, "bud/public/public.go")
	is.NoErr(err)
	is.Equal(string(code), "public")
	des, err = fs.ReadDir(bfs, "bud/public")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(fileCount, 1)
	is.Equal(dirCount, 1)
	// Update and try again
	bfs.Update("bud/public")
	code, err = fs.ReadFile(bfs, "bud/public/public.go")
	is.NoErr(err)
	is.Equal(string(code), "public")
	des, err = fs.ReadDir(bfs, "bud/public")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(fileCount, 1)
	is.Equal(dirCount, 2)
}

func TestCachingThroughNestedDir(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	dirCount := 0
	fileCount := 0
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	bfs.GenerateDir("bud", func(fsys *budfs.FS, dir *budfs.Dir) error {
		dir.GenerateDir("public", func(fsys *budfs.FS, dir *budfs.Dir) error {
			dirCount++
			dir.GenerateFile("public.go", func(fsys *budfs.FS, file *budfs.File) error {
				fileCount++
				file.Data = []byte("public")
				return nil
			})
			return nil
		})
		return nil
	})
	code, err := fs.ReadFile(bfs, "bud/public/public.go")
	is.NoErr(err)
	is.Equal(string(code), "public")
	des, err := fs.ReadDir(bfs, "bud/public")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(fileCount, 1)
	is.Equal(dirCount, 1)
}

func TestFileFSUpdate(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	count := 1
	td := testdir.New(dir)
	td.Files["a.txt"] = strconv.Itoa(count)
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	bfs.GenerateFile("b.txt", func(fsys *budfs.FS, file *budfs.File) error {
		code, err := fs.ReadFile(fsys, "a.txt")
		if err != nil {
			return err
		}
		file.Data = []byte(code)
		return nil
	})
	code, err := fs.ReadFile(bfs, "b.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	// Update the count
	count++
	is.NoErr(os.WriteFile(filepath.Join(dir, "a.txt"), []byte(strconv.Itoa(count)), 0644))
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

func TestUpdateGenDep(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	count := 1
	bfs.GenerateFile("a.txt", func(fsys *budfs.FS, file *budfs.File) error {
		file.Data = []byte(strconv.Itoa(count))
		count++
		return nil
	})
	bfs.GenerateFile("b.txt", func(fsys *budfs.FS, file *budfs.File) error {
		code, err := fs.ReadFile(fsys, "a.txt")
		if err != nil {
			return err
		}
		file.Data = []byte(code)
		return nil
	})
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
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["a.txt"] = "a"
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	des, err := fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 3) // includes go.mod and package.json
	// Create a new file
	is.NoErr(os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b"), 0644))
	// Cached
	des, err = fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 3)
	// Create the file
	bfs.Create("b.txt")
	// Try again
	des, err = fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 4)
}

func TestDirGenCreate(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["docs/a.txt"] = "a"
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	bfs.GenerateFile("bud/docs.txt", func(fsys *budfs.FS, file *budfs.File) error {
		des, err := fs.ReadDir(fsys, "docs")
		if err != nil {
			return err
		}
		file.Data = []byte(strconv.Itoa(len(des)))
		return nil
	})
	code, err := fs.ReadFile(bfs, "bud/docs.txt")
	is.NoErr(err)
	is.Equal(string(code), "1")
	// Add a file
	is.NoErr(os.WriteFile(filepath.Join(dir, "docs/b.txt"), []byte("b"), 0644))
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
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["a.txt"] = "a"
	td.Files["b.txt"] = "b"
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	des, err := fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 4)
	// Remove file
	is.NoErr(os.RemoveAll(filepath.Join(dir, "b.txt")))
	// Cached
	des, err = fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 4)
	// Create the file
	bfs.Delete("b.txt")
	// Try again
	des, err = fs.ReadDir(bfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 3)
}

func TestDirGenDelete(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["docs/a.txt"] = "a"
	td.Files["docs/b.txt"] = "b"
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	bfs.GenerateFile("bud/docs.txt", func(fsys *budfs.FS, file *budfs.File) error {
		des, err := fs.ReadDir(fsys, "docs")
		if err != nil {
			return err
		}
		file.Data = []byte(strconv.Itoa(len(des)))
		return nil
	})
	code, err := fs.ReadFile(bfs, "bud/docs.txt")
	is.NoErr(err)
	is.Equal(string(code), "2")
	// Add a file
	is.NoErr(os.RemoveAll(filepath.Join(dir, "docs/b.txt")))
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
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	tailwind := vfs.Memory{
		"tailwind/tailwind.css":  &fstest.MapFile{Data: []byte("/** tailwind **/")},
		"tailwind/preflight.css": &fstest.MapFile{Data: []byte("/** preflight **/")},
	}
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

// func TestServer(t *testing.T) {
// 	is := is.New(t)
// 	log := testlog.New()
// 	fsys := vfs.Memory{
// 		"view/index.svelte": &fstest.MapFile{Data: []byte("<h1>index</h1>")},
// 	}
// 	bfs := budfs.New(fsys, log)
// 	// Request the file
// 	r := httptest.NewRequest("GET", "/view/index.svelte", nil)
// 	w := httptest.NewRecorder()
// 	bfs.ServeHTTP(w, r)
// 	res := w.Result()
// 	is.Equal(res.StatusCode, 200)
// 	code, err := io.ReadAll(res.Body)
// 	is.NoErr(err)
// 	is.Equal(string(code), "<h1>index</h1>")
// 	// Change the file
// 	fsys["view/index.svelte"] = &fstest.MapFile{Data: []byte("<h1>index!</h1>")}
// 	// Request the file again
// 	r = httptest.NewRequest("GET", "/view/index.svelte", nil)
// 	w = httptest.NewRecorder()
// 	bfs.ServeHTTP(w, r)
// 	res = w.Result()
// 	is.Equal(res.StatusCode, 200)
// 	code, err = io.ReadAll(res.Body)
// 	is.NoErr(err)
// 	is.Equal(string(code), "<h1>index!</h1>")
// }

func TestDefer(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	cache := vcache.New()
	bfs := budfs.New(cache, module, log)
	called := 0
	bfs.GenerateFile("a.txt", func(fsys *budfs.FS, file *budfs.File) error {
		fsys.Defer(func() error {
			called++
			return nil
		})
		file.Data = []byte("b")
		return nil
	})
	code, err := fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "b")
	bfs.Update("a.txt")
	code, err = fs.ReadFile(bfs, "a.txt")
	is.NoErr(err)
	is.Equal(string(code), "b")
	is.Equal(called, 0)
	is.NoErr(bfs.Close())
	is.Equal(called, 2)
	is.NoErr(bfs.Close())
	is.Equal(called, 2)
}

func TestRemoteFS(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		ctx := context.Background()
		is := is.New(t)
		log := testlog.New()
		dir := t.TempDir()
		td := testdir.New(dir)
		err := td.Write(ctx)
		is.NoErr(err)
		module, err := gomod.Find(dir)
		is.NoErr(err)
		cache := vcache.New()
		bfs := budfs.New(cache, module, log)
		count := 1
		bfs.GenerateDir("bud/generator", func(fsys *budfs.FS, dir *budfs.Dir) error {
			dir.GenerateFile(dir.Relative(), func(fsys *budfs.FS, file *budfs.File) error {
				command := remotefs.Command{
					Env:    cmd.Env,
					Stderr: os.Stderr,
					Stdout: os.Stdout,
				}
				remotefs, err := command.Start(ctx, cmd.Path, cmd.Args[1:]...)
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
		})
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
		fsys := vfs.Memory{
			"a.txt": &fstest.MapFile{Data: []byte("a")},
			"b.txt": &fstest.MapFile{Data: []byte("b")},
		}
		err := remotefs.ServeFrom(ctx, fsys, "")
		is.NoErr(err)
	}
	testsub.Run(t, parent, child)
}

func TestMountRemoteFS(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		cache := vcache.New()
		bfs := budfs.New(cache, module, log)
		command := remotefs.Command{
			Env:    cmd.Env,
			Stderr: os.Stderr,
			Stdout: os.Stdout,
		}
		remotefs, err := command.Start(ctx, cmd.Path, cmd.Args[1:]...)
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
		is.Equal(string(code), "bb")
		// Update the file
		bfs.Update("bud/generator/a.txt")
		// Read again
		code, err = fs.ReadFile(bfs, "bud/generator/a.txt")
		is.NoErr(err)
		is.Equal(string(code), "a")
	}
	child := func(t testing.TB) {
		count := 1
		cache := vcache.New()
		bfs := budfs.New(cache, module, log)
		bfs.GenerateFile("a.txt", func(fsys *budfs.FS, file *budfs.File) error {
			file.Data = []byte(strings.Repeat(string("a"), count))
			count++
			return nil
		})
		bfs.GenerateFile("b.txt", func(fsys *budfs.FS, file *budfs.File) error {
			file.Data = []byte(strings.Repeat(string("b"), count))
			count++
			return nil
		})
		err := remotefs.ServeFrom(ctx, bfs, "")
		is.NoErr(err)
	}
	testsub.Run(t, parent, child)
}

type remoteService struct {
	cmd     *exec.Cmd
	process *remotefs.Process
}

func (s *remoteService) GenerateFile(fsys *budfs.FS, file *budfs.File) (err error) {
	// This remote service depends on the generators
	_, err = fs.Glob(fsys, "generator/*/*.go")
	if err != nil {
		return err
	}
	if s.process != nil {
		if err := s.process.Close(); err != nil {
			return err
		}
	}
	command := remotefs.Command{
		Env:    s.cmd.Env,
		Stderr: os.Stderr,
		Stdout: os.Stdout,
	}
	s.process, err = command.Start(fsys.Context(), s.cmd.Path, s.cmd.Args[1:]...)
	if err != nil {
		return err
	}
	fsys.Defer(func() error {
		return s.process.Close()
	})
	file.Data = []byte(s.process.URL())
	return nil
}

func TestRemoteService(t *testing.T) {
	ctx := context.Background()
	is := is.New(t)
	log := testlog.New()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["generator/tailwind/tailwind.go"] = "package tailwind"
	err := td.Write(ctx)
	is.NoErr(err)
	module, err := gomod.Find(dir)
	is.NoErr(err)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		cache := vcache.New()
		bfs := budfs.New(cache, module, log)
		defer bfs.Close()
		bfs.FileGenerator("bud/service/generator.url", &remoteService{cmd: cmd})
		// Return a URL to connect to
		url, err := fs.ReadFile(bfs, "bud/service/generator.url")
		is.NoErr(err)
		// Dial that URL
		client, err := remotefs.Dial(ctx, string(url))
		is.NoErr(err)
		defer client.Close()
		// Read the remote file
		code, err := fs.ReadFile(client, "a.txt")
		is.NoErr(err)
		is.Equal(string(code), "a")
		// Cached
		code, err = fs.ReadFile(client, "a.txt")
		is.NoErr(err)
		is.Equal(string(code), "a")
		// Uncached because it's a new file
		code, err = fs.ReadFile(client, "b.txt")
		is.NoErr(err)
		is.Equal(string(code), "bb")
		// Still cached
		url, err = fs.ReadFile(bfs, "bud/service/generator.url")
		is.NoErr(err)
		code, err = fs.ReadFile(client, "a.txt")
		is.NoErr(err)
		is.Equal(string(code), "a")
		code, err = fs.ReadFile(client, "b.txt")
		is.NoErr(err)
		is.Equal(string(code), "bb")
		// Update a dependency
		bfs.Update("generator/tailwind/tailwind.go")
		// Should lead to the generator service being uncached again
		url2, err := fs.ReadFile(bfs, "bud/service/generator.url")
		is.NoErr(err)
		is.True(!bytes.Equal(url, url2))
		// Dial the new URL
		client2, err := remotefs.Dial(ctx, string(url2))
		is.NoErr(err)
		defer client2.Close()
		// Still cached, even though the remote has been restarted
		code, err = fs.ReadFile(client2, "a.txt")
		is.NoErr(err)
		is.Equal(string(code), "a")
		code, err = fs.ReadFile(client2, "b.txt")
		is.NoErr(err)
		is.Equal(string(code), "bb")
	}
	child := func(t testing.TB) {
		count := 1
		cache := vcache.New()
		bfs := budfs.New(cache, module, log)
		defer bfs.Close()
		bfs.GenerateFile("a.txt", func(fsys *budfs.FS, file *budfs.File) error {
			file.Data = []byte(strings.Repeat(string("a"), count))
			count++
			return nil
		})
		bfs.GenerateFile("b.txt", func(fsys *budfs.FS, file *budfs.File) error {
			file.Data = []byte(strings.Repeat(string("b"), count))
			count++
			return nil
		})
		err := remotefs.ServeFrom(ctx, bfs, "")
		is.NoErr(err)
	}
	testsub.Run(t, parent, child)
}
