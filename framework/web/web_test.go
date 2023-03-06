package web_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testcli"
	"github.com/livebud/bud/package/testdir"
)

func TestEmptyBuild(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	td, err := testdir.Load()
	is.NoErr(err)
	is.NoErr(td.Write(ctx))
	cli := testcli.New(td.Directory())
	is.NoErr(td.NotExists("bud/internal/web"))
	result, err := cli.Run(ctx, "build")
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	// Empty builds generate the web directory
	is.NoErr(td.Exists("bud/internal/web"))
}
