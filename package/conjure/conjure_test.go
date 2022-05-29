package conjure_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"io/fs"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/conjure"
	"github.com/livebud/bud/package/merged"
)

func View() func(dir *conjure.Dir) error {
	return func(dir *conjure.Dir) error {
		dir.GenerateFile("index.svelte", func(file *conjure.File) error {
			file.Data = []byte(`<h1>index</h1>`)
			return nil
		})
		dir.GenerateFile("about/about.svelte", func(file *conjure.File) error {
			file.Data = []byte(`<h2>about</h2>`)
			return nil
		})
		return nil
	}
}

func TestFS(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.GenerateDir("duo/view", View())

	// 1. duo
	file, err := cfs.Open("duo")
	is.NoErr(err)
	rcfs, ok := file.(fs.ReadDirFile)
	is.True(ok)
	des, err := rcfs.ReadDir(-1)
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
	stat, err = fs.Stat(cfs, "duo")
	is.NoErr(err)
	is.Equal(stat.Name(), "duo")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir duo
	des, err = fs.ReadDir(cfs, "duo")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)

	// 2. duo/view
	file, err = cfs.Open("duo/view")
	is.NoErr(err)
	rcfs, ok = file.(fs.ReadDirFile)
	is.True(ok)
	des, err = rcfs.ReadDir(-1)
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
	stat, err = fs.Stat(cfs, "duo/view")
	is.NoErr(err)
	is.Equal(stat.Name(), "view")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir duo
	des, err = fs.ReadDir(cfs, "duo/view")
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
	file, err = cfs.Open("duo/view/about")
	is.NoErr(err)
	rcfs, ok = file.(fs.ReadDirFile)
	is.True(ok)
	des, err = rcfs.ReadDir(-1)
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
	stat, err = fs.Stat(cfs, "duo/view/about")
	is.NoErr(err)
	is.Equal(stat.Name(), "about")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir duo
	des, err = fs.ReadDir(cfs, "duo/view/about")
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
	file, err = cfs.Open("duo/view/index.svelte")
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
	stat, err = fs.Stat(cfs, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(stat.Name(), "index.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// ReadFile
	code, err := fs.ReadFile(cfs, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// 4. duo/view/about/about.svelte
	// Open
	file, err = cfs.Open("duo/view/about/about.svelte")
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
	stat, err = fs.Stat(cfs, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(stat.Name(), "about.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// ReadFile
	code, err = fs.ReadFile(cfs, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h2>about</h2>`)
}

func TestDir(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.GenerateDir("duo/view", func(dir *conjure.Dir) error {
		dir.GenerateDir("about", func(dir *conjure.Dir) error {
			dir.GenerateDir("me", func(dir *conjure.Dir) error {
				return nil
			})
			return nil
		})
		dir.GenerateDir("users/admin", func(dir *conjure.Dir) error {
			return nil
		})
		return nil
	})
	des, err := fs.ReadDir(cfs, "duo")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
	is.Equal(des[0].IsDir(), true)
	des, err = fs.ReadDir(cfs, "duo/view")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "about")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[1].Name(), "users")
	is.Equal(des[1].IsDir(), true)
	des, err = fs.ReadDir(cfs, "duo/view/about")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "me")
	is.Equal(des[0].IsDir(), true)
	des, err = fs.ReadDir(cfs, "duo/view/about/me")
	is.NoErr(err)
	is.Equal(len(des), 0)
	des, err = fs.ReadDir(cfs, "duo/view/users")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "admin")
	is.Equal(des[0].IsDir(), true)
	des, err = fs.ReadDir(cfs, "duo/view/users/admin")
	is.NoErr(err)
	is.Equal(len(des), 0)
}

func TestGenerateFileError(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.GenerateFile("duo/main.go", func(file *conjure.File) error {
		return fs.ErrNotExist
	})
	code, err := fs.ReadFile(cfs, "duo/main.go")
	is.True(err != nil)
	is.Equal(err.Error(), `conjure: generate "duo/main.go". file does not exist`)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(code, nil)
}

func TestServeFile(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.ServeFile("duo/view", func(file *conjure.File) error {
		file.Data = []byte(file.Path() + `'s data`)
		return nil
	})
	des, err := fs.ReadDir(cfs, "duo/view")
	is.True(errors.Is(err, fs.ErrInvalid))
	is.Equal(len(des), 0)

	// _index.svelte
	file, err := cfs.Open("duo/view/_index.svelte")
	is.NoErr(err)
	stat, err := file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "_index.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(29))
	is.Equal(stat.Sys(), nil)
	code, err := fs.ReadFile(cfs, "duo/view/_index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `duo/view/_index.svelte's data`)

	// about/_about.svelte
	file, err = cfs.Open("duo/view/about/_about.svelte")
	is.NoErr(err)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "_about.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(35))
	is.Equal(stat.Sys(), nil)
	code, err = fs.ReadFile(cfs, "duo/view/about/_about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `duo/view/about/_about.svelte's data`)
}

func TestHTTP(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.ServeFile("bud/view", func(file *conjure.File) error {
		file.Data = []byte(file.Path() + `'s data`)
		return nil
	})
	hfs := http.FS(cfs)

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
	r := httptest.NewRequest("GET", "/bud/view/_index.svelte", nil)
	handler(w, r)

	response := w.Result()
	body, err := ioutil.ReadAll(response.Body)
	is.NoErr(err)
	is.Equal(string(body), `bud/view/_index.svelte's data`)
	is.Equal(response.StatusCode, 200)
}

func rootless(fpath string) string {
	parts := strings.Split(fpath, string(filepath.Separator))
	return path.Join(parts[1:]...)
}

func TestTargetPath(t *testing.T) {
	is := is.New(t)
	// Test inner file and rootless
	cfs := conjure.New()
	cfs.GenerateDir("duo/view", func(dir *conjure.Dir) error {
		dir.GenerateFile("about/about.svelte", func(file *conjure.File) error {
			file.Data = []byte(rootless(file.Path()))
			return nil
		})
		return nil
	})
	code, err := fs.ReadFile(cfs, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), "view/about/about.svelte")
}

func TestDynamicDir(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.GenerateDir("duo/view", func(dir *conjure.Dir) error {
		doms := []string{"about/about.svelte", "index.svelte"}
		for _, dom := range doms {
			dom := dom
			dir.GenerateFile(dom, func(file *conjure.File) error {
				file.Data = []byte(`<h1>` + dom + `</h1>`)
				return nil
			})
		}
		return nil
	})
	des, err := fs.ReadDir(cfs, "duo/view")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "about")
	is.Equal(des[1].Name(), "index.svelte")
	des, err = fs.ReadDir(cfs, "duo/view/about")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "about.svelte")
}

func TestBases(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.GenerateDir("duo/view", func(dir *conjure.Dir) error {
		return nil
	})
	cfs.GenerateDir("duo/controller", func(dir *conjure.Dir) error {
		return nil
	})
	stat, err := fs.Stat(cfs, "duo/controller")
	is.NoErr(err)
	is.Equal(stat.Name(), "controller")
	stat, err = fs.Stat(cfs, "duo/view")
	is.NoErr(err)
	is.Equal(stat.Name(), "view")
}

func TestDirPath(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.GenerateDir("duo/view", func(dir *conjure.Dir) error {
		dir.GenerateDir("public", func(dir *conjure.Dir) error {
			dir.GenerateFile("favicon.ico", func(file *conjure.File) error {
				file.Data = []byte("cool_favicon.ico")
				return nil
			})
			return nil
		})
		return nil
	})
	cfs.GenerateDir("duo", func(dir *conjure.Dir) error {
		dir.GenerateDir("controller", func(dir *conjure.Dir) error {
			dir.GenerateFile("controller.go", func(file *conjure.File) error {
				file.Data = []byte("package controller")
				return nil
			})
			return nil
		})
		return nil
	})
	code, err := fs.ReadFile(cfs, "duo/view/public/favicon.ico")
	is.NoErr(err)
	is.Equal(string(code), "cool_favicon.ico")
	code, err = fs.ReadFile(cfs, "duo/controller/controller.go")
	is.NoErr(err)
	is.Equal(string(code), "package controller")
}

func TestDirMerge(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.GenerateDir("duo/view", func(dir *conjure.Dir) error {
		dir.GenerateFile("index.svelte", func(file *conjure.File) error {
			file.Data = []byte(`<h1>index</h1>`)
			return nil
		})
		dir.GenerateDir("somedir", func(dir *conjure.Dir) error {
			return nil
		})
		return nil
	})
	cfs.GenerateFile("duo/view/view.go", func(file *conjure.File) error {
		file.Data = []byte(`package view`)
		return nil
	})
	cfs.GenerateFile("duo/view/plugin.go", func(file *conjure.File) error {
		file.Data = []byte(`package plugin`)
		return nil
	})
	// duo/view
	des, err := fs.ReadDir(cfs, "duo/view")
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
	cfs := conjure.New()
	cfs.GenerateDir("duo/view", View())

	// Add the controller
	cfs.GenerateDir("duo/controller", func(dir *conjure.Dir) error {
		dir.GenerateFile("controller.go", func(file *conjure.File) error {
			file.Data = []byte(`package controller`)
			return nil
		})
		return nil
	})

	des, err := fs.ReadDir(cfs, "duo")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "controller")
	is.Equal(des[1].Name(), "view")

	// Read from view
	code, err := fs.ReadFile(cfs, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// Read from controller
	code, err = fs.ReadFile(cfs, "duo/controller/controller.go")
	is.NoErr(err)
	is.Equal(string(code), `package controller`)
}

type commandGenerator struct {
	Input string
}

func (c *commandGenerator) GenerateFile(file *conjure.File) error {
	file.Data = []byte(c.Input + c.Input)
	return nil
}

func (c *commandGenerator) GenerateDir(dir *conjure.Dir) error {
	dir.GenerateFile("index.svelte", func(file *conjure.File) error {
		file.Data = []byte(c.Input + c.Input)
		return nil
	})
	return nil
}

func (c *commandGenerator) ServeFile(file *conjure.File) error {
	file.Data = []byte(c.Input + "/" + file.Path())
	return nil
}

func TestFileGenerator(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.FileGenerator("duo/command/command.go", &commandGenerator{Input: "a"})
	code, err := fs.ReadFile(cfs, "duo/command/command.go")
	is.NoErr(err)
	is.Equal(string(code), "aa")
}

func TestDirGenerator(t *testing.T) {
	is := is.New(t)
	// Add the view
	cfs := conjure.New()
	cfs.DirGenerator("duo/view", &commandGenerator{Input: "a"})
	code, err := fs.ReadFile(cfs, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), "aa")
}

func TestFileServer(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.FileServer("duo/view", &commandGenerator{Input: "a"})
	code, err := fs.ReadFile(cfs, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), "a/duo/view/index.svelte")
}

