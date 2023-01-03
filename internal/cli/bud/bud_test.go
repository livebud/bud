package bud_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
)

func TestHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	result, err := cli.Run(ctx, "--help")
	is.NoErr(err)
	is.Equal(result.Stderr(), "")
	is.In(result.Stdout(), "bud")
	is.In(result.Stdout(), "build")
	is.In(result.Stdout(), "run")
	is.In(result.Stdout(), "new")
	is.In(result.Stdout(), "tool")
	is.In(result.Stdout(), "version")
	is.NoErr(td.NotExists("bud/app"))
}

func TestChdir(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(".")
	result, err := cli.Run(ctx, "--chdir", dir, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists(
		"bud/internal/app",
		"bud/app",
	))
}

func TestChdirHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(".")
	result, err := cli.Run(ctx, "--chdir", dir, "--help")
	is.NoErr(err)
	is.Equal(result.Stderr(), "")
	is.In(result.Stdout(), "  bud")
	is.In(result.Stdout(), "  build")
	is.In(result.Stdout(), "  run")
	is.In(result.Stdout(), "  new")
	is.In(result.Stdout(), "  tool")
	is.In(result.Stdout(), "  version")
	is.NoErr(td.NotExists("bud/app"))
}

func TestOutsideModule(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	cli := testcli.New(dir)
	result, err := cli.Run(ctx)
	is.NoErr(err)
	is.Equal(result.Stderr(), "")
	is.In(result.Stdout(), "bud")
	is.In(result.Stdout(), "-C, --chdir")
	is.In(result.Stdout(), "-h, --help")
	is.In(result.Stdout(), "-L, --log")
	is.In(result.Stdout(), "  build")
	is.In(result.Stdout(), "  run")
	is.In(result.Stdout(), "  new")
	is.In(result.Stdout(), "  tool")
	is.In(result.Stdout(), "  version")
	is.NoErr(td.NotExists(
		"bud/internal/app",
		"bud/app",
	))
}

func TestOutsideModuleHelp(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	cli := testcli.New(dir)
	result, err := cli.Run(ctx, "--help")
	is.NoErr(err)
	is.Equal(result.Stderr(), "")
	is.In(result.Stdout(), "bud")
	is.In(result.Stdout(), "-C, --chdir")
	is.In(result.Stdout(), "-h, --help")
	is.In(result.Stdout(), "-L, --log")
	is.In(result.Stdout(), "  build")
	is.In(result.Stdout(), "  run")
	is.In(result.Stdout(), "  new")
	is.In(result.Stdout(), "  tool")
	is.In(result.Stdout(), "  version")
	is.NoErr(td.NotExists(
		"bud/internal/app",
		"bud/app",
	))
}

// TODO: This test might go away at some point. Right now testcli isn't setup
// to run `bud build`, then start the built process. It's setup for running
// commands to completion and `bud run`. This is alright because `bud run`
// currently can behave exactly like `bud build` by setting the right flags.
// This test aims to maintain this.
func TestBuildRunAlignment(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(dir)
	buildResult, err := cli.Run(ctx, "build", "--help")
	is.NoErr(err)
	is.Equal(buildResult.Stderr(), "")
	runResult, err := cli.Run(ctx, "run", "--help")
	is.NoErr(err)
	is.Equal(runResult.Stderr(), "")
	is.In(buildResult.Stdout(), "build")
	is.In(buildResult.Stdout(), "--minify")
	is.In(buildResult.Stdout(), "--embed")
	is.In(runResult.Stdout(), "run")
	is.In(runResult.Stdout(), "--minify")
	is.In(runResult.Stdout(), "--embed")
}
