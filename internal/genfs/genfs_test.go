package genfs_test

import (
	"errors"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/livebud/bud/internal/genfs"
	"github.com/livebud/bud/internal/is"
)

type tailwind struct {
}

func (t *tailwind) GenerateFile(file *genfs.File) error {
	file.Data = []byte("/* tailwind */")
	return nil
}

type svelte struct {
}

func (s *svelte) GenerateFile(file *genfs.File) error {
	file.Data = []byte("/* svelte */")
	return nil
}

func TestFS(t *testing.T) {
	is := is.New(t)
	// fsys := fstest.MapFS{
	// 	"bud/public/index.html": &fstest.MapFile{Data: []byte("<h1>hello</h1>")},
	// }
	// log := log.Discard
	genfs := genfs.New()
	genfs.FileGenerator("bud/public/tailwind/tailwind.css", &tailwind{})
	genfs.FileGenerator("bud/view/index.svelte", &svelte{})

	// .
	des, err := fs.ReadDir(genfs, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)

	// bud
	is.Equal(des[0].Name(), "bud")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err := des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Mode(), fs.ModeDir)
	stat, err := fs.Stat(genfs, "bud")
	is.NoErr(err)
	is.Equal(stat.Mode(), fs.ModeDir)

	// bud/public
	des, err = fs.ReadDir(genfs, "bud")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "public")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "public")
	stat, err = fs.Stat(genfs, "bud/public")
	is.NoErr(err)
	is.Equal(stat.Name(), "public")

	// return errors for non-existent files
	_, err = genfs.Open("bud\\public")
	is.True(errors.Is(err, fs.ErrNotExist))

	// bud/public/tailwind
	des, err = fs.ReadDir(genfs, "bud/public/tailwind")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "tailwind.css")
	is.Equal(des[0].IsDir(), false)

	// read generated data
	data, err := fs.ReadFile(genfs, "bud/public/index.html")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.True(data == nil)
	data, err = fs.ReadFile(genfs, "bud/public/tailwind/tailwind.css")
	is.NoErr(err)
	is.Equal(string(data), "/* tailwind */")
	data, err = fs.ReadFile(genfs, "bud/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(data), "/* svelte */")

	// run the TestFS compliance test suite
	is.NoErr(fstest.TestFS(genfs, "bud/public/tailwind/tailwind.css", "bud/view/index.svelte"))
}

func view() func(dir *genfs.Dir) error {
	return func(dir *genfs.Dir) error {
		dir.GenerateFile("index.svelte", func(file *genfs.File) error {
			file.Data = []byte(`<h1>index</h1>`)
			return nil
		})
		dir.GenerateFile("about/about.svelte", func(file *genfs.File) error {
			file.Data = []byte(`<h2>about</h2>`)
			return nil
		})
		return nil
	}
}

func TestViewFS(t *testing.T) {
	is := is.New(t)
	gen := genfs.New()
	gen.GenerateDir("bud/view", view())

	// bud
	des, err := fs.ReadDir(gen, "bud")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err := des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "view")

	// bud/view
	stat, err := fs.Stat(gen, "bud/view")
	is.NoErr(err)
	is.Equal(stat.Name(), "view")
	is.Equal(stat.IsDir(), true)
	is.Equal(stat.Mode(), fs.ModeDir)

	_, err = gen.Open("about")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))

	_, err = gen.Open("bud/view/.")
	is.True(err != nil)
	is.True(errors.Is(err, fs.ErrNotExist))

	code, err := fs.ReadFile(gen, "bud/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), "<h1>index</h1>")
	code, err = fs.ReadFile(gen, "bud/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), "<h2>about</h2>")

	des, err = fs.ReadDir(gen, "bud/view/about")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "about.svelte")
	is.Equal(des[0].IsDir(), false)
	is.Equal(des[0].Type(), fs.FileMode(0))
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.Name(), "about.svelte")
	is.Equal(fi.Mode(), fs.FileMode(0))
	stat, err = fs.Stat(gen, "bud/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(stat.Name(), "about.svelte")

	is.NoErr(fstest.TestFS(gen, "bud/view/index.svelte", "bud/view/about/about.svelte"))
}

