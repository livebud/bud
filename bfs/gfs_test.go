package bfs_test

import (
	"errors"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/bfs"
)

func View() map[string]bfs.Generator {
	return map[string]bfs.Generator{
		"duo/view": bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
			dir.Entry("index.svelte", bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
				file.Write([]byte(`<h1>index</h1>`))
				return nil
			}))
			dir.Entry("about/about.svelte", bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
				file.Write([]byte(`<h2>about</h2>`))
				return nil
			}))
			return nil
		}),
	}
}

func TestFS(t *testing.T) {
	is := is.New(t)

	// 1. duo
	df := bfs.New(nil)
	df.Add(View())
	file, err := df.Open("duo")
	is.NoErr(err)
	rdf, ok := file.(fs.ReadDirFile)
	is.True(ok)
	des, err := rdf.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err := des[0].Info()
	is.NoErr(err)
	is.Equal(fi.IsDir(), true)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.ModeDir)
	is.Equal(fi.Name(), "view")
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	stat, err := file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "duo")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// Stat duo
	df = bfs.New(nil)
	df.Add(View())
	stat, err = fs.Stat(df, "duo")
	is.NoErr(err)
	is.Equal(stat.Name(), "duo")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir duo
	df = bfs.New(nil)
	df.Add(View())
	des, err = fs.ReadDir(df, "duo")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)

	// 2. duo/view
	df = bfs.New(nil)
	df.Add(View())
	file, err = df.Open("duo/view")
	is.NoErr(err)
	rdf, ok = file.(fs.ReadDirFile)
	is.True(ok)
	des, err = rdf.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "about")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "about")
	is.Equal(fi.IsDir(), true)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.ModeDir)
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	is.Equal(des[1].Name(), "index.svelte")
	is.Equal(des[1].IsDir(), false)
	is.Equal(des[1].Type(), fs.FileMode(0))
	fi, err = des[1].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "index.svelte")
	is.Equal(fi.IsDir(), false)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0))
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "view")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// Stat duo
	df = bfs.New(nil)
	df.Add(View())
	stat, err = fs.Stat(df, "duo/view")
	is.NoErr(err)
	is.Equal(stat.Name(), "view")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir duo
	df = bfs.New(nil)
	df.Add(View())
	des, err = fs.ReadDir(df, "duo/view")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "about")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "about")
	is.Equal(fi.IsDir(), true)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.ModeDir)
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	is.Equal(des[1].Name(), "index.svelte")
	is.Equal(des[1].IsDir(), false)
	is.Equal(des[1].Type(), fs.FileMode(0))
	fi, err = des[1].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "index.svelte")
	is.Equal(fi.IsDir(), false)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0))
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)

	// 3. duo/view/about
	df = bfs.New(nil)
	df.Add(View())
	file, err = df.Open("duo/view/about")
	is.NoErr(err)
	rdf, ok = file.(fs.ReadDirFile)
	is.True(ok)
	des, err = rdf.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "about.svelte")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[0].Type(), fs.FileMode(0))
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "about.svelte")
	is.Equal(fi.IsDir(), false)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0))
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "about")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// Stat duo
	df = bfs.New(nil)
	df.Add(View())
	stat, err = fs.Stat(df, "duo/view/about")
	is.NoErr(err)
	is.Equal(stat.Name(), "about")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir duo
	df = bfs.New(nil)
	df.Add(View())
	des, err = fs.ReadDir(df, "duo/view/about")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "about.svelte")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[0].Type(), fs.FileMode(0))
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "about.svelte")
	is.Equal(fi.IsDir(), false)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0))
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)

	// 4. duo/view/index.svelte
	// Open
	df = bfs.New(nil)
	df.Add(View())
	file, err = df.Open("duo/view/index.svelte")
	is.NoErr(err)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "index.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// Stat
	df = bfs.New(nil)
	df.Add(View())
	stat, err = fs.Stat(df, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(stat.Name(), "index.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// ReadFile
	df = bfs.New(nil)
	df.Add(View())
	code, err := fs.ReadFile(df, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// 4. duo/view/about/about.svelte
	// Open
	df = bfs.New(nil)
	df.Add(View())
	file, err = df.Open("duo/view/about/about.svelte")
	is.NoErr(err)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "about.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// Stat
	df = bfs.New(nil)
	df.Add(View())
	stat, err = fs.Stat(df, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(stat.Name(), "about.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// ReadFile
	df = bfs.New(nil)
	df.Add(View())
	code, err = fs.ReadFile(df, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h2>about</h2>`)
}

func TestInstance(t *testing.T) {
	is := is.New(t)
	df := bfs.New(nil)
	df.Add(View())

	// 1. duo
	file, err := df.Open("duo")
	is.NoErr(err)
	rdf, ok := file.(fs.ReadDirFile)
	is.True(ok)
	des, err := rdf.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err := des[0].Info()
	is.NoErr(err)
	is.Equal(fi.IsDir(), true)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.ModeDir)
	is.Equal(fi.Name(), "view")
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	stat, err := file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "duo")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// Stat duo
	stat, err = fs.Stat(df, "duo")
	is.NoErr(err)
	is.Equal(stat.Name(), "duo")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir duo
	des, err = fs.ReadDir(df, "duo")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)

	// 2. duo/view
	file, err = df.Open("duo/view")
	is.NoErr(err)
	rdf, ok = file.(fs.ReadDirFile)
	is.True(ok)
	des, err = rdf.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "about")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "about")
	is.Equal(fi.IsDir(), true)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.ModeDir)
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	is.Equal(des[1].Name(), "index.svelte")
	is.Equal(des[1].IsDir(), false)
	is.Equal(des[1].Type(), fs.FileMode(0))
	fi, err = des[1].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "index.svelte")
	is.Equal(fi.IsDir(), false)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0))
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "view")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// Stat duo
	stat, err = fs.Stat(df, "duo/view")
	is.NoErr(err)
	is.Equal(stat.Name(), "view")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir duo
	des, err = fs.ReadDir(df, "duo/view")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "about")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "about")
	is.Equal(fi.IsDir(), true)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.ModeDir)
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	is.Equal(des[1].Name(), "index.svelte")
	is.Equal(des[1].IsDir(), false)
	is.Equal(des[1].Type(), fs.FileMode(0))
	fi, err = des[1].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "index.svelte")
	is.Equal(fi.IsDir(), false)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0))
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)

	// 3. duo/view/about
	file, err = df.Open("duo/view/about")
	is.NoErr(err)
	rdf, ok = file.(fs.ReadDirFile)
	is.True(ok)
	des, err = rdf.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "about.svelte")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[0].Type(), fs.FileMode(0))
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "about.svelte")
	is.Equal(fi.IsDir(), false)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0))
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "about")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// Stat duo
	stat, err = fs.Stat(df, "duo/view/about")
	is.NoErr(err)
	is.Equal(stat.Name(), "about")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir duo
	des, err = fs.ReadDir(df, "duo/view/about")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "about.svelte")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[0].Type(), fs.FileMode(0))
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "about.svelte")
	is.Equal(fi.IsDir(), false)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0))
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)

	// 4. duo/view/index.svelte
	// Open
	file, err = df.Open("duo/view/index.svelte")
	is.NoErr(err)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "index.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// Stat
	stat, err = fs.Stat(df, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(stat.Name(), "index.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// ReadFile
	code, err := fs.ReadFile(df, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// 4. duo/view/about/about.svelte
	// Open
	file, err = df.Open("duo/view/about/about.svelte")
	is.NoErr(err)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "about.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// Stat
	stat, err = fs.Stat(df, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(stat.Name(), "about.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// ReadFile
	code, err = fs.ReadFile(df, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h2>about</h2>`)
}

func TestDirFS(t *testing.T) {
	is := is.New(t)
	df := bfs.New(os.DirFS("."))
	df.Add(map[string]bfs.Generator{
		"duo/dfs/bfs.go": bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
			if err := bfs.Exists(f, "bfs.go"); err != nil {
				return err
			}
			file.Write([]byte(`package bfs`))
			return nil
		}),
		"duo/public/public.go": bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
			if err := bfs.Exists(f, "public/public.go"); err != nil {
				return err
			}
			file.Write([]byte(`package public`))
			return nil
		}),
	})
	code, err := fs.ReadFile(df, "duo/dfs/bfs.go")
	is.NoErr(err)
	is.Equal(string(code), `package bfs`)
	code, err = fs.ReadFile(df, "duo/public/public.go")
	is.Equal(err.Error(), "open duo/public/public.go > exists public/public.go > open public/public.go > file does not exist")
	is.Equal(code, nil)
	stat, err := fs.Stat(df, "bfs.go")
	is.NoErr(err)
	is.Equal(stat.Name(), "bfs.go")
	is.Equal(stat.IsDir(), false)
}
func TestGenerateFileError(t *testing.T) {
	is := is.New(t)
	df := bfs.New(os.DirFS("."))
	df.Add(map[string]bfs.Generator{
		"duo/main.go": bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
			return fs.ErrNotExist
		}),
	})
	code, err := fs.ReadFile(df, "duo/main.go")
	is.True(err != nil)
	is.Equal(err.Error(), "open duo/main.go > file does not exist")
	is.Equal(code, nil)
}
func TestFileUnderneath(t *testing.T) {
	is := is.New(t)
	df := bfs.New(os.DirFS("."))
	df.Add(map[string]bfs.Generator{
		"bfs.go": bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
			return fs.ErrNotExist
		}),
	})
	stat, err := fs.Stat(df, "bfs.go")
	is.True(err != nil)
	is.Equal(err.Error(), "open bfs.go > file does not exist")
	is.Equal(stat, nil)
}
func TestServeFile(t *testing.T) {
	is := is.New(t)
	df := bfs.New(os.DirFS("."))
	df.Add(map[string]bfs.Generator{
		"duo/view": bfs.ServeFile(func(f bfs.FS, file *bfs.File) error {
			file.Write([]byte(file.Path() + `'s data`))
			return nil
		}),
	})
	des, err := fs.ReadDir(df, "duo/view")
	is.True(errors.Is(err, fs.ErrInvalid))
	is.Equal(len(des), 0)

	// _index.svelte
	file, err := df.Open("duo/view/_index.svelte")
	is.NoErr(err)
	stat, err := file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "_index.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(29))
	is.Equal(stat.Sys(), nil)
	code, err := fs.ReadFile(df, "duo/view/_index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `duo/view/_index.svelte's data`)

	// about/_about.svelte
	file, err = df.Open("duo/view/about/_about.svelte")
	is.NoErr(err)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "_about.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(35))
	is.Equal(stat.Sys(), nil)
	code, err = fs.ReadFile(df, "duo/view/about/_about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `duo/view/about/_about.svelte's data`)
}
func TestHTTP(t *testing.T) {
	is := is.New(t)
	df := bfs.New(os.DirFS("."))
	df.Add(map[string]bfs.Generator{
		"duo/view": bfs.ServeFile(func(f bfs.FS, file *bfs.File) error {
			file.Write([]byte(file.Path() + `'s data`))
			return nil
		}),
	})
	hfs := http.FS(df)

	handler := func(w http.ResponseWriter, r *http.Request) {
		file, err := hfs.Open(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		stat, err := file.Stat()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Add("Content-Type", "text/javascript")
		http.ServeContent(w, r, r.URL.Path, stat.ModTime(), file)
	}

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/duo/view/_index.svelte", nil)
	handler(w, r)

	response := w.Result()
	body, err := ioutil.ReadAll(response.Body)
	is.NoErr(err)
	is.Equal(string(body), `duo/view/_index.svelte's data`)
	is.Equal(response.StatusCode, 200)
}

func TestDirects(t *testing.T) {
	t.SkipNow()
	// is := is.New(t)
	// dfs := bfs.New(nil, &bfs.Dir{
	// 	Entries: map[string]bfs.Generator{
	// 		"duo/view": &bfs.Dir{
	// 			Entries: map[string]bfs.Generator{
	// 				"index.svelte": &bfs.File{
	// 					Data: []byte(`<h1>index</h1>`),
	// 				},
	// 			},
	// 		},
	// 	},
	// })
	// code, err := fs.ReadFile(dfs, "duo/view/index.svelte")
	// is.NoErr(err)
	// is.Equal(string(code), "<h1>index</h1>")
}

func rootless(fpath string) string {
	parts := strings.Split(fpath, string(filepath.Separator))
	return path.Join(parts[1:]...)
}

func TestTargetPath(t *testing.T) {
	is := is.New(t)
	// Test inner file and rootless
	df := bfs.New(nil)
	df.Add(map[string]bfs.Generator{
		"duo/view": bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
			is.Equal(dir.Path(), "duo/view")
			dir.Entry("about/about.svelte", bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
				is.Equal(file.Path(), "duo/view/about/about.svelte")
				file.Write([]byte(rootless(file.Path())))
				return nil
			}))
			return nil
		}),
	})
	code, err := fs.ReadFile(df, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), "view/about/about.svelte")
}

