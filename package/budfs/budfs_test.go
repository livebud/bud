package budfs_test

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"net/http/httptest"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/testsub"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/budfs"
	"github.com/livebud/bud/package/log/testlog"
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

func TestSubprocess(t *testing.T) {
	is := is.New(t)
	parent := func(t testing.TB, cmd *exec.Cmd) {
		log := testlog.New()
		fsys := fstest.MapFS{}
		bfs := budfs.New(fsys, log)
		// Prep the count
		cmd.Env = append(cmd.Env, "BUD_COUNT=1")
		// Create the generator
		bfs.FileGenerator("index.svelte", budfs.GenerateFile(func(fsys budfs.FS, file *budfs.File) error {
			var stderr bytes.Buffer
			cmd.Stderr = &stderr
			if err := cmd.Run(); err != nil {
				return err
			}
			// Reset for next time
			next := exec.Command(cmd.Path, cmd.Args[1:]...)
			// Update the counter for the next run
			next.Env = append(cmd.Env, "BUD_COUNT=2")
			next.Stdout = cmd.Stdout
			next.Stderr = cmd.Stderr
			next.Stdin = cmd.Stdin
			next.ExtraFiles = cmd.ExtraFiles
			next.Dir = cmd.Dir
			cmd = next
			// Respond with bytes
			file.Data = stderr.Bytes()
			return nil
		}))
		// Read the file, triggering the subprocess to run
		code, err := fs.ReadFile(bfs, "index.svelte")
		is.NoErr(err)
		is.Equal(string(code), "<h1>1</h1>")
		// Read the file, but it should be cached
		code, err = fs.ReadFile(bfs, "index.svelte")
		is.NoErr(err)
		is.Equal(string(code), "<h1>1</h1>")
		// Mark the file as updated
		bfs.Update("index.svelte")
		// Try again
		code, err = fs.ReadFile(bfs, "index.svelte")
		is.NoErr(err)
		is.Equal(string(code), "<h1>2</h1>")
	}
	child := func(t testing.TB) {
		fmt.Fprintf(os.Stderr, "<h1>"+os.Getenv("BUD_COUNT")+"</h1>")
	}
	testsub.Run(t, parent, child)
}
