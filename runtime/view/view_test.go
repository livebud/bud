package view_test

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/matryer/is"
	"github.com/matthewmueller/gotext"
	"gitlab.com/mnm/bud/pkg/vfs"
	"gitlab.com/mnm/bud/runtime/view"
)

// func TestSvelteView(t *testing.T) {
// 	is := is.New(t)
// 	cwd, err := os.Getwd()
// 	is.NoErr(err)
// 	dir := filepath.Join(cwd, "_tmp")
// 	is.NoErr(os.RemoveAll(dir))
// 	defer func() {
// 		if !t.Failed() {
// 			is.NoErr(os.RemoveAll(dir))
// 		}
// 	}()
// 	memfs := vfs.Memory{
// 		"view/index.svelte": &fstest.MapFile{
// 			Data: []byte(`<h1>hi world</h1>`),
// 		},
// 	}
// 	is.NoErr(vfs.WriteAll(memfs, ".", dir))
// 	dirfs := os.DirFS(dir)
// 	svelte := svelte.New(&svelte.Input{
// 		VM:  v8.New(),
// 		Dev: true,
// 	})
// 	dfs := dfs.New(dirfs, map[string]dfs.Generator{
// 		"duo/view/_ssr.js": ssr.Generator(dirfs, svelte, dir),
// 	})
// 	// vm := v8.New()
// 	// view := view.New(dfs, vm)

// 	// Install svelte
// 	err = npm.Install(dir, "svelte@3.42.3")
// 	is.NoErr(err)
// 	// Read the wrapped version of index.svelte with node_modules rewritten
// 	code, err := fs.ReadFile(dfs, "duo/view/_ssr.js")
// 	is.NoErr(err)
// 	is.True(strings.Contains(string(code), `create_ssr_component(`))
// 	is.True(strings.Contains(string(code), `<h1>hi world</h1>`))
// 	is.True(strings.Contains(string(code), `views["/"] = `))
// }

func TestMain(m *testing.M) {
	os.RemoveAll("_tmp")
	code := m.Run()
	if code == 0 {
		os.RemoveAll("_tmp")
	}
	os.Exit(code)
}

func testDir(t testing.TB) (string, func()) {
	t.Helper()
	is := is.New(t)
	dir := filepath.Join("_tmp", gotext.Camel(t.Name()))
	wd, err := os.Getwd()
	is.NoErr(err)
	absdir := filepath.Join(wd, dir)
	is.NoErr(os.RemoveAll(dir))
	return absdir, func() {
		if !t.Failed() {
			is.NoErr(os.RemoveAll(dir))
		}
	}
}

type RenderTest struct {
	Name   string
	Files  map[string]string
	Path   string
	Props  string
	Error  string
	Expect func(is *is.I, r *view.Response)
}

var renderTests = []*RenderTest{
	{
		Name: "index.svelte",
		Files: map[string]string{
			"view/index.svelte": `<h1>hi world</h1>`,
		},
		Path:  "/",
		Props: `{}`,
		Expect: func(is *is.I, r *view.Response) {
			is.Equal(200, r.Status)
		},
	},
}

func TestRunRender(t *testing.T) {
	for _, test := range renderTests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			is := is.New(t)
			dir, removeDir := testDir(t)
			defer removeDir()
			// Prepare main.go and files
			mem := vfs.Memory{}
			for path, data := range test.Files {
				mem[path] = &fstest.MapFile{Data: []byte(data)}
			}
			is.NoErr(vfs.WriteAll(".", dir, mem))
			// Sync the filesystem
			// testRender(is, dir, Runner(dir), test)
		})
	}
}