func TestDynamicDir(t *testing.T) {
	is := is.New(t)
	df := bfs.New(nil)
	df.Add(map[string]bfs.Generator{
		"duo/view": bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
			doms := []string{"about/about.svelte", "index.svelte"}
			for _, dom := range doms {
				dom := dom
				dir.Entry(dom, bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
					file.Write([]byte(`<h1>` + dom + `</h1>`))
					return nil
				}))
			}
			return nil
		}),
	})
	des, err := fs.ReadDir(df, "duo/view")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "about")
	is.Equal(des[1].Name(), "index.svelte")
	des, err = fs.ReadDir(df, "duo/view/about")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "about.svelte")
}

func TestBases(t *testing.T) {
	is := is.New(t)
	df := bfs.New(nil)
	df.Add(map[string]bfs.Generator{
		"duo/view": bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
			return nil
		}),
		"duo/controller": bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
			return nil
		}),
	})
	stat, err := fs.Stat(df, "duo/controller")
	is.NoErr(err)
	is.Equal(stat.Name(), "controller")
	stat, err = fs.Stat(df, "duo/view")
	is.NoErr(err)
	is.Equal(stat.Name(), "view")
}

func TestDirPath(t *testing.T) {
	is := is.New(t)
	df := bfs.New(nil)
	df.Add(map[string]bfs.Generator{
		"duo/view": bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
			is.Equal(dir.Path(), "duo/view")
			dir.Entry("public", bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
				is.Equal(dir.Path(), "duo/view/public")
				dir.Entry("favicon.ico", bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
					is.Equal(file.Path(), "duo/view/public/favicon.ico")
					file.Write([]byte("cool_favicon.ico"))
					return nil
				}))
				return nil
			}))
			return nil
		}),
		"duo": bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
			is.Equal(dir.Path(), "duo")
			dir.Entry("controller", bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
				is.Equal(dir.Path(), "duo/controller")
				dir.Entry("controller.go", bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
					is.Equal(file.Path(), "duo/controller/controller.go")
					file.Write([]byte("package controller"))
					return nil
				}))
				return nil
			}))
			return nil
		}),
	})
	code, err := fs.ReadFile(df, "duo/view/public/favicon.ico")
	is.NoErr(err)
	is.Equal(string(code), "cool_favicon.ico")
	code, err = fs.ReadFile(df, "duo/controller/controller.go")
	is.NoErr(err)
	is.Equal(string(code), "package controller")
}