func TestAll(t *testing.T) {
	is := is.New(t)
	gen := genfs.New()
	gen.GenerateDir("bud/view", view())

	// .
	file, err := gen.Open(".")
	is.NoErr(err)
	rgen, ok := file.(fs.ReadDirFile)
	is.True(ok)
	des, err := rgen.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "bud")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err := des[0].Info()
	is.NoErr(err)
	is.Equal(fi.IsDir(), true)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.ModeDir)
	is.Equal(fi.Name(), "bud")
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	stat, err := file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), ".")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// Stat .
	stat, err = fs.Stat(gen, ".")
	is.NoErr(err)
	is.Equal(stat.Name(), ".")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir .
	des, err = fs.ReadDir(gen, ".")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "bud")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)

	// bud
	file, err = gen.Open("bud")
	is.NoErr(err)
	rgen, ok = file.(fs.ReadDirFile)
	is.True(ok)
	des, err = rgen.ReadDir(-1)
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.IsDir(), true)
	is.True(fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.ModeDir)
	is.Equal(fi.Name(), "view")
	is.Equal(fi.Size(), int64(0))
	is.Equal(fi.Sys(), nil)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "bud")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// Stat bud
	stat, err = fs.Stat(gen, "bud")
	is.NoErr(err)
	is.Equal(stat.Name(), "bud")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir bud
	des, err = fs.ReadDir(gen, "bud")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[0].Type(), fs.ModeDir)

	// bud/view
	file, err = gen.Open("bud/view")
	is.NoErr(err)
	rgen, ok = file.(fs.ReadDirFile)
	is.True(ok)
	des, err = rgen.ReadDir(-1)
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
	is.Equal(fi.Size(), int64(14))
	is.Equal(fi.Sys(), nil)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "view")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// Stat bud
	stat, err = fs.Stat(gen, "bud/view")
	is.NoErr(err)
	is.Equal(stat.Name(), "view")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir bud
	des, err = fs.ReadDir(gen, "bud/view")
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
	is.Equal(fi.Size(), int64(14))
	is.Equal(fi.Sys(), nil)

	// bud/view/about
	file, err = gen.Open("bud/view/about")
	is.NoErr(err)
	rgen, ok = file.(fs.ReadDirFile)
	is.True(ok)
	des, err = rgen.ReadDir(-1)
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
	is.Equal(fi.Size(), int64(14))
	is.Equal(fi.Sys(), nil)
	stat, err = file.Stat()
	is.NoErr(err)
	is.Equal(stat.Name(), "about")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// Stat bud
	stat, err = fs.Stat(gen, "bud/view/about")
	is.NoErr(err)
	is.Equal(stat.Name(), "about")
	is.Equal(stat.Mode(), fs.ModeDir)
	is.True(stat.IsDir())
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(0))
	is.Equal(stat.Sys(), nil)
	// ReadDir bud
	des, err = fs.ReadDir(gen, "bud/view/about")
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
	is.Equal(fi.Size(), int64(14))
	is.Equal(fi.Sys(), nil)

	// bud/view/index.svelte
	// Open
	file, err = gen.Open("bud/view/index.svelte")
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
	stat, err = fs.Stat(gen, "bud/view/index.svelte")
	is.NoErr(err)
	is.Equal(stat.Name(), "index.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// ReadFile
	code, err := fs.ReadFile(gen, "bud/view/index.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h1>index</h1>`)

	// bud/view/about/about.svelte
	// Open
	file, err = gen.Open("bud/view/about/about.svelte")
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
	stat, err = fs.Stat(gen, "bud/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(stat.Name(), "about.svelte")
	is.Equal(stat.Mode(), fs.FileMode(0))
	is.Equal(stat.IsDir(), false)
	is.True(stat.ModTime().IsZero())
	is.Equal(stat.Size(), int64(14))
	is.Equal(stat.Sys(), nil)
	// ReadFile
	code, err = fs.ReadFile(gen, "bud/view/about/about.svelte")
	is.NoErr(err)
	is.Equal(string(code), `<h2>about</h2>`)

	// Run TestFS
	err = fstest.TestFS(gen, "bud", "bud/view", "bud/view/index.svelte", "bud/view/about/about.svelte")
	is.NoErr(err)
}

func TestDir(t *testing.T) {
	is := is.New(t)
	gen := genfs.New()
	gen.GenerateDir("bud/view", func(dir *genfs.Dir) error {
		dir.GenerateDir("about", func(dir *genfs.Dir) error {
			dir.GenerateDir("me", func(dir *genfs.Dir) error {
				return nil
			})
			return nil
		})
		dir.GenerateDir("users/admin", func(dir *genfs.Dir) error {
			return nil
		})
		return nil
	})
	des, err := fs.ReadDir(gen, "bud")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "view")
	is.Equal(des[0].IsDir(), true)
	des, err = fs.ReadDir(gen, "bud/view")
	is.NoErr(err)
	is.Equal(len(des), 2)
	is.Equal(des[0].Name(), "about")
	is.Equal(des[0].IsDir(), true)
	is.Equal(des[1].Name(), "users")
	is.Equal(des[1].IsDir(), true)
	des, err = fs.ReadDir(gen, "bud/view/about")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "me")
	is.Equal(des[0].IsDir(), true)
	des, err = fs.ReadDir(gen, "bud/view/about/me")
	is.NoErr(err)
	is.Equal(len(des), 0)
	des, err = fs.ReadDir(gen, "bud/view/users")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "admin")
	is.Equal(des[0].IsDir(), true)
	des, err = fs.ReadDir(gen, "bud/view/users/admin")
	is.NoErr(err)
	is.Equal(len(des), 0)
}

