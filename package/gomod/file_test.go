package gomod_test

import (
	"path/filepath"
	"testing"

	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/modcache"
	"github.com/matryer/is"
)

func TestAddRequire(t *testing.T) {
	is := is.New(t)
	modPath := filepath.Join(t.TempDir(), "go.mod")
	module, err := gomod.Parse(modPath, []byte(`module app.test`))
	is.NoErr(err)
	modFile := module.File()
	modFile.AddRequire("mod.test/two", "v2")
	modFile.AddRequire("mod.test/one", "v1.2.4")
	is.Equal(string(modFile.Format()), `module app.test

require (
	mod.test/two v2
	mod.test/one v1.2.4
)
`)
}

func TestAddReplace(t *testing.T) {
	is := is.New(t)
	modPath := filepath.Join(t.TempDir(), "go.mod")
	module, err := gomod.Parse(modPath, []byte(`module app.test`))
	is.NoErr(err)
	modFile := module.File()
	modFile.AddReplace("mod.test/two", "", "mod.test/twotwo", "")
	modFile.AddReplace("mod.test/one", "", "mod.test/oneone", "")
	is.Equal(string(modFile.Format()), `module app.test

replace mod.test/two => mod.test/twotwo

replace mod.test/one => mod.test/oneone
`)
}

func TestRequire(t *testing.T) {
	is := is.New(t)
	modPath := filepath.Join(t.TempDir(), "go.mod")
	m1, err := gomod.Parse(modPath, []byte(`module app.test`))
	is.NoErr(err)
	m2, err := gomod.Parse(modPath, []byte(`module app.test`))
	is.NoErr(err)
	is.Equal(string(m1.Hash()), string(m2.Hash()))
	m3, err := gomod.Parse(modPath, []byte(`module apptest`))
	is.NoErr(err)
	is.True(string(m2.Hash()) != string(m3.Hash()))
}

func TestLocalResolveDirectory(t *testing.T) {
	is := is.New(t)
	appDir := t.TempDir()
	modCache := modcache.Default()
	modPath := filepath.Join(appDir, "go.mod")
	modData := []byte(`module app.test`)
	module, err := gomod.Parse(modPath, modData)
	is.NoErr(err)
	modFile := module.File()
	modFile.AddRequire("github.com/livebud/bud-test-plugin", "v0.0.2")
	dir, err := module.ResolveDirectory("github.com/livebud/bud-test-plugin")
	is.NoErr(err)
	is.Equal(dir, filepath.Join(modCache.Directory(), "github.com/livebud", "bud-test-plugin@v0.0.2"))
}

func TestFindRequire(t *testing.T) {
	is := is.New(t)
	modPath := filepath.Join(t.TempDir(), "go.mod")
	module, err := gomod.Parse(modPath, []byte("module app.test\nrequire mod.test/one v1.2.4"))
	is.NoErr(err)
	modFile := module.File()
	modFile.AddRequire("mod.test/two", "v2")
	require := modFile.Require("mod.test/two")
	is.True(require != nil)
	is.Equal(require.Path, "mod.test/two")
	is.Equal(require.Version, "v2")
	require = modFile.Require("mod.test/one")
	is.True(require != nil)
	is.Equal(require.Path, "mod.test/one")
	is.Equal(require.Version, "v1.2.4")
	require = modFile.Require("mod.test/three")
	is.True(require == nil)
}
