package gen_test

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
	"gitlab.com/mnm/bud/gen"
	"gitlab.com/mnm/bud/vfs"
)

func View() map[string]gen.Generator {
	return map[string]gen.Generator{
		"duo/view": gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
			dir.Entry("index.svelte", gen.GenerateFile(func(f gen.F, file *gen.File) error {
				file.Write([]byte(`<h1>index</h1>`))
				return nil
			}))
			dir.Entry("about/about.svelte", gen.GenerateFile(func(f gen.F, file *gen.File) error {
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
	df := gen.New(nil)
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
	df = gen.New(nil)
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
	df = gen.New(nil)
	df.Add(View())
	des, err = fs.ReadDir(df, "duo")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)

	// 2. duo/view
	df = gen.New(nil)
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
	df = gen.New(nil)
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
	df = gen.New(nil)
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
	df = gen.New(nil)
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
	df = gen.New(nil)
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
	df = gen.New(nil)
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
	df = gen.New(nil)
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
	df = gen.New(nil)
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
	df = gen.New(nil)
	df.Add(View())
	code, err := fs.ReadFile(df, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// 4. duo/view/about/about.svelte
	// Open
	df = gen.New(nil)
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
	df = gen.New(nil)
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
	df = gen.New(nil)
	df.Add(View())
	code, err = fs.ReadFile(df, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h2>about</h2>`)
}

func TestInstance(t *testing.T) {
	is := is.New(t)
	df := gen.New(nil)
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
	df := gen.New(os.DirFS("."))
	df.Add(map[string]gen.Generator{
		"duo/dfs/gen.go": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			if _, err := fs.Stat(f, "gen.go"); err != nil {
				return err
			}
			file.Write([]byte(`package gen`))
			return nil
		}),
		"duo/public/public.go": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			if _, err := fs.Stat(f, "public/public.go"); err != nil {
				return err
			}
			file.Write([]byte(`package public`))
			return nil
		}),
	})
	code, err := fs.ReadFile(df, "duo/dfs/gen.go")
	is.NoErr(err)
	is.Equal(string(code), `package gen`)
	code, err = fs.ReadFile(df, "duo/public/public.go")
	is.Equal(err.Error(), "open duo/public/public.go > open public/public.go > file does not exist")
	is.Equal(code, nil)
	stat, err := fs.Stat(df, "gen.go")
	is.NoErr(err)
	is.Equal(stat.Name(), "gen.go")
	is.Equal(stat.IsDir(), false)
}
func TestGenerateFileError(t *testing.T) {
	is := is.New(t)
	df := gen.New(os.DirFS("."))
	df.Add(map[string]gen.Generator{
		"duo/main.go": gen.GenerateFile(func(f gen.F, file *gen.File) error {
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
	df := gen.New(os.DirFS("."))
	df.Add(map[string]gen.Generator{
		"bfs.go": gen.GenerateFile(func(f gen.F, file *gen.File) error {
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
	df := gen.New(os.DirFS("."))
	df.Add(map[string]gen.Generator{
		"duo/view": gen.ServeFile(func(f gen.F, file *gen.File) error {
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
	df := gen.New(os.DirFS("."))
	df.Add(map[string]gen.Generator{
		"duo/view": gen.ServeFile(func(f gen.F, file *gen.File) error {
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
	df := gen.New(nil)
	df.Add(map[string]gen.Generator{
		"duo/view": gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
			is.Equal(dir.Path(), "duo/view")
			dir.Entry("about/about.svelte", gen.GenerateFile(func(f gen.F, file *gen.File) error {
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
	df := gen.New(nil)
	df.Add(map[string]gen.Generator{
		"duo/view": gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
			doms := []string{"about/about.svelte", "index.svelte"}
			for _, dom := range doms {
				dom := dom
				dir.Entry(dom, gen.GenerateFile(func(f gen.F, file *gen.File) error {
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
	df := gen.New(nil)
	df.Add(map[string]gen.Generator{
		"duo/view": gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
			return nil
		}),
		"duo/controller": gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
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
	df := gen.New(nil)
	df.Add(map[string]gen.Generator{
		"duo/view": gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
			is.Equal(dir.Path(), "duo/view")
			dir.Entry("public", gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
				is.Equal(dir.Path(), "duo/view/public")
				dir.Entry("favicon.ico", gen.GenerateFile(func(f gen.F, file *gen.File) error {
					is.Equal(file.Path(), "duo/view/public/favicon.ico")
					file.Write([]byte("cool_favicon.ico"))
					return nil
				}))
				return nil
			}))
			return nil
		}),
		"duo": gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
			is.Equal(dir.Path(), "duo")
			dir.Entry("controller", gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
				is.Equal(dir.Path(), "duo/controller")
				dir.Entry("controller.go", gen.GenerateFile(func(f gen.F, file *gen.File) error {
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
	df := gen.New(os.DirFS("."))
	df.Add(map[string]gen.Generator{
		"duo/view": gen.ServeFile(func(f gen.F, file *gen.File) error {
			source := rootless(file.Path())
			file.Write([]byte(source + `'s data`))
			file.Watch(source, gen.WriteEvent|gen.RemoveEvent)
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
	df.Trigger("view/index.svelte", gen.WriteEvent)
	select {
	default:
		t.Fatal("Write event expected")
	case event := <-subs.Wait():
		is.Equal(string(event), "Write")
	}
}
func TestWatchDir(t *testing.T) {
	is := is.New(t)
	df := gen.New(os.DirFS("."))
	df.Add(map[string]gen.Generator{
		"duo/view": gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
			dir.Entry("index.svelte", gen.GenerateFile(func(f gen.F, file *gen.File) error {
				file.Write([]byte(`<h2>index</h2>`))
				return nil
			}))
			dir.Entry("about/about.svelte", gen.GenerateFile(func(f gen.F, file *gen.File) error {
				file.Write([]byte(`<h2>about</h2>`))
				return nil
			}))
			dir.Watch("view/{**,*}.{svelte,jsx}", gen.CreateEvent|gen.RemoveEvent)
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
	df.Trigger("view/edit.svelte", gen.CreateEvent)
	select {
	default:
		t.Fatal("Write event expected")
	case event := <-subs.Wait():
		is.Equal(string(event), "Create")
	}
}
func TestDirMerge(t *testing.T) {
	is := is.New(t)
	df := gen.New(nil)
	df.Add(map[string]gen.Generator{
		"duo/view": gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
			dir.Entry("index.svelte", gen.GenerateFile(func(f gen.F, file *gen.File) error {
				file.Write([]byte(`<h1>index</h1>`))
				return nil
			}))
			dir.Entry("somedir", gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
				return nil
			}))
			return nil
		}),
		"duo/view/view.go": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte(`package view`))
			return nil
		}),
		"duo/view/plugin.go": gen.GenerateFile(func(f gen.F, file *gen.File) error {
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
	df := gen.New(nil)
	df.Add(View())

	// Add the controller
	df.Add(map[string]gen.Generator{
		"duo/controller": gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
			dir.Entry("controller.go", gen.GenerateFile(func(f gen.F, file *gen.File) error {
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

func (c *commandGenerator) GenerateFile(f gen.F, file *gen.File) error {
	file.Write([]byte(c.Input + c.Input))
	return nil
}

func TestFileGenerator(t *testing.T) {
	is := is.New(t)
	// Add the view
	df := gen.New(nil)
	df.Add(map[string]gen.Generator{
		"duo/command/command.go": gen.FileGenerator(&commandGenerator{Input: "a"}),
	})
	code, err := fs.ReadFile(df, "duo/command/command.go")
	is.NoErr(err)
	is.Equal(string(code), "aa")
}

type viewGenerator struct {
	Input string
}

func (c *viewGenerator) GenerateDir(f gen.F, dir *gen.Dir) error {
	dir.Entry("index.svelte", gen.GenerateFile(func(f gen.F, file *gen.File) error {
		file.Write([]byte(c.Input + c.Input))
		return nil
	}))
	return nil
}

func (c *viewGenerator) ServeFile(f gen.F, file *gen.File) error {
	file.Write([]byte(c.Input + c.Input))
	return nil
}

func TestDirGenerator(t *testing.T) {
	is := is.New(t)
	// Add the view
	df := gen.New(nil)
	df.Add(map[string]gen.Generator{
		"duo/view": gen.DirGenerator(&viewGenerator{Input: "a"}),
	})
	code, err := fs.ReadFile(df, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), "aa")
}

func TestFileServer(t *testing.T) {
	is := is.New(t)
	// Add the view
	df := gen.New(nil)
	df.Add(map[string]gen.Generator{
		"duo/view": gen.FileServer(&viewGenerator{Input: "a"}),
	})
	code, err := fs.ReadFile(df, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), "aa")
}

func TestServeFS(t *testing.T) {
	is := is.New(t)
	df := gen.New(nil)
	dirfs := os.DirFS(".")
	df.Add(map[string]gen.Generator{
		"duo/view": gen.ServeFS(dirfs),
	})
	// duo/view
	edes, err := fs.ReadDir(dirfs, ".")
	is.NoErr(err)
	ades, err := fs.ReadDir(df, "duo/view")
	is.NoErr(err)
	is.Equal(len(edes), len(ades))
	es, err := fs.Stat(dirfs, ".")
	is.NoErr(err)
	as, err := fs.Stat(df, "duo/view")
	is.NoErr(err)
	is.Equal(as.Name(), "view")
	is.Equal(as.IsDir(), true)
	is.Equal(as.ModTime(), es.ModTime())
	is.Equal(as.Mode(), es.Mode())
	is.Equal(as.Size(), es.Size())
	// ./dir.go
	ed, err := fs.ReadFile(dirfs, "dir.go")
	is.NoErr(err)
	ad, err := fs.ReadFile(df, "duo/view/dir.go")
	is.NoErr(err)
	is.Equal(ed, ad)
	es, err = fs.Stat(dirfs, "dir.go")
	is.NoErr(err)
	as, err = fs.Stat(df, "duo/view/dir.go")
	is.NoErr(err)
	is.Equal(as.Name(), "dir.go")
	is.Equal(as.IsDir(), false)
	is.Equal(as.ModTime(), es.ModTime())
	is.Equal(as.Mode(), es.Mode())
	is.Equal(as.Size(), es.Size())
}

func TestDotReadDirEmpty(t *testing.T) {
	is := is.New(t)
	gfs := gen.New(os.DirFS(t.TempDir()))
	gfs.Add(map[string]gen.Generator{
		"bud/generate/main.go": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte("package main"))
			return nil
		}),
		"go.mod": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte("module pkg"))
			return nil
		}),
	})
	des, err := fs.ReadDir(gfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 2)
}

func TestDotReadDirFiles(t *testing.T) {
	is := is.New(t)
	tmp := t.TempDir()
	err := os.WriteFile(filepath.Join(tmp, "a.txt"), []byte("a"), 0644)
	is.NoErr(err)
	err = os.WriteFile(filepath.Join(tmp, "b.txt"), []byte("b"), 0644)
	is.NoErr(err)
	gfs := gen.New(os.DirFS(tmp))
	gfs.Add(map[string]gen.Generator{
		"bud/generate/main.go": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte("package main"))
			return nil
		}),
		"go.mod": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte("module pkg"))
			return nil
		}),
	})
	des, err := fs.ReadDir(gfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 4)
}