func TestGenerateFileError(t *testing.T) {
	is := is.New(t)
	gen := genfs.New()
	gen.GenerateFile("bud/main.go", func(file *genfs.File) error {
		return fs.ErrNotExist
	})
	code, err := fs.ReadFile(gen, "bud/main.go")
	is.True(err != nil)
	is.Equal(err.Error(), `genfs: error generating "bud/main.go". file does not exist`)
	is.True(errors.Is(err, fs.ErrNotExist))
	is.Equal(code, nil)
}

// func TestServeFile(t *testing.T) {
// 	is := is.New(t)
// 	gen := genfs.New()
// 	gen.GenerateDir("bud/view", func(dir *genfs.Dir) error {
// 		switch dir.Rel() {
// 		case ".":
// 			return fs.ErrInvalid
// 		case "_index.svelte":
// 			fmt.Println("HER?")
// 			dir.GenerateFile("_index.svelte", func(file *genfs.File) error {
// 				fmt.Println("OK!")
// 				file.Data = []byte(dir.Target() + "'s data")
// 				return nil
// 			})
// 			return nil
// 		case "about/_about.svelte":
// 			return fs.ErrNotExist
// 		default:
// 			return fs.ErrNotExist
// 		}
// 	})
// 	des, err := fs.ReadDir(gen, "bud/view")
// 	is.True(errors.Is(err, fs.ErrInvalid))
// 	is.Equal(len(des), 0)

