package pluginfs_test

import (
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/testplugin"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/pluginfs"
	"github.com/livebud/bud/package/vfs"
	"github.com/matryer/is"
)

func TestMergeModules(t *testing.T) {
	is := is.New(t)
	appDir := t.TempDir()
	dep, err := testplugin.Plugin()
	is.NoErr(err)
	err = vfs.Write(appDir, vfs.Map{
		"public/normalize.css": []byte(`/* normalize */`),
		"go.mod": []byte(`
			module app.com
			require ` + dep.Path + ` ` + dep.Version + `
		`),
	})
	is.NoErr(err)
	module, err := gomod.Find(appDir)
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
	is.Equal(fi.Mode(), fs.FileMode(0555|fs.ModeDir))
	is.Equal(fi.Name(), "tailwind")
	is.True(fi.Size() != 0)
	is.True(fi.Sys() != nil)
	fi, err = fs.Stat(pfs, "public/tailwind")
	is.NoErr(err)
	is.Equal(fi.IsDir(), true)
	is.True(!fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0555|fs.ModeDir))
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
	is.Equal(fi.Mode(), fs.FileMode(0444))
	is.Equal(fi.Name(), "preflight.css")
	is.Equal(fi.Size(), int64(14))
	is.True(fi.Sys() != nil)
	fi, err = fs.Stat(pfs, "public/tailwind/preflight.css")
	is.NoErr(err)
	is.Equal(fi.IsDir(), false)
	is.True(!fi.ModTime().IsZero())
	is.Equal(fi.Mode(), fs.FileMode(0444))
	is.Equal(fi.Name(), "preflight.css")
	is.Equal(fi.Size(), int64(14))
	is.True(fi.Sys() != nil)
	code, err = fs.ReadFile(pfs, "public/tailwind/preflight.css")
	is.NoErr(err)
	is.Equal(string(code), `/* tailwind */`)
}

func TestMultiple(t *testing.T) {
	is := is.New(t)
	appDir := t.TempDir()
	dep, err := testplugin.Plugin()
	is.NoErr(err)
	nested, err := testplugin.NestedPlugin()
	is.NoErr(err)
	err = vfs.Write(appDir, vfs.Map{
		"go.mod": []byte(`
			module app.com
			require ` + dep.Path + ` ` + dep.Version + `
			require ` + nested.Path + ` ` + nested.Version + `
		`),
	})
	is.NoErr(err)
	module, err := gomod.Find(appDir)
	is.NoErr(err)
	pfs, err := pluginfs.Load(module)
	is.NoErr(err)

	// From bud-test-plugin
	code, err := fs.ReadFile(pfs, "public/tailwind/preflight.css")
	is.NoErr(err)
	is.Equal(string(code), `/* tailwind */`)

	// From bud-test-nested-plugin
	code, err = fs.ReadFile(pfs, "public/admin.css")
	is.NoErr(err)
	is.Equal(string(code), `/* admin.css */`)
}

func TestPluginDir(t *testing.T) {
	t.SkipNow()
}
