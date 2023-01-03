package config_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/livebud/bud/internal/config"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
	"github.com/livebud/bud/package/gomod"
)

func TestVersionAlignment(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	td.Modules["github.com/livebud/bud"] = "v0.1.7"
	is.NoErr(td.Write(ctx))
	module, err := gomod.Find(dir)
	is.NoErr(err)
	err = config.EnsureVersionAlignment(ctx, module, "0.1.8")
	is.NoErr(err)
	modFile, err := os.ReadFile(td.Path("go.mod"))
	is.NoErr(err)
	module, err = gomod.Parse(td.Path("go.mod"), modFile)
	is.NoErr(err)
	version := module.File().Require("github.com/livebud/bud")
	is.Equal(version.Version, "v0.1.8")
}

func TestGoVersion(t *testing.T) {
	is := is.New(t)
	is.NoErr(config.CheckGoVersion("go1.17"))
	is.NoErr(config.CheckGoVersion("go1.18"))
	is.True(errors.Is(config.CheckGoVersion("go1.16"), config.ErrMinGoVersion))
	is.True(errors.Is(config.CheckGoVersion("go1.16.5"), config.ErrMinGoVersion))
	is.True(errors.Is(config.CheckGoVersion("go1.8"), config.ErrMinGoVersion))
	is.NoErr(config.CheckGoVersion("abc123"))
}