// 	// _index.svelte
// 	file, err := gen.Open("bud/view/_index.svelte")
// 	is.NoErr(err)
// 	stat, err := file.Stat()
// 	is.NoErr(err)
// 	is.Equal(stat.Name(), "_index.svelte")
// 	is.Equal(stat.Mode(), fs.FileMode(0))
// 	is.Equal(stat.IsDir(), false)
// 	is.True(stat.ModTime().IsZero())
// 	is.Equal(stat.Size(), int64(29))
// 	is.Equal(stat.Sys(), nil)
// 	code, err := fs.ReadFile(gen, "bud/view/_index.svelte")
// 	is.NoErr(err)
// 	is.Equal(string(code), `bud/view/_index.svelte's data`)

// 	// about/_about.svelte
// 	file, err = gen.Open("bud/view/about/_about.svelte")
// 	is.NoErr(err)
// 	stat, err = file.Stat()
// 	is.NoErr(err)
// 	is.Equal(stat.Name(), "_about.svelte")
// 	is.Equal(stat.Mode(), fs.FileMode(0))
// 	is.Equal(stat.IsDir(), false)
// 	is.True(stat.ModTime().IsZero())
// 	is.Equal(stat.Size(), int64(35))
// 	is.Equal(stat.Sys(), nil)
// 	code, err = fs.ReadFile(gen, "bud/view/about/_about.svelte")
// 	is.NoErr(err)
// 	is.Equal(string(code), `bud/view/about/_about.svelte's data`)
// }

// func TestHTTP(t *testing.T) {
// 	is := is.New(t)
// 	cfs := conjure.New()
// 	cfs.ServeFile("bud/view", func(file *genfs.File) error {
// 		file.Data = []byte(file.Path() + `'s data`)
// 		return nil
// 	})
// 	hfs := http.FS(cfs)

// 	handler := func(w http.ResponseWriter, r *http.Request) {
// 		file, err := hfs.Open(r.URL.Path)
// 		if err != nil {
// 			http.Error(w, err.Error(), 500)
// 			return
// 		}
// 		stat, err := file.Stat()
// 		if err != nil {
// 			http.Error(w, err.Error(), 500)
// 			return
// 		}
// 		w.Header().Add("Content-Type", "text/javascript")
// 		http.ServeContent(w, r, r.URL.Path, stat.ModTime(), file)
// 	}

// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest("GET", "/bud/view/_index.svelte", nil)
// 	handler(w, r)

// 	response := w.Result()
// 	body, err := ioutil.ReadAll(response.Body)
// 	is.NoErr(err)
// 	is.Equal(string(body), `bud/view/_index.svelte's data`)
// 	is.Equal(response.StatusCode, 200)
// }

// func rootless(fpath string) string {
// 	parts := strings.Split(fpath, string(filepath.Separator))
// 	return path.Join(parts[1:]...)
// }

// func TestTargetPath(t *testing.T) {
// 	is := is.New(t)
// 	// Test inner file and rootless
// 	cfs := conjure.New()
// 	cfs.GenerateDir("bud/view", func(dir *genfs.Dir) error {
// 		dir.GenerateFile("about/about.svelte", func(file *genfs.File) error {
// 			file.Data = []byte(rootless(file.Path()))
// 			return nil
// 		})
// 		return nil
// 	})
// 	code, err := fs.ReadFile(cfs, "bud/view/about/about.svelte")
// 	is.NoErr(err)
// 	is.Equal(string(code), "view/about/about.svelte")
// }

// func TestDynamicDir(t *testing.T) {
// 	is := is.New(t)
// 	cfs := conjure.New()
// 	cfs.GenerateDir("bud/view", func(dir *genfs.Dir) error {
// 		doms := []string{"about/about.svelte", "index.svelte"}
// 		for _, dom := range doms {
// 			dom := dom
// 			dir.GenerateFile(dom, func(file *genfs.File) error {
// 				file.Data = []byte(`<h1>` + dom + `</h1>`)
// 				return nil
// 			})
// 		}
// 		return nil
// 	})
// 	des, err := fs.ReadDir(cfs, "bud/view")
// 	is.NoErr(err)
// 	is.Equal(len(des), 2)
// 	is.Equal(des[0].Name(), "about")
// 	is.Equal(des[1].Name(), "index.svelte")
// 	des, err = fs.ReadDir(cfs, "bud/view/about")
// 	is.NoErr(err)
// 	is.Equal(len(des), 1)
// 	is.Equal(des[0].Name(), "about.svelte")
// }

