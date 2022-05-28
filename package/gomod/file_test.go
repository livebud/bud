package gomod_test

import (
	"path/filepath"
	"testing"

	"github.com/livebud/bud/package/modcache"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/package/gomod"
)

func TestAddRequire(t *testing.T) {
	is := is.New(t)
	module, err := gomod.Parse("go.mod", []byte(`module app.test`))
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
	module, err := gomod.Parse("go.mod", []byte(`module app.test`))
	is.NoErr(err)
	modFile := module.File()
	modFile.AddReplace("mod.test/two", "", "mod.test/twotwo", "")
	modFile.AddReplace("mod.test/one", "", "mod.test/oneone", "")
	is.Equal(string(modFile.Format()), `module app.test

replace mod.test/two => mod.test/twotwo

replace mod.test/one => mod.test/oneone
`)
}

func TestLocalResolveDirectory(t *testing.T) {
	is := is.New(t)
	module, err := gomod.Parse("go.mod", []byte(`module app.test`))
	is.NoErr(err)
	modFile := module.File()
	modFile.AddRequire("github.com/livebud/bud-test-plugin", "v0.0.2")
	dir, err := module.ResolveDirectory("github.com/livebud/bud-test-plugin")
	is.NoErr(err)
	is.Equal(dir, filepath.Join(modcache.Default().Directory(), "github.com/livebud", "bud-test-plugin@v0.0.2"))
}

func TestFindRequire(t *testing.T) {
	is := is.New(t)
	module, err := gomod.Parse("go.mod", []byte("module app.test\nrequire mod.test/one v1.2.4"))
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

func TestRequires(t *testing.T) {
	is := is.New(t)
	module, err := gomod.Parse("go.mod", []byte(`
		module app.test
		require (
			github.com/andybalholm/cascadia v1.3.1 // indirect
			github.com/ajg/form v1.5.2-0.20200323032839-9aeb3cf462e1
			go.kuoruan.net/v8go-polyfills v0.5.0
			github.com/PuerkitoBio/goquery v1.8.0
			rogchap.com/v8go v0.7.0
		)
	`))
	is.NoErr(err)
	reqs := module.File().Requires()
	is.Equal(len(reqs), 5)
	is.Equal(reqs[0].Mod.Path, `github.com/PuerkitoBio/goquery`)
	is.Equal(reqs[1].Mod.Path, `github.com/ajg/form`)
	is.Equal(reqs[2].Mod.Path, `go.kuoruan.net/v8go-polyfills`)
	is.Equal(reqs[3].Mod.Path, `rogchap.com/v8go`)
	is.Equal(reqs[4].Mod.Path, `github.com/andybalholm/cascadia`)
}
