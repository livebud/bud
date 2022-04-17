package pluginfs_test

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/modcache"
	"github.com/livebud/bud/package/pluginfs"
	"github.com/livebud/bud/package/vfs"
	"github.com/matryer/is"
)

func TestMergeModules(t *testing.T) {
	is := is.New(t)
	cacheDir := t.TempDir()
	modCache := modcache.New(cacheDir)
	preflight := `/* tailwind */`
	err := modCache.Write(map[string]modcache.Files{
		"github.com/livebud/bud-tailwind@v0.0.1": modcache.Files{
			"public/tailwind/preflight.css": preflight,
		},
	})
	is.NoErr(err)
	appDir := t.TempDir()
	err = vfs.Write(appDir, vfs.Map{
		"public/normalize.css": []byte(`/* normalize */`),
		"go.mod":               []byte("module app.com\nrequire github.com/livebud/bud-tailwind v0.0.1"),
	})
	is.NoErr(err)
	module, err := gomod.Find(appDir, gomod.WithModCache(modCache))
	is.NoErr(err)
	pfs, err := pluginfs.Load(module)
	is.NoErr(err)

	des, err := fs.ReadDir(pfs, "public")
	is.NoErr(err)
	is.Equal(len(des), 2)

	// normalize.css
	is.Equal(des[0].Name(), "normalize.css")
	is.Equal(des[0].IsDir(), false)
	// TODO: see if we can fix the FileMode
	is.Equal(des[0].Type(), fs.FileMode(0))
	fi, err := des[0].Info()
	is.NoErr(err)
	is.Equal(fi.IsDir(), false)
	is.True(!fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0644))
	is.Equal(fi.Name(), "normalize.css")
	is.Equal(fi.Size(), int64(15))
	is.True(fi.Sys() != nil)
	fi, err = fs.Stat(pfs, "public/normalize.css")
	is.NoErr(err)
	is.Equal(fi.IsDir(), false)
	is.True(!fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0644))
	is.Equal(fi.Name(), "normalize.css")
	is.Equal(fi.Size(), int64(15))
	is.True(fi.Sys() != nil)
	code, err := fs.ReadFile(pfs, "public/normalize.css")
	is.NoErr(err)
	is.Equal(string(code), `/* normalize */`)

	// tailwind
	is.Equal(des[1].Name(), "tailwind")
	is.Equal(des[1].IsDir(), true)
	// TODO: see if we can fix the FileMode
	is.Equal(des[1].Type(), fs.FileMode(fs.ModeDir))
	fi, err = des[1].Info()
	is.NoErr(err)
	is.Equal(fi.IsDir(), true)
	is.True(!fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0755|fs.ModeDir))
	is.Equal(fi.Name(), "tailwind")
	is.True(fi.Size() != 0)
	is.True(fi.Sys() != nil)
	fi, err = fs.Stat(pfs, "public/tailwind")
	is.NoErr(err)
	is.Equal(fi.IsDir(), true)
	is.True(!fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0755|fs.ModeDir))
	is.Equal(fi.Name(), "tailwind")
	is.True(fi.Size() != 0)
	is.True(fi.Sys() != nil)

	// tailwind/preflight.css
	des, err = fs.ReadDir(pfs, "public/tailwind")
	is.NoErr(err)
	is.Equal(len(des), 1)
	is.Equal(des[0].Name(), "preflight.css")
	is.Equal(des[0].IsDir(), false)
	// TODO: see if we can fix the FileMode
	is.Equal(des[0].Type(), fs.FileMode(0))
	fi, err = des[0].Info()
	is.NoErr(err)
	is.Equal(fi.IsDir(), false)
	is.True(!fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0644))
	is.Equal(fi.Name(), "preflight.css")
	is.Equal(fi.Size(), int64(14))
	is.True(fi.Sys() != nil)
	fi, err = fs.Stat(pfs, "public/tailwind/preflight.css")
	is.NoErr(err)
	is.Equal(fi.IsDir(), false)
	is.True(!fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0644))
	is.Equal(fi.Name(), "preflight.css")
	is.Equal(fi.Size(), int64(14))
	is.True(fi.Sys() != nil)
	code, err = fs.ReadFile(pfs, "public/tailwind/preflight.css")
	is.NoErr(err)
	is.Equal(string(code), `/* tailwind */`)
}

func TestPlugin(t *testing.T) {
	t.SkipNow()
}
