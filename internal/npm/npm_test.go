package npm_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/npm"
	"gitlab.com/mnm/bud/pkg/vfs"
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
	err = npm.Link("../../livebud", "_tmp")
	is.NoErr(err)
	exists(t, filepath.Join("_tmp", "node_modules", "livebud", "package.json"))
	exists(t, filepath.Join("_tmp", "node_modules", "livebud", "runtime", "svelte", "index.ts"))
	exists(t, filepath.Join("_tmp", "node_modules", "livebud", "runtime", "hot", "index.ts"))
	exists(t, filepath.Join("_tmp", "node_modules", "livebud", "runtime", "index.ts"))
	exists(t, filepath.Join("_tmp", "node_modules", "livebud", "runtime", "jsx", "index.ts"))
}