// func TestBases(t *testing.T) {
// 	is := is.New(t)
// 	cfs := conjure.New()
// 	cfs.GenerateDir("bud/view", func(dir *genfs.Dir) error {
// 		return nil
// 	})
// 	cfs.GenerateDir("bud/controller", func(dir *genfs.Dir) error {
// 		return nil
// 	})
// 	stat, err := fs.Stat(cfs, "bud/controller")
// 	is.NoErr(err)
// 	is.Equal(stat.Name(), "controller")
// 	stat, err = fs.Stat(cfs, "bud/view")
// 	is.NoErr(err)
// 	is.Equal(stat.Name(), "view")
// }

// func TestDirPath(t *testing.T) {
// 	is := is.New(t)
// 	cfs := conjure.New()
// 	cfs.GenerateDir("bud/view", func(dir *genfs.Dir) error {
// 		dir.GenerateDir("public", func(dir *genfs.Dir) error {
// 			dir.GenerateFile("favicon.ico", func(file *genfs.File) error {
// 				file.Data = []byte("cool_favicon.ico")
// 				return nil
// 			})
// 			return nil
// 		})
// 		return nil
// 	})
// 	cfs.GenerateDir("bud", func(dir *genfs.Dir) error {
// 		dir.GenerateDir("controller", func(dir *genfs.Dir) error {
// 			dir.GenerateFile("controller.go", func(file *genfs.File) error {
// 				file.Data = []byte("package controller")
// 				return nil
// 			})
// 			return nil
// 		})
// 		return nil
// 	})
// 	code, err := fs.ReadFile(cfs, "bud/view/public/favicon.ico")
// 	is.NoErr(err)
// 	is.Equal(string(code), "cool_favicon.ico")
// 	code, err = fs.ReadFile(cfs, "bud/controller/controller.go")
// 	is.NoErr(err)
// 	is.Equal(string(code), "package controller")
// }

// func TestDirMerge(t *testing.T) {
// 	is := is.New(t)
// 	cfs := conjure.New()
// 	cfs.GenerateDir("bud/view", func(dir *genfs.Dir) error {
// 		dir.GenerateFile("index.svelte", func(file *genfs.File) error {
// 			file.Data = []byte(`<h1>index</h1>`)
// 			return nil
// 		})
// 		dir.GenerateDir("somedir", func(dir *genfs.Dir) error {
// 			return nil
// 		})
// 		return nil
// 	})
// 	cfs.GenerateFile("bud/view/view.go", func(file *genfs.File) error {
// 		file.Data = []byte(`package view`)
// 		return nil
// 	})
// 	cfs.GenerateFile("bud/view/plugin.go", func(file *genfs.File) error {
// 		file.Data = []byte(`package plugin`)
// 		return nil
// 	})
// 	// bud/view
// 	des, err := fs.ReadDir(cfs, "bud/view")
// 	is.NoErr(err)
// 	is.Equal(len(des), 4)
// 	is.Equal(des[0].Name(), "index.svelte")
// 	is.Equal(des[0].IsDir(), false)
// 	is.Equal(des[1].Name(), "plugin.go")
// 	is.Equal(des[1].IsDir(), false)
// 	is.Equal(des[2].Name(), "somedir")
// 	is.Equal(des[2].IsDir(), true)
// 	is.Equal(des[3].Name(), "view.go")
// 	is.Equal(des[3].IsDir(), false)
// }

// func TestAddGenerator(t *testing.T) {
// 	is := is.New(t)
// 	// Add the view
// 	cfs := conjure.New()
// 	cfs.GenerateDir("bud/view", View())

