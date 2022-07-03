package create_test

import (
	"context"
	"testing"

	"github.com/livebud/bud/internal/cli/testcli"
	"github.com/livebud/bud/internal/is"
	"github.com/livebud/bud/internal/testdir"
)

func TestCreateOutsideGoPathError(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	err := td.Write(ctx)
	is.NoErr(err)
	cli := testcli.New(dir)
	is.NoErr(td.NotExists(".gitignore"))
	result, err := cli.Run(ctx, "create", dir)
	is.In(err.Error(), `Try again using the module <path> name.`)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.NotExists(".gitignore"))
}

func TestCreateOutsideGoPathModulePath(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	cli := testcli.New(dir)
	is.NoErr(td.NotExists(".gitignore"))
	result, err := cli.Run(ctx, "create", "--module=github.com/my/app", dir)
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
	is.NoErr(td.Exists(".gitignore"))
	is.NoErr(td.Exists("go.sum"))
	is.NoErr(td.Exists("package.json"))
	is.NoErr(td.Exists("package-lock.json"))
}
