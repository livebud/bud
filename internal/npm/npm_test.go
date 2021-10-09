package npm_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-duo/bud/internal/npm"
	"github.com/go-duo/bud/internal/vfs"
	"github.com/matryer/is"
)

func exists(t testing.TB, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}

func TestInstallSvelte(t *testing.T) {
	is := is.New(t)
	is.NoErr(os.RemoveAll("_tmp"))
	defer func() {
		if !t.Failed() {
			is.NoErr(os.RemoveAll("_tmp"))
		}
	}()
	err := npm.Install("_tmp", "svelte@3.42.3", "uid@2.0.0")
	is.NoErr(err)
	exists(t, filepath.Join("_tmp", "node_modules", "svelte", "package.json"))
	exists(t, filepath.Join("_tmp", "node_modules", "uid", "package.json"))
	exists(t, filepath.Join("_tmp", "node_modules", "svelte", "internal", "index.js"))
}
func TestLinkLiveBud(t *testing.T) {
	is := is.New(t)
	is.NoErr(os.RemoveAll("_tmp"))
	defer func() {
		if !t.Failed() {
			is.NoErr(os.RemoveAll("_tmp"))
		}
	}()
	err := vfs.WriteAll(".", "_tmp", vfs.Memory{
		"package.json": &vfs.File{Data: []byte(`{}`)},
	})
	is.NoErr(err)
	err = npm.Link("../../budjs", "_tmp")
	is.NoErr(err)
	exists(t, filepath.Join("_tmp", "node_modules", "livebud", "package.json"))
	exists(t, filepath.Join("_tmp", "node_modules", "livebud", "svelte.ts"))
	exists(t, filepath.Join("_tmp", "node_modules", "livebud", "hot.ts"))
	exists(t, filepath.Join("_tmp", "node_modules", "livebud", "index.ts"))
	exists(t, filepath.Join("_tmp", "node_modules", "livebud", "jsx.ts"))
}