func TestDotReadDirEmpty(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.GenerateFile("bud/generate/main.go", func(file *conjure.File) error {
		file.Data = []byte("package main")
		return nil
	})
	cfs.GenerateFile("go.mod", func(file *conjure.File) error {
		file.Data = []byte("module pkg")
		return nil
	})
	des, err := fs.ReadDir(cfs, ".")
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
	cfs := conjure.New()
	mapfs := fstest.MapFS{
		"a.txt": &fstest.MapFile{Data: []byte("a"), Mode: 0644},
		"b.txt": &fstest.MapFile{Data: []byte("b"), Mode: 0644},
	}
	cfs.GenerateFile("bud/generate/main.go", func(file *conjure.File) error {
		file.Data = []byte("package main")
		return nil
	})
	cfs.GenerateFile("go.mod", func(file *conjure.File) error {
		file.Data = []byte("module pkg")
		return nil
	})
	fsys := merged.Merge(cfs, mapfs)
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 4)
}

func TestReadDirDuplicates(t *testing.T) {
	is := is.New(t)
	mapfs := fstest.MapFS{
		"go.mod": &fstest.MapFile{Data: []byte(`module app.com`)},
	}
	cfs := conjure.New()
	cfs.GenerateFile("go.mod", func(file *conjure.File) error {
		file.Data = []byte("module app.cool")
		return nil
	})
	fsys := merged.Merge(cfs, mapfs)
	des, err := fs.ReadDir(fsys, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "go.mod")
	code, err := fs.ReadFile(fsys, "go.mod")
	is.NoErr(err)
	is.Equal(string(code), "module app.cool")
}

