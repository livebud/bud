package cli_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testcli"
	"github.com/livebud/bud/package/testdir"
)

func TestBuildEmpty(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	err = td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(td.Directory())
	is.NoErr(td.NotExists("bud/app"))
	result, err := cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("bud/app"))
}

func TestBuildTwice(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(td.Directory())
	is.NoErr(td.NotExists("bud/app"))
	result, err := cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("bud/app"))
	result, err = cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists("bud/app"))
}