// 	// Add the controller
// 	cfs.GenerateDir("bud/controller", func(dir *genfs.Dir) error {
// 		dir.GenerateFile("controller.go", func(file *genfs.File) error {
// 			file.Data = []byte(`package controller`)
// 			return nil
// 		})
// 		return nil
// 	})

// 	des, err := fs.ReadDir(cfs, "bud")
// 	is.NoErr(err)
// 	is.Equal(len(des), 2)
// 	is.Equal(des[0].Name(), "controller")
// 	is.Equal(des[1].Name(), "view")

// 	// Read from view
// 	code, err := fs.ReadFile(cfs, "bud/view/index.svelte")
// 	is.NoErr(err)
// 	is.Equal(string(code), `<h1>index</h1>`)

// 	// Read from controller
// 	code, err = fs.ReadFile(cfs, "bud/controller/controller.go")
// 	is.NoErr(err)
// 	is.Equal(string(code), `package controller`)
// }

// type commandGenerator struct {
// 	Input string
// }

// func (c *commandGenerator) GenerateFile(file *genfs.File) error {
// 	file.Data = []byte(c.Input + c.Input)
// 	return nil
// }

// func (c *commandGenerator) GenerateDir(dir *genfs.Dir) error {
// 	dir.GenerateFile("index.svelte", func(file *genfs.File) error {
// 		file.Data = []byte(c.Input + c.Input)
// 		return nil
// 	})
// 	return nil
// }

// func (c *commandGenerator) ServeFile(file *genfs.File) error {
// 	file.Data = []byte(c.Input + "/" + file.Path())
// 	return nil
// }

// func TestFileGenerator(t *testing.T) {
// 	is := is.New(t)
// 	cfs := conjure.New()
// 	cfs.FileGenerator("bud/command/command.go", &commandGenerator{Input: "a"})
// 	code, err := fs.ReadFile(cfs, "bud/command/command.go")
// 	is.NoErr(err)
// 	is.Equal(string(code), "aa")
// }

// func TestDirGenerator(t *testing.T) {
// 	is := is.New(t)
// 	// Add the view
// 	cfs := conjure.New()
// 	cfs.DirGenerator("bud/view", &commandGenerator{Input: "a"})
// 	code, err := fs.ReadFile(cfs, "bud/view/index.svelte")
// 	is.NoErr(err)
// 	is.Equal(string(code), "aa")
// }

// func TestFileServer(t *testing.T) {
// 	is := is.New(t)
// 	cfs := conjure.New()
// 	cfs.FileServer("bud/view", &commandGenerator{Input: "a"})
// 	code, err := fs.ReadFile(cfs, "bud/view/index.svelte")
// 	is.NoErr(err)
// 	is.Equal(string(code), "a/bud/view/index.svelte")
// }

// func TestDotReadDirEmpty(t *testing.T) {
// 	is := is.New(t)
// 	cfs := conjure.New()
// 	cfs.GenerateFile("bud/generate/main.go", func(file *genfs.File) error {
// 		file.Data = []byte("package main")
// 		return nil
// 	})
// 	cfs.GenerateFile("go.mod", func(file *genfs.File) error {
// 		file.Data = []byte("module pkg")
// 		return nil
// 	})
// 	des, err := fs.ReadDir(cfs, ".")
// 	is.NoErr(err)
// 	is.Equal(len(des), 2)
// }

