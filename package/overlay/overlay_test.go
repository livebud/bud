package overlay_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"io/fs"

	"gitlab.com/mnm/bud/package/overlay"

	"github.com/matryer/is"
	"gitlab.com/mnm/bud/pkg/gomod"
	"gitlab.com/mnm/bud/pkg/modcache"
)

const modFile = `
module app.com

require gitlab.com/mnm/bud-tailwind v0.0.1
require gitlab.com/mnm/bud-lambda v1.0.0
`

func TestPlugins(t *testing.T) {
	is := is.New(t)
	cacheDir := t.TempDir()
	modCache := modcache.New(cacheDir)
	err := modCache.Write(map[string]modcache.Files{
		"gitlab.com/mnm/bud-tailwind@v0.0.1": modcache.Files{
			"public/tailwind/preflight.css": `/* tailwind */`,
		},
		"gitlab.com/mnm/bud-lambda@v1.0.0": modcache.Files{
			"command/lambda/lambda.go": `package lambda`,
		},
	})
	is.NoErr(err)
	appDir := t.TempDir()
	err = os.WriteFile(filepath.Join(appDir, "go.mod"), []byte(modFile), 0644)
	is.NoErr(err)
	module, err := gomod.Find(appDir, gomod.WithModCache(modCache))
	is.NoErr(err)
	ofs, err := overlay.Load(module)
	is.NoErr(err)
	code, err := fs.ReadFile(ofs, "public/tailwind/preflight.css")
	is.NoErr(err)
	is.Equal(string(code), `/* tailwind */`)
	code, err = fs.ReadFile(ofs, "command/lambda/lambda.go")
	is.NoErr(err)
	is.Equal(string(code), `package lambda`)
}

type ctxKey string

func TestContextPropagation(t *testing.T) {
	t.SkipNow()
	is := is.New(t)
	appDir := t.TempDir()
	err := os.WriteFile(filepath.Join(appDir, "go.mod"), []byte(`module app.com`), 0644)
	is.NoErr(err)
	module, err := gomod.Find(appDir)
	is.NoErr(err)
	ofs, err := overlay.Load(module)
	is.NoErr(err)
	ofs.GenerateFile("public/normalize.css", func(ctx context.Context, fsys overlay.F, file *overlay.File) error {
		test := ctx.Value(ctxKey("test")).(string)
		is.Equal(test, "test")
		file.Data = []byte("/* normalize */")
		return nil
	})
	// ctx := context.WithValue(context.Background(), ctxKey("test"), "test")
	code, err := fs.ReadFile(ofs, "public/normalize.css")
	is.NoErr(err)
	is.Equal(string(code), `/* normalize */`)
}
