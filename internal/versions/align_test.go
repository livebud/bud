package versions_test

import (
	"context"
	"os"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/internal/versions"
	"github.com/livebud/bud/package/gomod"
)

func TestAlignRuntime(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Modules["github.com/livebud/bud"] = "v0.1.7"
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	err = versions.AlignRuntime(ctx, module, "0.1.8")
	is.NoErr(err)
	modFile, err := os.ReadFile(td.Path("go.mod"))
	is.NoErr(err)
	module, err = gomod.Parse(td.Path("go.mod"), modFile)
	is.NoErr(err)
	version := module.File().Require("github.com/livebud/bud")
	is.Equal(version.Version, "v0.1.8")
}
