package testdir_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/internal/testdir"
	"gitlab.com/mnm/bud/pkg/modcache"
)

func exists(fsys fs.FS, paths ...string) error {
	for _, path := range paths {
		if _, err := fs.Stat(fsys, path); err != nil {
			return err
		}
	}
	return nil
}

func TestDir(t *testing.T) {
	is := is.New(t)
	td := testdir.New()
	td.Modules = map[string]modcache.Files{
		"gitlab.com/mnm/bud-tailwind@v0.0.1": modcache.Files{
			"public/tailwind/preflight.css": `/* tailwind */`,
		},
	}
	td.Files["action/action.go"] = `package action`
	td.BFiles["public/favicon.ico"] = []byte{0x00}
	td.NodeModules["svelte"] = `3.46.4`
	dir := t.TempDir()
	err := td.Write(dir)
	is.NoErr(err)
	err = exists(os.DirFS(dir),
		"action/action.go",
		".mod/gitlab.com/mnm/bud-tailwind@v0.0.1/public/tailwind/preflight.css",
		"node_modules/svelte/package.json",
		"package.json",
		"go.mod",
	)
	is.NoErr(err)
}