func TestWatchFile(t *testing.T) {
	is := is.New(t)
	df := bfs.New(os.DirFS("."))
	df.Add(map[string]bfs.Generator{
		"duo/view": bfs.ServeFile(func(f bfs.FS, file *bfs.File) error {
			source := rootless(file.Path())
			file.Write([]byte(source + `'s data`))
			file.Watch(source, bfs.WriteEvent|bfs.RemoveEvent)
			return nil
		}),
	})
	subs, err := df.Subscribe("duo/view/index.svelte")
	is.NoErr(err)
	select {
	case <-subs.Wait():
		t.Fatal("No event expected")
	default:
	}
	df.Trigger("view/index.svelte", bfs.WriteEvent)
	select {
	default:
		t.Fatal("Write event expected")
	case event := <-subs.Wait():
		is.Equal(string(event), "Write")
	}
}
func TestWatchDir(t *testing.T) {
	is := is.New(t)
	df := bfs.New(os.DirFS("."))
	df.Add(map[string]bfs.Generator{
		"duo/view": bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
			dir.Entry("index.svelte", bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
				file.Write([]byte(`<h2>index</h2>`))
				return nil
			}))
			dir.Entry("about/about.svelte", bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
				file.Write([]byte(`<h2>about</h2>`))
				return nil
			}))
			dir.Watch("view/{**,*}.{svelte,jsx}", bfs.CreateEvent|bfs.RemoveEvent)
			return nil
		}),
	})
	subs, err := df.Subscribe("duo/view")
	is.NoErr(err)
	select {
	case <-subs.Wait():
		t.Fatal("No event expected")
	default:
	}
	df.Trigger("view/edit.svelte", bfs.CreateEvent)
	select {
	default:
		t.Fatal("Write event expected")
	case event := <-subs.Wait():
		is.Equal(string(event), "Create")
	}
}
func TestDirMerge(t *testing.T) {
	is := is.New(t)
	df := bfs.New(nil)
	df.Add(map[string]bfs.Generator{
		"duo/view": bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
			dir.Entry("index.svelte", bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
				file.Write([]byte(`<h1>index</h1>`))
				return nil
			}))
			dir.Entry("somedir", bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
				return nil
			}))
			return nil
		}),
		"duo/view/view.go": bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
			file.Write([]byte(`package view`))
			return nil
		}),
		"duo/view/plugin.go": bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
			file.Write([]byte(`package plugin`))
			return nil
		}),
	})

	// duo/view
	des, err := fs.ReadDir(df, "duo/view")
	is.NoErr(err)
	is.Equal(len(des), 4)
	is.Equal(des[0].Name(), "index.svelte")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[1].Name(), "plugin.go")
	is.Equal(des[1].IsDir(), false)
	is.Equal(des[2].Name(), "somedir")
	is.Equal(des[2].IsDir(), true)
	is.Equal(des[3].Name(), "view.go")
	is.Equal(des[3].IsDir(), false)
}

