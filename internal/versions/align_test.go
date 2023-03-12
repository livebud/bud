package versions_test

import (
	"context"
	"io/fs"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/gomod"
	"github.com/livebud/bud/package/testdir"
)

func TestAlignRuntime(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	td.Modules["github.com/livebud/bud"] = "v0.1.7"
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(td.Directory())
	is.NoErr(err)
	err = versions.AlignRuntime(ctx, module, "0.1.8")
	is.NoErr(err)
	modFile, err := fs.ReadFile(td, "go.mod")
	is.NoErr(err)
	module, err = gomod.Parse("go.mod", modFile)
	is.NoErr(err)
	version := module.File().Require("github.com/livebud/bud")
	is.Equal(version.Version, "v0.1.8")
}