func TestReadDirDuplicates(t *testing.T) {
	is := is.New(t)
	fsys := vfs.Map{
		"go.mod": []byte(`module app.com`),
	}
	genfs := gen.New(fsys)
	genfs.Add(map[string]gen.Generator{
		"go.mod": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			file.Write([]byte("module app.cool"))
			return nil
		}),
	})
	des, err := fs.ReadDir(genfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "go.mod")
	code, err := fs.ReadFile(genfs, "go.mod")
	is.NoErr(err)
	is.Equal(string(code), "module app.cool")
}

func TestFileSkip(t *testing.T) {
	is := is.New(t)
	genfs := gen.New(nil)
	genfs.Add(map[string]gen.Generator{
		"go.mod": gen.GenerateFile(func(f gen.F, file *gen.File) error {
			return file.Skip()
		}),
	})
	des, err := fs.ReadDir(genfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "go.mod")
	code, err := fs.ReadFile(genfs, "go.mod")
	is.True(errors.Is(err, fs.ErrNotExist))
	is.True(strings.Contains(err.Error(), `"go.mod"`))
	is.Equal(code, nil)
}

func TestDirSkip(t *testing.T) {
	is := is.New(t)
	genfs := gen.New(nil)
	genfs.Add(map[string]gen.Generator{
		"bud": gen.GenerateDir(func(f gen.F, dir *gen.Dir) error {
			return dir.Skip()
		}),
	})
	des, err := fs.ReadDir(genfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "bud")
	code, err := fs.ReadFile(genfs, "bud")
	is.True(errors.Is(err, fs.ErrNotExist))
	is.True(strings.Contains(err.Error(), `"bud"`))
	is.Equal(code, nil)
}
