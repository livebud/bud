package pluginfs_test

import (
	"context"
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/testdir"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/pluginfs"
)

func TestMergeModules(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Files["public/normalize.css"] = `/* normalize */`
	td.Modules["github.com/livebud/bud-test-plugin"] = `v0.0.8`
	err := td.Write(ctx)
	is.NoErr(err)

	module, err := gomod.Find(dir)
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
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)

	td.Files["public/normalize.css"] = `/* normalize */`
	td.Modules["github.com/livebud/bud-test-plugin"] = `v0.0.9`
	td.Modules["github.com/livebud/bud-test-nested-plugin"] = `v0.0.5`
	err := td.Write(ctx)
	is.NoErr(err)

	module, err := gomod.Find(dir)
	is.NoErr(err)
	pfs, err := pluginfs.Load(module)
	is.NoErr(err)

	// From bud-test-plugin
	code, err := fs.ReadFile(pfs, "public/base.css")
	is.NoErr(err)
	is.Equal(string(code), `/* base */`)

	// From bud-test-nested-plugin
	code, err = fs.ReadFile(pfs, "public/admin.css")
	is.NoErr(err)
	is.Equal(string(code), `/* admin.css */`)
}

func TestConflicts(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)

	td.Files["public/admin.css"] = `/* app admin.css */`
	td.Modules["github.com/livebud/bud-test-plugin"] = `v0.0.9`
	td.Modules["github.com/livebud/bud-test-nested-plugin"] = `v0.0.5`
	err := td.Write(ctx)
	is.NoErr(err)

	module, err := gomod.Find(dir)
	is.NoErr(err)
	pfs, err := pluginfs.Load(module)
	is.NoErr(err)

	// 1. Prefer application files
	code, err := fs.ReadFile(pfs, "public/admin.css")
	is.NoErr(err)
	is.Equal(string(code), `/* app admin.css */`)

	// 2. Prefer modules that are higher on the list alphanumerically.
	// In this case, `github.com/livebud/bud-test-nested-plugin` has priority over
	// `github.com/livebud/bud-test-plugin`.
	code, err = fs.ReadFile(pfs, "public/tailwind/preflight.css")
	is.NoErr(err)
	is.Equal(string(code), `/* conflicting preflight */`)

	// 3. Finally, ensure that non-conflicting files are still readable
	code, err = fs.ReadFile(pfs, "public/base.css")
	is.NoErr(err)
	is.Equal(string(code), `/* base */`)
}

func TestLocalPluginDir(t *testing.T) {
	t.SkipNow()
}
