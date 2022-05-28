package npm_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/npm"
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