func TestAddGenerator(t *testing.T) {
	is := is.New(t)
	// Add the view
	df := bfs.New(nil)
	df.Add(View())

	// Add the controller
	df.Add(map[string]bfs.Generator{
		"duo/controller": bfs.GenerateDir(func(f bfs.FS, dir *bfs.Dir) error {
			dir.Entry("controller.go", bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
				file.Write([]byte(`package controller`))
				return nil
			}))
			return nil
		}),
	})

	des, err := fs.ReadDir(df, "duo")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "controller")
	is.Equal(des[1].Name(), "view")

	// Read from view
	code, err := fs.ReadFile(df, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// Read from controller
	code, err = fs.ReadFile(df, "duo/controller/controller.go")
	is.NoErr(err)
	is.Equal(string(code), `package controller`)
}

type commandGenerator struct {
	Input string
}

func (c *commandGenerator) GenerateFile(f bfs.FS, file *bfs.File) error {
	file.Write([]byte(c.Input + c.Input))
	return nil
}

func TestFileGenerator(t *testing.T) {
	is := is.New(t)
	// Add the view
	df := bfs.New(nil)
	df.Add(map[string]bfs.Generator{
		"duo/command/command.go": bfs.FileGenerator(&commandGenerator{Input: "a"}),
	})
	code, err := fs.ReadFile(df, "duo/command/command.go")
	is.NoErr(err)
	is.Equal(string(code), "aa")
}

type viewGenerator struct {
	Input string
}

func (c *viewGenerator) GenerateDir(f bfs.FS, dir *bfs.Dir) error {
	dir.Entry("index.svelte", bfs.GenerateFile(func(f bfs.FS, file *bfs.File) error {
		file.Write([]byte(c.Input + c.Input))
		return nil
	}))
	return nil
}

func (c *viewGenerator) ServeFile(f bfs.FS, file *bfs.File) error {
	file.Write([]byte(c.Input + c.Input))
	return nil
}

func TestDirGenerator(t *testing.T) {
	is := is.New(t)
	// Add the view
	df := bfs.New(nil)
	df.Add(map[string]bfs.Generator{
		"duo/view": bfs.DirGenerator(&viewGenerator{Input: "a"}),
	})
	code, err := fs.ReadFile(df, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), "aa")
}

func TestFileServer(t *testing.T) {
	is := is.New(t)
	// Add the view
	df := bfs.New(nil)
	df.Add(map[string]bfs.Generator{
		"duo/view": bfs.FileServer(&viewGenerator{Input: "a"}),
	})
	code, err := fs.ReadFile(df, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), "aa")
}

func TestServeDir(t *testing.T) {
	// TODO: ServeDir test
}