// func TestDotReadDirFiles(t *testing.T) {
// 	is := is.New(t)
// 	tmp := t.TempDir()
// 	err := os.WriteFile(filepath.Join(tmp, "a.txt"), []byte("a"), 0644)
// 	is.NoErr(err)
// 	err = os.WriteFile(filepath.Join(tmp, "b.txt"), []byte("b"), 0644)
// 	is.NoErr(err)
// 	cfs := conjure.New()
// 	mapfs := fstest.MapFS{
// 		"a.txt": &fstest.MapFile{Data: []byte("a"), Mode: 0644},
// 		"b.txt": &fstest.MapFile{Data: []byte("b"), Mode: 0644},
// 	}
// 	cfs.GenerateFile("bud/generate/main.go", func(file *genfs.File) error {
// 		file.Data = []byte("package main")
// 		return nil
// 	})
// 	cfs.GenerateFile("go.mod", func(file *genfs.File) error {
// 		file.Data = []byte("module pkg")
// 		return nil
// 	})
// 	fsys := merged.Merge(cfs, mapfs)
// 	des, err := fs.ReadDir(fsys, ".")
// 	is.NoErr(err)
// 	is.Equal(len(des), 4)
// }

// func TestReadDirDuplicates(t *testing.T) {
// 	is := is.New(t)
// 	mapfs := fstest.MapFS{
// 		"go.mod": &fstest.MapFile{Data: []byte(`module app.com`)},
// 	}
// 	cfs := conjure.New()
// 	cfs.GenerateFile("go.mod", func(file *genfs.File) error {
// 		file.Data = []byte("module app.cool")
// 		return nil
// 	})
// 	fsys := merged.Merge(cfs, mapfs)
// 	des, err := fs.ReadDir(fsys, ".")
// 	is.NoErr(err)
// 	is.Equal(len(des), 1)
// 	is.Equal(des[0].Name(), "go.mod")
// 	code, err := fs.ReadFile(fsys, "go.mod")
// 	is.NoErr(err)
// 	is.Equal(string(code), "module app.cool")
// }

// func TestEmbedOpen(t *testing.T) {
// 	is := is.New(t)
// 	now := time.Now()
// 	cfs := conjure.New()
// 	cfs.FileGenerator("bud/view/index.svelte", &conjure.Embed{
// 		Data:    []byte(`<h1>index</h1>`),
// 		Mode:    fs.FileMode(0644),
// 		ModTime: now,
// 	})
// 	cfs.FileGenerator("bud/view/about/about.svelte", &conjure.Embed{
// 		Data:    []byte(`<h1>about</h1>`),
// 		Mode:    fs.FileMode(0644),
// 		ModTime: now,
// 	})
// 	cfs.FileGenerator("bud/public/favicon.ico", &conjure.Embed{
// 		Data:    []byte(`favicon.ico`),
// 		Mode:    fs.FileMode(0644),
// 		ModTime: now,
// 	})
// 	// bud/view/index.svelte
// 	code, err := fs.ReadFile(cfs, "bud/view/index.svelte")
// 	is.NoErr(err)
// 	is.Equal(string(code), `<h1>index</h1>`)
// 	stat, err := fs.Stat(cfs, "bud/view/index.svelte")
// 	is.NoErr(err)
// 	is.Equal(stat.ModTime(), now)
// 	is.Equal(stat.Mode(), fs.FileMode(0644))
// 	is.Equal(stat.IsDir(), false)

// 	// bud/view/about/about.svelte
// 	code, err = fs.ReadFile(cfs, "bud/view/about/about.svelte")
// 	is.NoErr(err)
// 	is.Equal(string(code), `<h1>about</h1>`)
// 	stat, err = fs.Stat(cfs, "bud/view/about/about.svelte")
// 	is.NoErr(err)
// 	is.Equal(stat.ModTime(), now)
// 	is.Equal(stat.Mode(), fs.FileMode(0644))
// 	is.Equal(stat.IsDir(), false)

// 	// bud/public/favicon.ico
// 	code, err = fs.ReadFile(cfs, "bud/public/favicon.ico")
// 	is.NoErr(err)
// 	is.Equal(string(code), `favicon.ico`)
// 	stat, err = fs.Stat(cfs, "bud/public/favicon.ico")
// 	is.NoErr(err)
// 	is.Equal(stat.ModTime(), now)
// 	is.Equal(stat.Mode(), fs.FileMode(0644))
// 	is.Equal(stat.IsDir(), false)