func TestEmbedOpen(t *testing.T) {
	is := is.New(t)
	now := time.Now()
	cfs := conjure.New()
	cfs.FileGenerator("duo/view/index.svelte", &conjure.Embed{
		Data:    []byte(`<h1>index</h1>`),
		Mode:    fs.FileMode(0644),
		ModTime: now,
	})
	cfs.FileGenerator("duo/view/about/about.svelte", &conjure.Embed{
		Data:    []byte(`<h1>about</h1>`),
		Mode:    fs.FileMode(0644),
		ModTime: now,
	})
	cfs.FileGenerator("duo/public/favicon.ico", &conjure.Embed{
		Data:    []byte(`favicon.ico`),
		Mode:    fs.FileMode(0644),
		ModTime: now,
	})
	// duo/view/index.svelte
	code, err := fs.ReadFile(cfs, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)
	stat, err := fs.Stat(cfs, "duo/view/index.svelte")
	is.NoErr(err)
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.IsDir(), false)

	// duo/view/about/about.svelte
	code, err = fs.ReadFile(cfs, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>about</h1>`)
	stat, err = fs.Stat(cfs, "duo/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.IsDir(), false)

	// duo/public/favicon.ico
	code, err = fs.ReadFile(cfs, "duo/public/favicon.ico")
	is.NoErr(err)
	is.Equal(string(code), `favicon.ico`)
	stat, err = fs.Stat(cfs, "duo/public/favicon.ico")
	is.NoErr(err)
	is.Equal(stat.ModTime(), now)
	is.Equal(stat.Mode(), fs.FileMode(0644))
	is.Equal(stat.IsDir(), false)

	// duo/public
	// TODO: consider locking this down, though this might be taken care of higher
	// up in the stack.
	des, err := fs.ReadDir(cfs, "duo/public")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "favicon.ico")
}

func TestGoModGoMod(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.GenerateFile("go.mod", func(file *conjure.File) error {
		file.Data = []byte("module app.com\nrequire mod.test/module v1.2.4")
		return nil
	})
	stat, err := fs.Stat(cfs, "go.mod/go.mod")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(stat, nil)
	stat, err = fs.Stat(cfs, "go.mod")
	is.NoErr(err)
	is.Equal(stat.Name(), "go.mod")
}

func TestGoModGoModEmbed(t *testing.T) {
	is := is.New(t)
	cfs := conjure.New()
	cfs.FileGenerator("go.mod", &conjure.Embed{
		Data: []byte("module app.com\nrequire mod.test/module v1.2.4"),
	})
	stat, err := fs.Stat(cfs, "go.mod/go.mod")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(stat, nil)
	stat, err = fs.Stat(cfs, "go.mod")
	is.NoErr(err)
	is.Equal(stat.Name(), "go.mod")
}
