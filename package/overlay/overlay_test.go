package overlay_test

import (
	"os"
	"path/filepath"
	"testing"

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

func TestAll(t *testing.T) {
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
	_ = ofs
}
