package create_test

import (
	"context"
	"path/filepath"
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
	cli := testcli.New(dir)
	result, err := cli.Run(ctx, "create", "--module=github.com/my/app", filepath.Join(dir, "app"))
	is.NoErr(err)
	is.Equal(result.Stdout(), "")
	is.Equal(result.Stderr(), "")
}