// 	// bud/public
// 	// TODO: consider locking this down, though this might be taken care of higher
// 	// up in the stack.
// 	des, err := fs.ReadDir(cfs, "bud/public")
// 	is.NoErr(err)
// 	is.Equal(len(des), 1)
// 	is.Equal(des[0].Name(), "favicon.ico")
// }

// func TestGoModGoMod(t *testing.T) {
// 	is := is.New(t)
// 	cfs := conjure.New()
// 	cfs.GenerateFile("go.mod", func(file *genfs.File) error {
// 		file.Data = []byte("module app.com\nrequire mod.test/module v1.2.4")
// 		return nil
// 	})
// 	stat, err := fs.Stat(cfs, "go.mod/go.mod")
// 	is.True(err != nil)
// 	is.True(errors.Is(err, fs.ErrNotExist))
// 	is.Equal(stat, nil)
// 	stat, err = fs.Stat(cfs, "go.mod")
// 	is.NoErr(err)
// 	is.Equal(stat.Name(), "go.mod")
// }

// func TestGoModGoModEmbed(t *testing.T) {
// 	is := is.New(t)
// 	cfs := conjure.New()
// 	cfs.FileGenerator("go.mod", &conjure.Embed{
// 		Data: []byte("module app.com\nrequire mod.test/module v1.2.4"),
// 	})
// 	stat, err := fs.Stat(cfs, "go.mod/go.mod")
// 	is.True(err != nil)
// 	is.True(errors.Is(err, fs.ErrNotExist))
// 	is.Equal(stat, nil)
// 	stat, err = fs.Stat(cfs, "go.mod")
// 	is.NoErr(err)
// 	is.Equal(stat.Name(), "go.mod")
// }

// func TestMount(t *testing.T) {
// 	is := is.New(t)
// 	now := time.Now()
// 	cfs := conjure.New()
// 	cfs.GenerateDir("bud/view", View())
// 	gfs := conjure.New()
// 	cfs.Mount("bud/generator", gfs)
// 	gfs.FileGenerator("bud/generator/tailwind/tailwind.css", &conjure.Embed{
// 		Data:    []byte(`/** tailwind **/`),
// 		Mode:    fs.FileMode(0644),
// 		ModTime: now,
// 	})
// 	des, err := fs.ReadDir(cfs, "bud")
// 	is.NoErr(err)
// 	is.Equal(len(des), 2)
// 	is.Equal(des[0].Name(), "generator")
// 	is.Equal(des[0].IsDir(), true)
// 	is.Equal(des[1].Name(), "view")
// 	is.Equal(des[1].IsDir(), true)
// 	des, err = fs.ReadDir(cfs, "bud/generator")
// 	is.NoErr(err)
// 	is.Equal(len(des), 1)
// 	is.Equal(des[0].Name(), "tailwind")
// 	is.Equal(des[0].IsDir(), true)
// 	des, err = fs.ReadDir(cfs, "bud/generator/tailwind")
// 	is.NoErr(err)
// 	is.Equal(len(des), 1)
// 	is.Equal(des[0].Name(), "tailwind.css")
// 	is.Equal(des[0].IsDir(), false)
// 	is.Equal(des[0].Type(), fs.FileMode(0))
// 	fi, err := des[0].Info()
// 	is.NoErr(err)
// 	is.True(fi.ModTime().IsZero())
// 	is.Equal(fi.Mode(), fs.FileMode(0))
// 	is.Equal(fi.IsDir(), false)
// 	is.Equal(fi.Size(), int64(0))
// 	code, err := fs.ReadFile(cfs, "bud/generator/tailwind/tailwind.css")
// 	is.NoErr(err)
// 	is.Equal(string(code), `/** tailwind **/`)
// }
